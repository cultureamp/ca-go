package consumer

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

type batchProcessor struct {
	consumerID       string
	batchSize        int
	handlerExecutor  *handlerExecutor
	reader           Reader
	fetched          chan kafka.Message
	processed        chan kafka.Message
	stop             chan struct{}
	getOrderingKeyFn GetOrderingKey
	fetchCancel      func()
	debugLogger      DebugLogger
	debugKeyVals     []any
}

func newBatchProcessor(consumerID string, debugLogger DebugLogger, reader Reader, handlerExecutor *handlerExecutor, getOrderingKeyFn GetOrderingKey, batchSize int) *batchProcessor {
	return &batchProcessor{
		consumerID:       consumerID,
		batchSize:        batchSize,
		handlerExecutor:  handlerExecutor,
		reader:           reader,
		fetched:          make(chan kafka.Message, batchSize),
		processed:        make(chan kafka.Message, batchSize),
		stop:             make(chan struct{}),
		getOrderingKeyFn: getOrderingKeyFn,
		debugLogger:      debugLogger,
		debugKeyVals:     []any{"consumerId", consumerID, "batchSize", batchSize},
	}
}

func (b *batchProcessor) process(ctx context.Context, handler Handler) error {
	defer b.reset()
	errg, errgCtx := errgroup.WithContext(ctx)
	errg.Go(func() error {
		return b.startFetching(errgCtx)
	})

	errg.Go(func() error {
		return b.startProcessing(errgCtx, handler)
	})

	if err := errg.Wait(); err != nil {
		return err
	}
	close(b.processed)
	b.debugLogger.Print(fmt.Sprintf("Fetching and handling finished for %d messages in batch", len(b.processed)), b.debugKeyVals...)

	var commits []kafka.Message

	for msg := range b.processed {
		commits = append(commits, msg)
	}

	if err := b.reader.CommitMessages(ctx, commits...); err != nil {
		return fmt.Errorf("unable to commit messages for batch: %w", err)
	}
	b.debugLogger.Print(fmt.Sprintf("Committed %d message offsets for batch", len(commits)), b.debugKeyVals...)
	return nil
}

func (b *batchProcessor) startFetching(ctx context.Context) error {
	var fetchContext context.Context
	fetchContext, b.fetchCancel = context.WithCancel(ctx)

	b.debugLogger.Print("Fetching messages for batch", b.debugKeyVals...)
	defer b.debugLogger.Print("Finished fetching messages for batch", b.debugKeyVals...)

	batchSize := b.batchSize - len(b.fetched)
	for i := 0; i < batchSize; i++ {
		msg, err := b.reader.FetchMessage(fetchContext)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				select {
				case <-b.stop:
					b.debugLogger.Print("Fetch message for batch stopped", append([]any{"messagesFetched", i}, b.debugKeyVals...)...)
					return nil
				default:
				}
			} else if !errors.Is(err, io.EOF) {
				err = fmt.Errorf("unable to fetch message: %w", err)
			}
			return err
		}
		b.debugLogger.Print("Fetched message", append([]any{"partition", msg.Partition, "offset", msg.Offset}, b.debugKeyVals...)...)
		b.fetched <- msg
	}
	return nil
}

func (b *batchProcessor) nextMessage() (kafka.Message, bool) {
	msg, ok := <-b.fetched
	return msg, ok
}

func (b *batchProcessor) hasNext() bool {
	return len(b.fetched) > 0
}

func (b *batchProcessor) stopFetching() {
	b.fetchCancel()
	close(b.stop)
}

func (b *batchProcessor) startProcessing(ctx context.Context, handler Handler) error {
	messagesReceived := 0
	messagesHandled := new(safeCounter)
	orderedChans := make(map[string]chan kafka.Message)
	handleErrg, handleCtx := errgroup.WithContext(ctx)

processLoop:
	for i := 0; i < b.batchSize; i++ {
		msg, ok := b.nextMessage()
		if !ok {
			break
		}

		messagesReceived++
		key := b.getOrderingKeyFn(ctx, msg)
		b.debugLogger.Print("Queuing message for handling",
			append([]any{"partition", msg.Partition, "offset", msg.Offset, "orderingKey", key}, b.debugKeyVals...)...,
		)

		if orderedChan, ok := orderedChans[key]; ok {
			orderedChan <- msg
			continue
		}

		msgCh := make(chan kafka.Message, b.batchSize)
		msgCh <- msg
		orderedChans[key] = msgCh

		handleErrg.Go(func() error {
			debugKeyVals := append([]any{"partition", msg.Partition, "offset", msg.Offset, "orderingKey", key}, b.debugKeyVals...)
			b.debugLogger.Print(fmt.Sprintf("Handling messages for ordering key %s", key), debugKeyVals...)
			defer b.debugLogger.Print(fmt.Sprintf("Finished handling all messages for ordering key %s", key), debugKeyVals...)

			for m := range msgCh {
				b.debugLogger.Print("Executing message handler", debugKeyVals...)
				if err := b.handlerExecutor.execute(handleCtx, m, handler); err != nil {
					b.debugLogger.Print("Message handler execution failed", debugKeyVals...)
					return err
				}
				b.processed <- m
				messagesHandled.inc()
				b.debugLogger.Print("Finished message handler execution", debugKeyVals...)
			}
			return nil
		})

		// End batch process early if there are no new messages available to
		// be fetched. This is to avoid unnecessary lag.
		for !b.hasNext() {
			handled := messagesHandled.val()
			if messagesReceived == handled {
				b.debugLogger.Print("Stopping batch early: all current messages handled and no new messages to fetch",
					append([]any{"messagesHandled", handled}, b.debugKeyVals)...,
				)
				b.stopFetching()
				break processLoop
			}
		}
	}

	for _, msgCh := range orderedChans {
		close(msgCh)
	}
	err := handleErrg.Wait()
	b.debugLogger.Print("Finished handling messages for batch",
		append([]any{"messagesHandled", messagesHandled.val()}, b.debugKeyVals)...,
	)
	return err
}

func (b *batchProcessor) reset() {
	b.processed = make(chan kafka.Message, b.batchSize)
	b.stop = make(chan struct{})
}

func (b *batchProcessor) close() {
	close(b.fetched)
}
