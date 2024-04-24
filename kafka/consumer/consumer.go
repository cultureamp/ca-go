package consumer

import (
	"context"

	"github.com/go-errors/errors"
)

// Consumer provides a high level API for consuming and handling messages from
// a Kafka topic.
type Consumer struct {
	conf *Config
}

// NewConsumer returns a new Consumer configured with the provided dialer and config.
func NewConsumer(opts ...Option) (*Consumer, error) {
	c := &Consumer{
		conf: newConfig(),
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.conf.mustProcess(); err != nil {
		return nil, errors.Errorf("bad consumer config: %w", err)
	}

	return c, nil
}

func (c *Consumer) Consume(ctx context.Context) error {
	group, err := newGroupConsumer(c.conf)
	if err != nil {
		return errors.Errorf("failed to create kafka consumer: %w", err)
	}

	return group.consume(ctx)
}
