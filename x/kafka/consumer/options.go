package consumer

import (
	"github.com/segmentio/kafka-go"
)

type Option func(consumer *Consumer)

// WithExplicitCommit enables offset commit only after a message is successfully
// handled.
//
// Do not use this option if the default behaviour of auto committing offsets on
// initial read (before handling the message) is required.
func WithExplicitCommit() Option {
	return func(consumer *Consumer) {
		consumer.withExplicitCommit = true
	}
}

// WithGroupBalancers adds a priority-ordered list of client-side consumer group
// balancing strategies that will be offered to the coordinator. The first strategy
// that all group members support will be chosen by the leader.
//
// Default: [Range, RoundRobin]
//
// Only used by consumer group.
func WithGroupBalancers(groupBalancers ...kafka.GroupBalancer) Option {
	return func(consumer *Consumer) {
		consumer.readerConfig.GroupBalancers = groupBalancers
	}
}

// WithHandlerBackOffRetry adds a back off retry policy on the consumer handler.
func WithHandlerBackOffRetry(backOffConstructor HandlerRetryBackOffConstructor) Option {
	return func(consumer *Consumer) {
		consumer.backOffConstructor = backOffConstructor
	}
}

// WithNotifyError adds the NotifyError function to the consumer for it to be invoked
// on each consumer handler error.
func WithNotifyError(notify NotifyError) Option {
	return func(consumer *Consumer) {
		consumer.notifyErr = notify
	}
}

// WithReaderLogger specifies a logger used to report internal consumer reader
// changes.
func WithReaderLogger(logger kafka.LoggerFunc) Option {
	return func(consumer *Consumer) {
		consumer.readerConfig.Logger = logger
	}
}

// WithReaderErrorLogger specifies a logger used to report internal consumer
// reader errors.
func WithReaderErrorLogger(logger kafka.LoggerFunc) Option {
	return func(consumer *Consumer) {
		consumer.readerConfig.ErrorLogger = logger
	}
}

// WithDataDogTracing adds Data Dog tracing to the consumer.
//
// A span is started each time a Kafka message is read and finished when the offset
// is committed. The consumer span can also be retrieved from within your handler
// using tracer.SpanFromContext.
func WithDataDogTracing() Option {
	return func(consumer *Consumer) {
		consumer.withDataDogTracing = true
	}
}

func WithMessageBatching(batchSize int, getOrderingKeyFn GetOrderingKey) Option {
	return func(consumer *Consumer) {
		consumer.batchSize = batchSize
		if getOrderingKeyFn != nil {
			consumer.getOrderingKeyFn = getOrderingKeyFn
		}
	}
}

// WithKafkaReader allows a custom reader to be injected into the Consumer/Group.
// Using this will ignore any other reader specific options passed in.
//
// It is highly recommended to not use this option unless injecting a mock reader
// implementation for testing.
func WithKafkaReader(readerFn func() Reader) Option {
	return func(consumer *Consumer) {
		consumer.reader = readerFn()
	}
}
