package consumer

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cultureamp/ca-go/log"
	"github.com/segmentio/kafka-go"
)

type AutoConsumer struct {
	topic   string
	client  ConsumerClient
	channel chan Message
}

type AutoConsumers map[string]*AutoConsumer

var DefaultKafkaConsumers = getInstance()

func getInstance() AutoConsumers {
	auto := make(AutoConsumers)
	return auto
}

func Run(topic string) <-chan Message {
	c, found := DefaultKafkaConsumers[topic]
	if !found {
		c = newAutoConsumer(topic)
		DefaultKafkaConsumers[topic] = c
		go c.run()
	}
	return c.channel
}

func newAutoConsumer(topic string) *AutoConsumer {
	var consumer ConsumerClient
	var channel chan Message

	if isTestMode() {
		consumer = newTestRunnerClient(topic)
		channel = make(chan Message, 1)
	} else {
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

		wantBackOff := NonStopExponentialBackOff
		wantBalancers := []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}}
		wantNotify := func(ctx context.Context, err error, msg Message) {
			log.Error("auto_consumer_notif_error", err).
				WithSystemTracing().
				Properties(log.SubDoc().
					Str("topic", msg.Topic).
					Str("key", string(msg.Key)).
					Str("value", string(msg.Value)),
				).Details("error consuming message")
		}
		wantReaderLogger := func(s string, i ...interface{}) {
			msg := fmt.Sprintf(s, i...)
			log.Info("auto_consumer_reader").
				WithSystemTracing().
				Details(msg)
		}
		wantReaderErrorLogger := func(s string, i ...interface{}) {
			msg := fmt.Sprintf(s, i...)
			err := errors.New(msg)
			log.Error("auto_consumer_reader_error", err).
				WithSystemTracing().
				Details(msg)
		}

		consumer = NewConsumer(kafka.DefaultDialer, cfg,
			WithExplicitCommit(),
			WithGroupBalancers(wantBalancers...),
			WithHandlerBackOffRetry(wantBackOff),
			WithNotifyError(wantNotify),
			WithReaderLogger(wantReaderLogger),
			WithReaderErrorLogger(wantReaderErrorLogger),
			WithDataDogTracing(),
		)
		channel = make(chan Message, 100)
	}

	dc := &AutoConsumer{
		topic:   topic,
		client:  consumer,
		channel: channel,
	}

	return dc
}

func (dc *AutoConsumer) run() {
	log.Debug("kafka_auto_consumer_run").
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("topic", dc.topic),
		).Details("auto consumer started")
	if err := dc.client.Run(context.Background(), dc.handle); err != nil {
		panic(err)
	}

	close(dc.channel)
}

func (dc *AutoConsumer) handle(ctx context.Context, msg Message) error {
	log.Debug("kafka_auto_consumer_handle").
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

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}
