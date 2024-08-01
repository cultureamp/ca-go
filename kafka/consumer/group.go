package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
)

type groupConsumer struct {
	conf       *Config
	client     kafkaClient
	dispatcher dispatcher
	decoder    decoder
	group      sarama.ConsumerGroup // Pete - this might make it too complex for others since it uses go routines
	logger     sarama.StdLogger
}

func newGroupConsumer(client kafkaClient, decoder decoder, dispatcher dispatcher, conf *Config) (*groupConsumer, error) {
	group, err := client.NewConsumerGroup(conf.brokers, conf.groupID, conf.saramaConfig)
	if err != nil {
		return nil, errors.Errorf("error creating consumer: %w", err)
	}

	return &groupConsumer{
		conf:       conf,
		client:     client,
		decoder:    decoder,
		dispatcher: dispatcher,
		group:      group,
		logger:     conf.stdLogger,
	}, nil
}

func (gc *groupConsumer) consume(ctx context.Context) error {
	// need to close() this groupConsumer or it will leak memory

	// `consume` should be called inside an infinite loop, when a
	// server-side rebalance happens, the consumer session will need to be
	// recreated to get the new claims
	for {
		handler := newHandler(gc.client, gc.decoder, gc.dispatcher, gc.conf.stdLogger)
		if err := gc.group.Consume(ctx, gc.conf.topics, handler); err != nil {
			if errFatal := gc.handleDispatchErrors(err); errFatal != nil {
				return errFatal
			}
		}

		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (gc *groupConsumer) handleDispatchErrors(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	if errors.Is(err, sarama.ErrClosedConsumerGroup) {
		return err
	}

	if errors.Is(err, errClosedMessageChannel) {
		return err
	}

	var target messageHandlerError
	if errors.As(err, &target) {
		// for any client dispatch errors, return if the conf was set to true otherwise, ignore.
		if gc.conf.returnOnClientDispatchError {
			return err
		}

		gc.logger.Printf("consumer group detected client dispatch failure: err='%s'. Trying to recover...", err)
		return nil
	}

	gc.logger.Printf("consumer group detected unexpected error: err='%s'. Trying to recover...", err)
	return nil
}

func (gc *groupConsumer) stop() error {
	// Cleans up memory inside sarama
	return gc.group.Close()
}
