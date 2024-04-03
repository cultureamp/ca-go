package consumer_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/kafka/consumer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPackageWithTimeout(t *testing.T) {
	defConsumers := consumer.DefaultConsumers
	consumer.DefaultConsumers = newPackageConsumerExampleMock(0)
	defer func() {
		consumer.DefaultConsumers = defConsumers
	}()

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
		assert.Fail(t, "should not have recieved any messages!")

	case <-time.After(time.Millisecond * 100):
		fmt.Println("Timeout recieved. Stopping.")
	}
}

func TestPackageWithDeadline(t *testing.T) {
	deadline := time.Millisecond * 100
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	defConsumers := consumer.DefaultConsumers
	consumer.DefaultConsumers = newPackageConsumerExampleMock(0)
	defer func() {
		consumer.DefaultConsumers = defConsumers
	}()

	ch, stop := consumer.Consume(ctx, "topic-name-with-context-deadline")
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
		assert.Fail(t, "should not have recieved any messages!")

	case <-ctx.Done():
		fmt.Println("Context deadline received. Stopping.")
	}
}

func TestPackageWithMock(t *testing.T) {
	const queueSize = 5

	defConsumers := consumer.DefaultConsumers
	consumer.DefaultConsumers = newPackageConsumerExampleMock(queueSize)
	defer func() {
		consumer.DefaultConsumers = defConsumers
	}()

	ctx := context.Background()
	ch, stop := consumer.Consume(ctx, "topic-name")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	// Read ALL the messages from the channel until it is closed
	i := 0
	for range ch {
		i++
	}
	assert.Equal(t, queueSize, i)
}

type packageConsumerExampleMock struct {
	queueSize int
}

func newPackageConsumerExampleMock(qs int) *packageConsumerExampleMock {
	return &packageConsumerExampleMock{
		queueSize: qs,
	}
}

func (c *packageConsumerExampleMock) Consume(ctx context.Context, topic string) (<-chan consumer.Message, consumer.StopFunc) {
	var channel chan consumer.Message

	channel = make(chan consumer.Message, c.queueSize)
	go c.run(channel)
	return channel, func() error { return nil }
}

func (c *packageConsumerExampleMock) run(channel chan consumer.Message) {
	defer close(channel)

	if c.queueSize == 0 {
		time.Sleep(time.Millisecond * 500)
	}

	// generate some messages and then exit when done
	for i := 0; i < c.queueSize; i++ {
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
