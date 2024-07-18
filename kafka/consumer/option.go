package consumer

import (
	"github.com/IBM/sarama"
)

type Option func(consumer *Subscriber)

func WithKafkaClient(client client) Option {
	return func(consumer *Subscriber) {
		consumer.client = client
	}
}

func WithAvroDecoder(decoder decoder) Option {
	return func(consumer *Subscriber) {
		consumer.decoder = decoder
	}
}

func WithHandler(handler Receiver) Option {
	return func(consumer *Subscriber) {
		consumer.receiver = handler
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

func WithSchemaRegistryURL(schemaRegistryURL string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.schemaRegistryURL = schemaRegistryURL
	}
}

func WithGroupID(id string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.groupID = id
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
