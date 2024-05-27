package consumer

import (
	"context"
	"sync"

	"github.com/go-errors/errors"
)

// Subscriber provides a high level API for consuming and handling messages from a Kafka topic.
// This implementation blocks on Consume() if you want a non-blocking version use Service.
type Subscriber struct {
	client client // Kafka client interfaces (Default: Sarama)
	conf   *Config

	groupMutex sync.Mutex
	group      *groupConsumer
}

// NewSubscriber returns a new Subscriber configured with the provided dialer and config.
func NewSubscriber(opts ...Option) (*Subscriber, error) {
	c := &Subscriber{
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

func (c *Subscriber) Subscribe(ctx context.Context) error {
	group, err := c.setupGroupConsumer()
	if err != nil {
		return err
	}

	// blocking call until either
	// 1. context is cancelled/done OR
	// 2. a server-side kafka rebalance happens OR
	// 3. client dispatch error occurs (and returnOnClientDispatchError=true in the)
	return group.consume(ctx)
}

func (c *Subscriber) setupGroupConsumer() (*groupConsumer, error) {
	c.groupMutex.Lock()
	defer c.groupMutex.Unlock()

	// if already consuming, do nothing
	if c.group != nil {
		return nil, errors.Errorf("consumer group already running! (forgot to call Stop()?)")
	}

	group, err := newGroupConsumer(c.client, c.conf)
	if err != nil {
		return nil, errors.Errorf("failed to create kafka consumer: %w", err)
	}

	c.group = group
	return group, nil
}

func (c *Subscriber) Stop() error {
	c.groupMutex.Lock()
	defer c.groupMutex.Unlock()

	// if already stopped, do nothing
	if c.group == nil {
		return nil
	}

	err := c.group.stop()
	c.group = nil

	return err
}
