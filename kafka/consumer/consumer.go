//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -destination=mock_reader_test.go -package consumer . Reader
package consumer

import (
	"context"
	"io"
	"time"

	"github.com/cultureamp/ca-go/log"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

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
	reader             Reader
	readerConfig       kafka.ReaderConfig
	withExplicitCommit bool
	stopCh             chan struct{}
	clientHandler      *readerMessageHandler
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
	if config.QueueCapacity < 1 {
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
		clientHandler: &readerMessageHandler{
			ConsumerID: config.ID,
			GroupID:    config.groupID,
			Notify:     func(ctx context.Context, err error, msg Message) {}, // default to noop
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	// Set the reader unless one was injected via the WithKafkaReader option.
	if c.reader == nil {
		if c.clientHandler.DataDogTracingEnabled {
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
	log.Debug("consumer_run").
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("id", c.id).
			Str("topic", c.readerConfig.Topic),
		).Details("running until context is cancelled, or an error occurs, or the consumer is Stop()'ed")

	// Run forever until we read from the stopCh or we have an error processing a message
	for {
		select {
		case <-c.stopCh:
			log.Info("consumer_run").
				WithSystemTracing().
				Properties(log.SubDoc().
					Str("id", c.id).
					Str("topic", c.readerConfig.Topic),
				).Details("stopped signal received")
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

	log.Debug("consumer_stop").
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("id", c.id).
			Str("topic", c.readerConfig.Topic),
		).Details("consumer has stopped")

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

	if err = c.clientHandler.execute(ctx, msg, handler); err != nil {
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

	if err = c.clientHandler.execute(ctx, msg, handler); err != nil {
		return errors.Errorf("unable to handle message: %w", err)
	}

	return nil
}
