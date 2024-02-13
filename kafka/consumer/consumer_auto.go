package consumer

import (
	"context"
	"errors"
	"fmt"
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

// AutoConsumers is a type that maps "topic name" to "consumer"
type AutoConsumers map[string]*AutoConsumer

func newAutoConsumer(topic string) *AutoConsumer {
	var channel chan Message

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "test1,test2" // revisit with default values
	}
	log.Debug("kafka_auto_consumer").
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("brokers", brokers).
			Str("topic", topic),
		).Details("creating new auto consumer")

	cfg := Config{
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
	}

	autoBackOff := NonStopExponentialBackOff
	autoBalancers := []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}}
	autoNotify := func(ctx context.Context, err error, msg Message) {
		log.Error("auto_consumer_notif_error", err).
			WithSystemTracing().
			Properties(log.SubDoc().
				Str("topic", msg.Topic).
				Str("key", string(msg.Key)).
				Str("value", string(msg.Value)),
			).Details("error consuming message")
	}
	autoReaderLogger := func(s string, i ...interface{}) {
		msg := fmt.Sprintf(s, i...)
		log.Info("auto_consumer_reader").
			WithSystemTracing().
			Details(msg)
	}
	autoReaderErrorLogger := func(s string, i ...interface{}) {
		msg := fmt.Sprintf(s, i...)
		err := errors.New(msg)
		log.Error("auto_consumer_reader_error", err).
			WithSystemTracing().
			Details(msg)
	}

	testRunnerKafkaReader := func() Reader {
		return newTestRunnerReader(topic)
	}

	var consumer *Consumer
	if isTestMode() {
		// running inside a test, configure different options including a testRunnerKafkaReader
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
			WithNotifyError(autoNotify),
			WithReaderLogger(autoReaderLogger),
			WithReaderErrorLogger(autoReaderErrorLogger),
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
	log.Debug("kafka_auto_consumer_run").
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("topic", dc.topic),
		).Details("auto consumer started")

	if err := dc.consumer.Run(ctx, dc.handleRetrievedMessage); err != nil {
		log.Error("kafka_auto_consumer_run", err).
			WithSystemTracing().
			Properties(log.SubDoc().
				Str("topic", dc.topic),
			).Details("auto consumer run error")
	}

	close(dc.channel)
}

func (dc *AutoConsumer) handleRetrievedMessage(ctx context.Context, msg Message) error {
	log.Debug("kafka_auto_consumer_handle_retrieved_message").
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("consumer_id", msg.Metadata.ConsumerID).
			Str("topic", msg.Topic).
			Int("partition", msg.Partition).
			Int64("offset", msg.Offset).
			Str("key", string(msg.Key)).
			Str("value", string(msg.Value)),
		).Details("handle message")

	dc.channel <- msg
	return nil
}
