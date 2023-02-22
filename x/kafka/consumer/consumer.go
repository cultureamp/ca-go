//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -destination=mock_reader_test.go -package consumer . Reader
package consumer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	kafkatrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/segmentio/kafka.go.v0"
)

// Metadata contains relevant handler metadata for received Kafka messages.
type Metadata struct {
	GroupID    string
	ConsumerID string
	Attempt    int
}

type Message struct {
	kafka.Message
	Metadata
}

// Handler specifies how a consumer should handle a received Kafka message.
type Handler func(ctx context.Context, msg Message) error

// NotifyError is a notify-on-error function used to report consumer handler errors.
type NotifyError func(ctx context.Context, err error, msg Message)

// GetOrderingKey specifies what key to store the Kafka message under when
// processing in batches. Ordering keys are used to spawn new goroutines that
// are responsible for processing each message for that key in order. An ordering
// key is also useful for decreasing/increasing processing concurrency within
// a batch.
type GetOrderingKey func(ctx context.Context, message kafka.Message) string

// Reader fetches and commits messages from a Kafka topic.
type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Config is a configuration object used to create a new Consumer.
type Config struct {
	ID      string // Default: UUID
	Brokers []string
	Topic   string

	MinBytes      int           // Default: 1MB
	MaxBytes      int           // Default: 10MB
	MaxWait       time.Duration // Default: 250ms
	QueueCapacity int           // Default: 100

	DebugLogger DebugLogger
	groupID     string
}

// Consumer provides a high level API for consuming and handling messages from
// a Kafka topic.
//
// It is worth noting that publishing failed messages to a dead letter queue is
// not supported and instead would need to be included in your handler implementation.
type Consumer struct {
	id                 string
	reader             Reader
	readerConfig       kafka.ReaderConfig
	withExplicitCommit bool
	batchSize          int
	fetchDuration      time.Duration
	getOrderingKeyFn   GetOrderingKey
	stopCh             chan struct{}
	handlerExecutor    *handlerExecutor
	debugLogger        DebugLogger
	debugKeyVals       []any
}

// NewConsumer returns a new Consumer configured with the provided dialer and config.
func NewConsumer(dialer *kafka.Dialer, config Config, opts ...Option) *Consumer {
	if config.ID == "" {
		config.ID = uuid.New().String()
	}
	if config.MaxBytes == 0 {
		config.MaxBytes = 1e6 // 1 MB
	}
	if config.MaxBytes == 0 {
		config.MaxBytes = 10e6 // 10 MB
	}
	if config.MaxWait == 0 {
		config.MaxWait = 250 * time.Millisecond
	}
	if config.QueueCapacity == 0 {
		config.QueueCapacity = 100
	}

	c := &Consumer{
		id:     config.ID,
		stopCh: make(chan struct{}),
		readerConfig: kafka.ReaderConfig{
			Brokers:               config.Brokers,
			GroupID:               config.groupID,
			Topic:                 config.Topic,
			Dialer:                dialer,
			WatchPartitionChanges: true,
			MaxBytes:              config.MaxBytes,
		},
		handlerExecutor: &handlerExecutor{
			ConsumerID: config.ID,
			GroupID:    config.groupID,
		},
		batchSize:        0,
		getOrderingKeyFn: func(ctx context.Context, message kafka.Message) string { return "" },
		debugLogger:      config.DebugLogger,
		debugKeyVals:     []any{"consumerId", config.ID},
	}

	if c.debugLogger == nil {
		c.debugLogger = noopDebugLogger{}
	}

	for _, opt := range opts {
		opt(c)
	}

	// Set the reader unless one was injected via the WithKafkaReader option.
	if c.reader == nil {
		if c.handlerExecutor.DataDogTracingEnabled {
			c.reader = kafkatrace.NewReader(c.readerConfig)
		} else {
			c.reader = kafka.NewReader(c.readerConfig)
		}
	}

	return c
}

// Run consumes and handles messages from the topic. The method call blocks until
// the context is canceled, the consumer is closed, or an error occurs.
func (c *Consumer) Run(ctx context.Context, handler Handler) error {
	c.debugLogger.Print("Running consumer", c.debugKeyVals...)
	bp := newBatchProcessor(batchProcessorConfig{
		consumerID:       c.id,
		batchSize:        c.batchSize,
		fetchDuration:    c.fetchDuration,
		debugLogger:      c.debugLogger,
		getOrderingKeyFn: c.getOrderingKeyFn,
		handlerExecutor:  c.handlerExecutor,
		reader:           c.reader,
	})

	for {
		select {
		case <-c.stopCh:
			c.debugLogger.Print("Consumer stopped", c.debugKeyVals...)
			return nil
		default:
		}

		if c.batchSize > 0 {
			if err := bp.process(ctx, handler); err != nil {
				return fmt.Errorf("consumer %s batch error: %w", c.id, err)
			}
		} else {
			if err := c.process(ctx, handler); err != nil {
				return fmt.Errorf("consumer %s error: %w", c.id, err)
			}
		}
	}
}

// Stop stops the consumer. It waits for the current message/batch (if any) to
// finish being handled before closing the reader stream, preventing the consumer
// from reading any more messages.
func (c *Consumer) Stop() error {
	close(c.stopCh)
	c.debugLogger.Print("Consumer stopped", c.debugKeyVals...)
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("unable to close consumer %s reader: %w", c.id, err)
	}
	c.debugLogger.Print("Consumer reader closed", c.debugKeyVals...)
	return nil
}

