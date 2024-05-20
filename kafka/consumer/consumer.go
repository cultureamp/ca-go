package consumer

import (
	"context"

	"github.com/go-errors/errors"
)

type KafkaConsumer interface {
	Start(ctx context.Context) error
}

// Consumer provides a high level API for consuming and handling messages from
// a Kafka topic.
type Consumer struct {
	client kafkaClient // Kafka client interfaces (Default: Sarama)
	conf   *Config
	group  *groupConsumer
}

// NewConsumer returns a new Consumer configured with the provided dialer and config.
func NewConsumer(opts ...Option) (*Consumer, error) {
	c := &Consumer{
		conf:   newConfig(),
		client: newSaramaClient(),
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.conf.mustProcess(); err != nil {
		return nil, errors.Errorf("bad consumer config: %w", err)
	}

	return c, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	if c.group != nil {
		return errors.Errorf("already running! (forgot to call Stop?)")
	}

	group, err := newGroupConsumer(c.client, c.conf)
	if err != nil {
		return errors.Errorf("failed to create kafka consumer: %w", err)
	}

	c.group = group
	// is this correct? run in a go-routine?
	err = group.consume(ctx)
	return err
}

func (c *Consumer) Stop() error {
	if c.group == nil {
		return nil
	}

	// err := c.group.stop()
	c.group = nil
	return nil
}
