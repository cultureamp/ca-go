package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
)

type groupConsumer struct {
	conf   *Config
	client kafkaClient
	group  sarama.ConsumerGroup
}

func newGroupConsumer(client kafkaClient, conf *Config) (*groupConsumer, error) {
	group, err := client.newConsumerGroup(conf.brokers, conf.groupId, conf.saramaConfig)
	if err != nil {
		return nil, errors.Errorf("error creating consumer: %w", err)
	}

	return &groupConsumer{
		conf:   conf,
		client: client,
		group:  group,
	}, nil
}

func (gc *groupConsumer) consume(ctx context.Context) error {
	// `Consume` should be called inside an infinite loop, when a
	// server-side rebalance happens, the consumer session will need to be
	// recreated to get the new claims
	for {
		receiver := newReceiver(gc.client, gc.conf.handler)
		if err := gc.group.Consume(ctx, gc.conf.topics, receiver); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return err
			}
			sarama.Logger.Printf("error from consumer: %w", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
