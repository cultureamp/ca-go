package consumer

import (
	"context"

	"github.com/go-errors/errors"
)

type Cleanup func() error

type KafkaConsumer interface {
	Consume(ctx context.Context) (Cleanup, error)
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

	if err := c.conf.shouldProcess(); err != nil {
		return nil, errors.Errorf("bad consumer config: %w", err)
	}

	return c, nil
}

func (c *Consumer) Consume(ctx context.Context) (Cleanup, error) {
	// if already consuming, do nothing
	if c.group != nil {
		return c.stop, nil
	}

	group, err := newGroupConsumer(c.client, c.conf)
	if err != nil {
		return c.stop, errors.Errorf("failed to create kafka consumer: %w", err)
	}
	c.group = group

	// blocking call until either
	// 1. context is cancelled OR
	// 2. a server-side kafka rebalance happens
	return c.stop, c.group.consume(ctx)
}

func (c *Consumer) stop() error {
	// if already stopped, do nothing
	if c.group == nil {
		return nil
	}

	return c.group.stop()
}
