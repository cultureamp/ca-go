package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/cultureamp/ca-go/kafka/consumer"
	"github.com/google/uuid"
)

func TestPackageWithMock(t *testing.T) {
	consumer.DefaultConsumers = new(packageConsumerExampleMock)

	ctx := context.Background()
	ch, stop := consumer.Consume(ctx, "topic-name")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	ok := true
	for ok {
		select {
		// Read the messages from the topic
		case msg, opened := <-ch:
			if opened {
				fmt.Printf("Channel open=%t, topic=%s message_offset=%v\n", opened, msg.Topic, msg.Offset)
			} else {
				// time to stop, channel is closed
				ok = false
			}
		}
	}
}

type packageConsumerExampleMock struct{}

func (c *packageConsumerExampleMock) Consume(ctx context.Context, topic string) (<-chan consumer.Message, consumer.StopFunc) {
	var channel chan consumer.Message

	channel = make(chan consumer.Message, 10)
	go c.run(channel)
	return channel, func() error { return nil }
}

func (c *packageConsumerExampleMock) run(channel chan consumer.Message) {
	defer close(channel)

	// generate some messages and then exit when done
	for i := 0; i < 10; i++ {
		m := c.newMessage()
		channel <- m
	}
}

func (c *packageConsumerExampleMock) newMessage() consumer.Message {
	msg := consumer.Message{}

	msg.Topic = "topic-name"
	msg.Offset = rand.Int63()           //#nosec G404 -- This for the test runner reader, not used in production.
	msg.Partition = rand.Intn(20-1) + 1 //#nosec G404 -- This for the test runner reader, not used in production.
	msg.Value = []byte(uuid.New().String())
	return msg
}
