package consumer

import (
	"github.com/IBM/sarama"
)

type Option func(consumer *Consumer)

func WithBrokers(brokers []string) Option {
	return func(consumer *Consumer) {
		consumer.conf.brokers = brokers
	}
}

// WithVersion sets the underlying Sarama version (Default: V2_1_0_0).
func WithVersion(version string) Option {
	return func(consumer *Consumer) {
		consumer.conf.version = version
	}
}

func WithTopics(topic string) Option {
	return func(consumer *Consumer) {
		consumer.conf.topic = topic
	}
}

// WithOldest sets the consumer initial offset from oldest (Default true).
func WithOldest(oldest bool) Option {
	return func(consumer *Consumer) {
		consumer.conf.oldest = oldest
	}
}

func WithLogging(logger sarama.StdLogger) Option {
	return func(consumer *Consumer) {
		consumer.conf.logger = logger
	}
}

// WithClientID sets the consumer id (Default new uiid).
func WithClientID(id string) Option {
	return func(consumer *Consumer) {
		consumer.conf.id = id
	}
}

func WithHandler(handler Handler) Option {
	return func(consumer *Consumer) {
		consumer.conf.handler = handler
	}
}
