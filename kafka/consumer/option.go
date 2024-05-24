package consumer

import (
	"github.com/IBM/sarama"
)

type Option func(consumer *Subscriber)

func WithKafkaClient(client kafkaClient) Option {
	return func(consumer *Subscriber) {
		consumer.client = client
	}
}

func WithBrokers(brokers []string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.brokers = brokers
	}
}

// WithVersion sets the underlying Sarama version (Default: V2_1_0_0).
func WithVersion(version string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.version = version
	}
}

func WithTopics(topics []string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.topics = topics
	}
}

// WithOldest sets the consumer initial offset from oldest (Default true).
func WithOldest(oldest bool) Option {
	return func(consumer *Subscriber) {
		consumer.conf.oldest = oldest
	}
}

func WithLogging(logger sarama.StdLogger) Option {
	return func(consumer *Subscriber) {
		consumer.conf.stdLogger = logger
	}
}

func WithDebugLogger(logger sarama.StdLogger) Option {
	return func(consumer *Subscriber) {
		consumer.conf.debugLogger = logger
	}
}

// WithConsumerID sets the consumer id (Default new uiid).
func WithConsumerID(id string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.id = id
	}
}

func WithHandler(handler Handler) Option {
	return func(consumer *Subscriber) {
		consumer.conf.handler = handler
	}
}

func WithGroupId(id string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.groupId = id
	}
}

func WithAssignor(assignor string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.assignor = assignor
	}
}

func WithReturnOnClientDispathError(returnOnError bool) Option {
	return func(consumer *Subscriber) {
		consumer.conf.returnOnClientDispatchError = returnOnError
	}
}
