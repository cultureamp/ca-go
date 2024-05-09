//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -destination=mock_reader_test.go -package consumer . Reader
package consumer

import (
	"context"
	"io"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	kafkatrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/segmentio/kafka.go.v0"
)

const (
	consumerMinBytes      = 1e6  // 1 MB
	consumerMaxBytes      = 10e6 // 10 MB
	consumerMaxWait       = 250 * time.Millisecond
	consumerQueueCapacity = 100
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
	groupID       string
}

// Consumer provides a high level API for consuming and handling messages from
// a Kafka topic.
//
// It is worth noting that publishing failed messages to a dead letter queue is
// not supported and instead would need to be included in your handler implementation.
type Consumer struct {
	id                 string
	conf               kafka.ReaderConfig
	reader             Reader
	withExplicitCommit bool
	stopCh             chan struct{}
	clientHandler      *messageHandler
}

// NewConsumer returns a new Consumer configured with the provided dialer and config.
func NewConsumer(config Config, opts ...Option) *Consumer {
	if config.ID == "" {
		config.ID = uuid.New().String()
	}
	if config.MinBytes == 0 {
		config.MinBytes = consumerMinBytes // 1 MB
	}
	if config.MaxBytes == 0 {
		config.MaxBytes = consumerMaxBytes // 10 MB
	}
	if config.MaxWait == 0 {
		config.MaxWait = consumerMaxWait // 250ms
	}
	if config.QueueCapacity < 1 {
		config.QueueCapacity = consumerQueueCapacity // 100
	}

	c := &Consumer{
		id:     config.ID,
		stopCh: make(chan struct{}),
		conf: kafka.ReaderConfig{
			Brokers:               config.Brokers,
			GroupID:               config.groupID,
			Topic:                 config.Topic,
			Dialer:                kafka.DefaultDialer,
			WatchPartitionChanges: true,
			MaxBytes:              config.MaxBytes,
			Logger:                kafka.LoggerFunc(func(string, ...interface{}) {}), // default to noop
			ErrorLogger:           kafka.LoggerFunc(func(string, ...interface{}) {}), // default to noop
		},
		clientHandler: &messageHandler{
			ConsumerID:   config.ID,
			GroupID:      config.groupID,
			clientNotify: func(_ context.Context, _ error, _ Message) {}, // default to noop
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	// Set the reader unless one was injected via the WithKafkaReader option.
	if c.reader == nil {
		if c.clientHandler.DataDogTracingEnabled {
			c.reader = kafkatrace.NewReader(c.conf)
		} else {
			c.reader = kafka.NewReader(c.conf)
		}
	}

	return c
}

// Run consumes and handles messages from the topic. The method call blocks until
// the context is canceled, the consumer is closed, or an error occurs.
func (c *Consumer) Run(ctx context.Context, handler Handler) error {
	c.conf.Logger.Printf(
		"consumer(%s:%s): running until context is cancelled, an error occurs, or the consumer is stopped",
		c.conf.Topic,
		c.id,
	)

	// Run forever until we read from the stopCh or we have an error processing a message
	for {
		select {
		case <-c.stopCh:
			c.conf.Logger.Printf(
				"consumer(%s:%s): stopped signal received",
				c.conf.Topic,
				c.id,
			)
			return nil
		default:
		}

		if err := c.retreiveNextMessage(ctx, handler); err != nil {
			return errors.Errorf("consumer error: %w", err)
		}
	}
}

// Stop stops the consumer. It waits for the current message (if any) to
// finish being handled before closing the reader stream, preventing the consumer
// from reading any more messages.
func (c *Consumer) Stop() error {
	close(c.stopCh)
	if err := c.reader.Close(); err != nil {
		return errors.Errorf("unable to close consumer reader: %w", err)
	}

	c.conf.Logger.Printf(
		"consumer(%s:%s): consumer has stopped",
		c.conf.Topic,
		c.id,
	)
	return nil
}

func (c *Consumer) retreiveNextMessage(ctx context.Context, handler Handler) error {
	if c.withExplicitCommit {
		return c.fetchNextMessage(ctx, handler)
	}

	return c.readNextMessage(ctx, handler)
}

func (c *Consumer) fetchNextMessage(ctx context.Context, handler Handler) error {
	var msg kafka.Message
	var err error

	msg, err = c.reader.FetchMessage(ctx)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return errors.Errorf("unable to fetch message: %w", err)
	}

	if err = c.clientHandler.dispatch(ctx, msg, handler); err != nil {
		return errors.Errorf("unable to handle message: %w", err)
	}

	if err = c.reader.CommitMessages(ctx, msg); err != nil {
		return errors.Errorf("unable to commit message: %w", err)
	}

	return nil
}

func (c *Consumer) readNextMessage(ctx context.Context, handler Handler) error {
	var msg kafka.Message
	var err error

	msg, err = c.reader.ReadMessage(ctx)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return errors.Errorf("unable to read message: %w", err)
	}

	if err = c.clientHandler.dispatch(ctx, msg, handler); err != nil {
		return errors.Errorf("unable to handle message: %w", err)
	}

	return nil
}
