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
	group, err := client.NewConsumerGroup(conf.brokers, conf.groupId, conf.saramaConfig)
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
	// need to close() this groupConsumer or it will leak memory

	// `consume` should be called inside an infinite loop, when a
	// server-side rebalance happens, the consumer session will need to be
	// recreated to get the new claims
	for {
		receiver := newConsumer(gc.client, gc.conf.handler, gc.conf.stdLogger)
		if err := gc.group.Consume(ctx, gc.conf.topics, receiver); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return err
			}
			// for any client dispatch errors, return is the conf was set to true
			// otherwise, ignore
			if gc.conf.returnOnClientDispatchError {
				return err
			}
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (gc *groupConsumer) stop() error {
	// Cleans up memory inside sarama
	return gc.group.Close()
}
