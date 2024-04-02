package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type autoConsumer struct {
	topic    string
	consumer *Consumer
	channel  chan Message
}

// autoConsumers is a type that maps "topic name" to "consumer".
type autoConsumers map[string]*autoConsumer

func newAutoConsumer(topic string, brokers []string) *autoConsumer {
	var channel chan Message

	cfg := Config{
		Brokers: brokers,
		Topic:   topic,
	}

	autoBackOff := NonStopExponentialBackOff
	autoBalancers := []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}}

	var consumer *Consumer
	if isTestMode() {
		// running inside a test, configure different options including a testRunnerKafkaReader
		testRunnerKafkaReader := func() Reader {
			return newTestRunnerReader(topic)
		}
		consumer = NewConsumer(kafka.DefaultDialer, cfg,
			WithExplicitCommit(),
			WithGroupBalancers(autoBalancers...),
			WithHandlerBackOffRetry(autoBackOff),
			WithKafkaReader(testRunnerKafkaReader),
		)
	} else {
		consumer = NewConsumer(kafka.DefaultDialer, cfg,
			WithExplicitCommit(),
			WithGroupBalancers(autoBalancers...),
			WithHandlerBackOffRetry(autoBackOff),
			WithLogger(new(autoKafkaLogger)),
			WithNotifyError(autoClientNotifyError),
			WithDataDogTracing(),
		)
	}

	channel = make(chan Message, consumer.conf.QueueCapacity)

	return &autoConsumer{
		topic:    topic,
		consumer: consumer,
		channel:  channel,
	}
}

// Run method call blocks until the context is canceled, the consumer is closed, or an error occurs.
func (ac *autoConsumer) run(ctx context.Context) {
	defer close(ac.channel)

	if err := ac.consumer.Run(ctx, ac.handleMessage); err != nil {
		ac.consumer.conf.ErrorLogger.Printf(
			"auto consumer(%s:%s): err: %v",
			ac.topic,
			ac.consumer.id,
			err,
		)
	}
}

func (ac *autoConsumer) stop() error {
	return ac.consumer.Stop()
}

func (ac *autoConsumer) handleMessage(ctx context.Context, msg Message) error {
	ac.consumer.conf.Logger.Printf(
		"auto consumer(%s:%s): writing message to channel...",
		ac.topic,
		ac.consumer.id,
	)
	ac.channel <- msg
	return nil
}