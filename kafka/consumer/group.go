package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
)

type groupConsumer struct {
	conf        *Config
	client      kafkaClient
	groupClient sarama.ConsumerGroup
}

func newGroupConsumer(client kafkaClient, conf *Config) (*groupConsumer, error) {
	group, err := client.NewConsumerGroup(conf.brokers, conf.groupId, conf.saramaConfig)
	if err != nil {
		return nil, errors.Errorf("error creating consumer: %w", err)
	}

	return &groupConsumer{
		conf:        conf,
		client:      client,
		groupClient: group,
	}, nil
}

func (gc *groupConsumer) consume(ctx context.Context) error {
	// todo need to close this groupConsumer, is this defer ok?
	defer gc.groupClient.Close()

	// `consume` should be called inside an infinite loop, when a
	// server-side rebalance happens, the consumer session will need to be
	// recreated to get the new claims
	for {
		receiver := newReceiver(gc.client, gc.conf.handler)
		if err := gc.groupClient.Consume(ctx, gc.conf.topics, receiver); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return err
			}
			if gc.conf.returnOnError {
				sarama.Logger.Printf("error from client. Exiting consume: '%s'", err.Error())
				return err
			}
			sarama.Logger.Printf("error from client: '%s'", err.Error())
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