func (c *Consumer) process(ctx context.Context, handler Handler) error {
	var msg kafka.Message
	var err error

	if c.withExplicitCommit {
		msg, err = c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("unable to fetch message: %w", err)
		}
	} else {
		msg, err = c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("unable to read message: %w", err)
		}
	}

	debugKeyVals := append([]any{"partition", msg.Partition, "offset", msg.Offset}, c.debugKeyVals...)
	c.debugLogger.Print("Fetched message", debugKeyVals...)

	if err = c.handlerExecutor.execute(ctx, msg, handler); err != nil {
		return fmt.Errorf("unable to handle message: %w", err)
	}
	c.debugLogger.Print("Message handler execution finished", debugKeyVals...)

	if c.withExplicitCommit {
		if err = c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("unable to commit message: %w", err)
		}
	}
	c.debugLogger.Print("Committed message offset", debugKeyVals...)

	return nil
}

// GroupConfig is a configuration object used to create a new Group. The default
// consumer count in a group is 1 unless specified otherwise.
type GroupConfig struct {
	Count   int
	Brokers []string
	Topic   string
	GroupID string

	MinBytes      int           // Default: 1MB
	MaxBytes      int           // Default: 10MB
	MaxWait       time.Duration // Default: 250ms
	QueueCapacity int           // Default: 100

	DebugLogger DebugLogger
}

// Group groups consumers together to concurrently consume and handle messages
// from a Kafka topic. Many groups with the same group ID are safe to use, which
// is particularly useful for groups across separate instances.
//
// It is worth noting that publishing failed messages to a dead letter queue is
// not supported and instead would need to be included in your handler implementation.
type Group struct {
	ID      string
	config  GroupConfig
	opts    []Option
	dialer  *kafka.Dialer
	stopChs []chan struct{}
}

// NewGroup returns a new Group configured with the provided dialer and config.
func NewGroup(dialer *kafka.Dialer, config GroupConfig, opts ...Option) *Group {
	if config.Count <= 0 {
		config.Count = 1
	}

	return &Group{
		ID:     fmt.Sprintf("%s-%s", strings.ToLower(config.GroupID), uuid.New().String()[:7]), // semi-random slug
		config: config,
		dialer: dialer,
		opts:   opts,
	}
}

// Run concurrently consumes and handles messages from the topic across all
// consumers in the group. The method call returns an error channel that is used
// to receive any consumer errors. The run process is only stopped if the context
// is canceled, the consumer has been closed, or all consumers in the group have
// errored.
func (g *Group) Run(ctx context.Context, handler Handler) <-chan error {
	var wg sync.WaitGroup
	errCh := make(chan error, g.config.Count)

	for i := 0; i < g.config.Count; i++ {
		wg.Add(1)

		// Consumers must be created and run in sequential order so that Kafka can
		// successfully re-balance the group as each is added. This unfortunately
		// prevents us from receiving a passed in list of consumers, which is
		// arguably a cleaner approach.
		cfg := Config{
			ID:          fmt.Sprintf("%s-%d", g.ID, i),
			Brokers:     g.config.Brokers,
			Topic:       g.config.Topic,
			MaxBytes:    g.config.MaxBytes,
			groupID:     g.config.GroupID,
			DebugLogger: g.config.DebugLogger,
		}
		c := NewConsumer(g.dialer, cfg, g.opts...)

		go func() {
			defer wg.Done()
			if err := c.Run(ctx, handler); err != nil {
				errCh <- fmt.Errorf("consumer %s for group %s failed: %w", c.id, g.ID, err)
			}
		}()
		g.stopChs = append(g.stopChs, c.stopCh)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return errCh
}

// Stop stops the group. It waits for the current message/batch (if any) in each
// consumer to finish being handled before closing the reader streams, preventing
// each consumer from reading any more messages.
func (g *Group) Stop() {
	for _, stopCh := range g.stopChs {
		close(stopCh)
	}
}

// DialerSCRAM512 returns a Kafka dialer configured with SASL authentication
// to securely transmit the provided credentials to Kafka using SCRAM-SHA-512.
func DialerSCRAM512(username string, password string) (*kafka.Dialer, error) {
	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		return nil, err
	}

	return &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS:           &tls.Config{MinVersion: tls.VersionTLS12},
	}, nil
}

type HandlerRetryBackOffConstructor func() backoff.BackOff

// NonStopExponentialBackOff is the suggested backoff retry strategy for consumers
// groups handling messages where ordering matters e.g. a data-capture stream from
// a database. This results in endless retries to prevent Kafka from re-balancing
// the group so that each consumer does not eventually experience the same error.
//
// Retry intervals: 500ms, 4s, 32s, 4m, 34m, 4.5h, 5h (max).
//
// The max interval of 5 hours is intended to leave enough time for manual
// intervention if necessary.
func NonStopExponentialBackOff() backoff.BackOff { //nolint:ireturn
	bo := backoff.NewExponentialBackOff()
	bo.RandomizationFactor = 0
	bo.MaxInterval = 5 * time.Hour
	bo.Multiplier = 8
	bo.MaxElapsedTime = 0
	return bo
}
