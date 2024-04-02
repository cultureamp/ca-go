package consumer

import (
	"context"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

type AutoConsumer struct {
	topic    string
	consumer *Consumer
	channel  chan Message
}

// AutoConsumers is a type that maps "topic name" to "consumer".
type AutoConsumers map[string]*AutoConsumer

func newAutoConsumer(topic string) *AutoConsumer {
	var channel chan Message

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "test1,test2" // revisit with default values
	}
	cfg := Config{
		Brokers: strings.Split(brokers, ","),
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

	dc := &AutoConsumer{
		topic:    topic,
		consumer: consumer,
		channel:  channel,
	}

	return dc
}

func (dc *AutoConsumer) run(ctx context.Context) {
	defer close(dc.channel)

	if err := dc.consumer.Run(ctx, dc.handleRetrievedMessage); err != nil {
		dc.consumer.conf.ErrorLogger.Printf(
			"auto consumer(%s:%s): err: %v",
			dc.topic,
			dc.consumer.id,
			err,
		)
	}
}

func (dc *AutoConsumer) handleRetrievedMessage(ctx context.Context, msg Message) error {
	dc.consumer.conf.Logger.Printf(
		"auto consumer(%s:%s): writing message to channel...",
		dc.topic,
		dc.consumer.id,
	)
	dc.channel <- msg
	return nil
}
