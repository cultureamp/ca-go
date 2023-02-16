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
	batchSize        int
	handlerExecutor  *handlerExecutor
	reader           Reader
	fetched          chan kafka.Message
	processed        chan kafka.Message
	stop             chan struct{}
	getOrderingKeyFn GetOrderingKey
	fetchCancel      func()
}

func newBatchProcessor(reader Reader, handlerExecutor *handlerExecutor, getOrderingKeyFn GetOrderingKey, batchSize int) *batchProcessor {
	return &batchProcessor{
		batchSize:        batchSize,
		handlerExecutor:  handlerExecutor,
		reader:           reader,
		fetched:          make(chan kafka.Message, batchSize),
		processed:        make(chan kafka.Message, batchSize),
		stop:             make(chan struct{}),
		getOrderingKeyFn: getOrderingKeyFn,
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

	var commits []kafka.Message

	for msg := range b.processed {
		commits = append(commits, msg)
	}

	if err := b.reader.CommitMessages(ctx, commits...); err != nil {
		return fmt.Errorf("unable to commit messages for batch: %w", err)
	}
	return nil
}

func (b *batchProcessor) startFetching(ctx context.Context) error {
	var fetchContext context.Context
	fetchContext, b.fetchCancel = context.WithCancel(ctx)

	for i := 0; i < b.batchSize; i++ {
		msg, err := b.reader.FetchMessage(fetchContext)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				select {
				case <-b.stop:
					return nil
				default:
				}
			} else if !errors.Is(err, io.EOF) {
				err = fmt.Errorf("unable to fetch message: %w", err)
			}
			return err
		}
		b.fetched <- msg
	}
	close(b.fetched)
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
	for {
		msg, ok := b.nextMessage()
		if !ok {
			break
		}

		messagesReceived++
		key := b.getOrderingKeyFn(ctx, msg)
		if orderedChan, ok := orderedChans[key]; ok {
			orderedChan <- msg
			continue
		}

		msgCh := make(chan kafka.Message, b.batchSize)
		msgCh <- msg
		orderedChans[key] = msgCh

		handleErrg.Go(func() error {
			for m := range msgCh {
				if err := b.handlerExecutor.execute(handleCtx, m, handler); err != nil {
					return err
				}
				b.processed <- m
				messagesHandled.inc()
			}
			return nil
		})

		// End batch process early if there are no new messages available to
		// be fetched. This is to avoid unnecessary lag.
		for !b.hasNext() {
			if messagesReceived == messagesHandled.val() {
				b.stopFetching()
				break processLoop
			}
		}
	}

	for _, msgCh := range orderedChans {
		close(msgCh)
	}
	err := handleErrg.Wait()
	return err
}

func (b *batchProcessor) reset() {
	b.fetched = make(chan kafka.Message, b.batchSize)
	b.processed = make(chan kafka.Message, b.batchSize)
	b.stop = make(chan struct{})
}
