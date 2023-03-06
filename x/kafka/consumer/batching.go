package consumer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

type batchProcessorConfig struct {
	consumerID       string
	batchSize        int
	fetchDuration    time.Duration
	debugLogger      DebugLogger
	getOrderingKeyFn GetOrderingKey
	handlerExecutor  *handlerExecutor
	reader           Reader
}

type batchProcessor struct {
	consumerID       string
	batchSize        int
	handlerExecutor  *handlerExecutor
	reader           Reader
	fetchTimeout     time.Duration
	fetched          chan kafka.Message
	processed        chan kafka.Message
	getOrderingKeyFn GetOrderingKey
	debugLogger      DebugLogger
	debugKeyVals     []any
}

func newBatchProcessor(config batchProcessorConfig) *batchProcessor {
	if config.fetchDuration == time.Duration(0) {
		config.fetchDuration = time.Second * 10
	}
	return &batchProcessor{
		consumerID:       config.consumerID,
		batchSize:        config.batchSize,
		handlerExecutor:  config.handlerExecutor,
		reader:           config.reader,
		fetchTimeout:     config.fetchDuration,
		getOrderingKeyFn: config.getOrderingKeyFn,
		debugLogger:      config.debugLogger,
	}
}

func (b *batchProcessor) process(ctx context.Context, handler Handler) error {
	b.processed = make(chan kafka.Message, b.batchSize)
	b.fetched = make(chan kafka.Message, b.batchSize)
	b.debugKeyVals = []any{"consumerId", b.consumerID, "batchSize", b.batchSize, "batchId", uuid.New().String()}

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

	b.debugLogger.Print(fmt.Sprintf("Preparing to commit offsets for %d messages in batch", len(b.processed)), b.debugKeyVals...)

	var commits []kafka.Message

	for msg := range b.processed {
		commits = append(commits, msg)
	}

	if err := b.reader.CommitMessages(ctx, commits...); err != nil {
		return fmt.Errorf("unable to commit messages for batch: %w", err)
	}
	b.debugLogger.Print(fmt.Sprintf("Committed offsets for %d messages in batch", len(commits)), b.debugKeyVals...)
	return nil
}

func (b *batchProcessor) startFetching(ctx context.Context) error {
	b.debugLogger.Print("Fetching messages for batch", b.debugKeyVals...)
	defer b.debugLogger.Print("Finished fetching messages for batch", b.debugKeyVals...)

	fetchCtx, cancel := context.WithTimeout(ctx, b.fetchTimeout)
	defer cancel()

	defer close(b.fetched)

	for i := 0; i < b.batchSize; i++ {
		msg, err := b.reader.FetchMessage(fetchCtx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				b.debugLogger.Print("Fetching stopped early due to context deadline exceeded", b.debugKeyVals...)
				return nil
			}
			return err
		}
		b.fetched <- msg
	}

	return nil
}

func (b *batchProcessor) startProcessing(ctx context.Context, handler Handler) error {
	orderedChans := make(map[string]chan kafka.Message)
	handleErrg, handleCtx := errgroup.WithContext(ctx)

processLoop:
	for i := 0; i < b.batchSize; i++ {
		var msg kafka.Message

		select {
		case <-ctx.Done():
			return ctx.Err()
		case fetched, ok := <-b.fetched:
			if !ok {
				break processLoop
			}
			msg = fetched
		}

		key := b.getOrderingKeyFn(ctx, msg)
		b.debugLogger.Print("Queuing message for handling",
			append([]any{"partition", msg.Partition, "offset", msg.Offset, "orderingKey", key}, b.debugKeyVals...)...,
		)

		if orderedChan, ok := orderedChans[key]; ok {
			orderedChan <- msg
		} else {
			orderdChan := make(chan kafka.Message, b.batchSize)
			orderdChan <- msg
			orderedChans[key] = orderdChan
			handleErrg.Go(func() error {
				return b.handleMessages(handleCtx, key, orderdChan, handler)
			})
		}
	}

	for _, msgCh := range orderedChans {
		close(msgCh)
	}

	b.debugLogger.Print("Waiting for all batch messages to be handled", b.debugKeyVals...)
	err := handleErrg.Wait()
	b.debugLogger.Print("Finished handling messages for batch", b.debugKeyVals...)
	return err
}

func (b *batchProcessor) handleMessages(ctx context.Context, orderingKey string, orderedChan chan kafka.Message, handler Handler) error {
	b.debugLogger.Print(fmt.Sprintf("Handling messages for ordering key %s", orderingKey))
	defer b.debugLogger.Print(fmt.Sprintf("Finished handling all messages for ordering key %s", orderingKey))

	for m := range orderedChan {
		debugKeyVals := append([]any{"partition", m.Partition, "offset", m.Offset, "orderingKey", orderingKey}, b.debugKeyVals...)
		b.debugLogger.Print("Executing message handler", debugKeyVals...)
		if err := b.handlerExecutor.execute(ctx, m, handler); err != nil {
			b.debugLogger.Print("Message handler execution failed", debugKeyVals...)
			return err
		}
		b.processed <- m
		b.debugLogger.Print("Finished message handler execution", debugKeyVals...)
	}
	return nil
}
