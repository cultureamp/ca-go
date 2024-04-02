package consumer_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/kafka/consumer"
	"github.com/google/uuid"
)

func TestPackage(t *testing.T) {
	ctx := context.Background()
	ch, stop := consumer.Consume(ctx, "topic-name")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	// Read the next message from the topic
	msg, ok := <-ch
	fmt.Printf("Channel open=%t, topic=%s message_offset=%v\n", ok, msg.Topic, msg.Offset)
}

func TestPackageWithTimeout(t *testing.T) {
	ctx := context.Background()
	ch, stop := consumer.Consume(ctx, "topic-name-with-timeout")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	select {
	case msg, ok := <-ch:
		fmt.Printf("Channel open=%t, topic=%s message_offset=%v\n", ok, msg.Topic, msg.Offset)

	case <-time.After(time.Duration(1) * time.Second):
		fmt.Println("No message received before timeout")
	}
}

func TestPackageWithDeadline(t *testing.T) {
	deadline := time.Duration(1) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	ch, stop := consumer.Consume(ctx, "topic-name-with-context-deadline")
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
		case msg, ok := <-ch:
			fmt.Printf("Channel open=%t, topic=%s message_offset=%v\n", ok, msg.Topic, msg.Offset)

		case <-ctx.Done():
			fmt.Println("Context deadline received. Stopping.")
			ok = false
		}
	}
}

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
