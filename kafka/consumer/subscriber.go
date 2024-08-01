package consumer

import (
	"context"
	"sync"

	"github.com/go-errors/errors"
)

// Subscriber provides a high level API for consuming and handling messages from a Kafka topic.
// This implementation blocks on ConsumeAll() if you want a non-blocking version use Service.
type Subscriber struct {
	conf *Config

	kakfaClient kafkaClient          // Kafka client (Default: Sarama)
	avroClient  schemaRegistryClient // Avro client (Default: srclient)
	decoder     decoder
	receiver    Receiver
	groupMutex  sync.Mutex
	group       *groupConsumer
}

// NewSubscriber returns a new Subscriber configured with the provided options.
// Note: The receiver MUST be set using WithHandler() option.
func NewSubscriber(opts ...Option) (*Subscriber, error) {
	c := &Subscriber{
		conf:        newConfig(),
		kakfaClient: newSaramaClient(),
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.conf.shouldProcess(); err != nil {
		return nil, errors.Errorf("bad consumer config: %w", err)
	}

	if c.receiver == nil {
		return nil, errors.Errorf("missing message handler")
	}

	if c.avroClient == nil {
		c.avroClient = newAvroSchemaRegistryClient(c.conf.schemaRegistryURL)
	}

	if c.decoder == nil {
		c.decoder = newAvroDecoder(c.avroClient)
	}

	return c, nil
}

// ConsumeAll consumes all the messages from the configured Kafka topic.
// Note: The is a blocking call until the context is done, cancelled, deadline reached or an error occurs.
func (c *Subscriber) ConsumeAll(ctx context.Context) error {
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

	dispatcher := newDispatcher(c.receiver)
	group, err := newGroupConsumer(c.kakfaClient, c.decoder, dispatcher, c.conf)
	if err != nil {
		return nil, errors.Errorf("failed to create kafka consumer: %w", err)
	}

	c.group = group
	return group, nil
}

// Stop terminates the subscriber and closes the underlying Kafka consumer.
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
