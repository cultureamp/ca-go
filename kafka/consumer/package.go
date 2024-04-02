package consumer

import (
	"context"
	"flag"
	"os"
	"strings"
)

var DefaultConsumers = getInstance()

func getInstance() TopicConsumer {
	def := newDefaultConsumer()
	return def
}

type StopFunc func() error

type TopicConsumer interface {
	Consume(ctx context.Context, topic string) (<-chan Message, StopFunc)
}

// Consume creates as new channel for the specified topic and starts sending Messages to it.
func Consume(ctx context.Context, topic string) (<-chan Message, StopFunc) {
	ch, stop := DefaultConsumers.Consume(ctx, topic)
	return ch, stop
}

type topicConsumer struct {
	brokers       []string
	autoConsumers autoConsumers
}

func newDefaultConsumer() *topicConsumer {
	topicConsumers := make(autoConsumers)

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "test1,test2" // revisit with default values
	}

	return &topicConsumer{
		brokers:       strings.Split(brokers, ","),
		autoConsumers: topicConsumers,
	}
}

// Consume creates as new consumer for the specified topic.
func (tc *topicConsumer) Consume(ctx context.Context, topic string) (<-chan Message, StopFunc) {
	ac := newAutoConsumer(topic, tc.brokers)
	tc.autoConsumers[topic] = ac
	go ac.run(ctx)
	return ac.channel, func() error { return tc.stop(topic) }
}

// stop sends a signal to the consumer to finish, returning an error if it failed to do so.
func (tc *topicConsumer) stop(topic string) error {
	ac, found := tc.autoConsumers[topic]
	if found {
		delete(tc.autoConsumers, topic)
		return ac.stop()
	}
	return nil
}

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		strings.Contains(argZero, "__debug_bin") || // vscode debug binary
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}
