package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
)

type groupConsumer struct {
	conf   *Config
	client client
	group  sarama.ConsumerGroup
	logger sarama.StdLogger
}

func newGroupConsumer(client client, conf *Config) (*groupConsumer, error) {
	group, err := client.NewConsumerGroup(conf.brokers, conf.groupId, conf.saramaConfig)
	if err != nil {
		return nil, errors.Errorf("error creating consumer: %w", err)
	}

	return &groupConsumer{
		conf:   conf,
		client: client,
		group:  group,
		logger: conf.stdLogger,
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

			if errors.Is(err, errClosedMessageChannel) {
				return err
			}

			var target dispatchHandlerError
			if errors.As(err, &target) {
				// for any client dispatch errors, return if the conf was set to true otherwise, ignore.
				if gc.conf.returnOnClientDispatchError {
					return err
				}

				gc.logger.Printf("consumer group detected client dispatch failure: err='%s'. Trying to recover...", err)
			}

			gc.logger.Printf("consumer group detected unexpected error: err='%s'. Trying to recover...", err)
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
