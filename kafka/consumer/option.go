package consumer

import (
	"time"

	"github.com/IBM/sarama"
)

type Option func(consumer *Subscriber)

// WithKafkaClient sets the consumer kafka client (Default: saramaClient).
func WithKafkaClient(client kafkaClient) Option {
	return func(consumer *Subscriber) {
		consumer.kakfaClient = client
	}
}

// WithAvroSchemaRegistryClient sets the consumer avro schema registry client
// (Default: avroSchemaRegistryClient "github.com/riferrei/srclient").
func WithAvroSchemaRegistryClient(client schemaRegistryClient) Option {
	return func(consumer *Subscriber) {
		consumer.avroClient = client
	}
}

// WithAvroDecoder sets the consumer avro decoder (Default: avroDecoder).
func WithAvroDecoder(decoder decoder) Option {
	return func(consumer *Subscriber) {
		consumer.decoder = decoder
	}
}

// WithHandler sets the consumer message handler that clients MUST implement.
func WithHandler(handler Receiver) Option {
	return func(consumer *Subscriber) {
		consumer.receiver = handler
	}
}

// WithBrokers sets the consumer brokers (Default: env var 'KAFKA_BROKERS').
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

// WithTopics sets the consumer topics (Default: env var 'KAFKA_TOPICS').
func WithTopics(topics []string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.topics = topics
	}
}

// WithOldest sets the consumer initial offset from oldest (Default: true).
func WithOldest(oldest bool) Option {
	return func(consumer *Subscriber) {
		consumer.conf.oldest = oldest
	}
}

// WithLogger sets the consumer logger (Default: sarama.StdLogger).
func WithLogging(logger sarama.StdLogger) Option {
	return func(consumer *Subscriber) {
		consumer.conf.stdLogger = logger
	}
}

// WithDebugLogger sets the consumer debug logger (Default: sarama.DebugLogger).
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

// WithSchemaRegistryURL sets the client avro registry URL (Default: env var 'SCHEMA_REGISTRY_URL').
func WithSchemaRegistryURL(schemaRegistryURL string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.schemaRegistryURL = schemaRegistryURL
	}
}

// WithGroupID sets the consumer groupID.
func WithGroupID(id string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.groupID = id
	}
}

// WithAssignor sets the consumer group partition assignor (Default: "range". Other options: "sticky" or "roundrobin").
func WithAssignor(assignor string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.assignor = assignor
	}
}

// WithReturnOnClientDispathError sets whether the consume should exit on receiver.dispatch error (Default false).
func WithReturnOnClientDispathError(returnOnError bool) Option {
	return func(consumer *Subscriber) {
		consumer.conf.returnOnClientDispatchError = returnOnError
	}
}

// WithChannelSize sets the consumer channel size (Default: 256).
func WithChannelMessageSize(size int) Option {
	return func(consumer *Subscriber) {
		consumer.conf.saramaConfig.ChannelBufferSize = size
	}
}

// WithMaxWaitTime sets the consumer max wait time (Default: 500ms).
func WithMaxWaitTime(waitTime time.Duration) Option {
	return func(consumer *Subscriber) {
		consumer.conf.saramaConfig.Consumer.MaxWaitTime = waitTime
	}
}

// WithFetchSize sets the consumer fetch size (Default: 10MB).
func WithFetchSize(fetchSize int) Option {
	return func(consumer *Subscriber) {
		consumer.conf.saramaConfig.Consumer.Fetch.Default = int32(fetchSize)
	}
}

// WithSaslUsername sets the consumer SASL username.
func WithSaslUsername(username string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.saramaConfig.Net.SASL.User = username
	}
}

func WithSaslPassword(password string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.saramaConfig.Net.SASL.Password = password
	}
}

func WithSaslAlgorithm(algo string) Option {
	return func(consumer *Subscriber) {
		consumer.conf.algorithm = algo
	}
}
