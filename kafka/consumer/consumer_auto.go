package consumer

import (
	"context"
	"os"
	"strings"

	"github.com/cultureamp/ca-go/log"
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
		autoNotify := func(ctx context.Context, err error, msg Message) {
			log.Error("auto_consumer_notify_error", err).
				WithSystemTracing().
				Properties(log.SubDoc().
					Str("topic", msg.Topic).
					Str("key", string(msg.Key)).
					Str("value", string(msg.Value)),
				).Details("error consuming message")
		}
		consumer = NewConsumer(kafka.DefaultDialer, cfg,
			WithExplicitCommit(),
			WithGroupBalancers(autoBalancers...),
			WithHandlerBackOffRetry(autoBackOff),
			WithKafkaLogger(new(autoClientLogger)),
			WithNotifyError(autoNotify),
			WithDataDogTracing(),
		)
	}

	channel = make(chan Message, consumer.readerConfig.QueueCapacity)

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
		log.Error("auto_consumer_run", err).
			WithSystemTracing().
			Properties(log.SubDoc().
				Str("topic", dc.topic),
			).Details("auto consumer run error")
	}
}

func (dc *AutoConsumer) handleRetrievedMessage(ctx context.Context, msg Message) error {
	log.Debug("auto_consumer_handle_message").
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("consumer_id", msg.Metadata.ConsumerID).
			Str("topic", msg.Topic).
			Int("partition", msg.Partition).
			Int64("offset", msg.Offset).
			Str("key", string(msg.Key)).
			Str("value", string(msg.Value)),
		).Details("writing message to channel...")

	dc.channel <- msg
	return nil
}
