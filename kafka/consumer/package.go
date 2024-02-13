package consumer

import (
	"context"
	"flag"
	"os"
	"strings"
)

var TopicConsumers = getInstance()

func getInstance() AutoConsumers {
	auto := make(AutoConsumers)
	return auto
}

// Consume reads messages from the topic until there is an error
// or if the ctx deadline is reached.
func Consume(ctx context.Context, topic string) <-chan Message {
	c := newAutoConsumer(topic)
	TopicConsumers[topic] = c
	go c.run(ctx)
	return c.channel
}

// Stop sends a signal to the consumer to finish returning an error if it failed to do so.
func Stop(topic string) error {
	c, found := TopicConsumers[topic]
	if found {
		delete(TopicConsumers, topic)
		return c.consumer.Stop()
	}
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
