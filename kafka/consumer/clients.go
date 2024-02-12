package consumer

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
)

type ConsumerClient interface {
	Run(ctx context.Context, handler Handler) error
}

type testRunnerClient struct {
	topic string
}

func newTestRunnerClient(topic string) *testRunnerClient {
	return &testRunnerClient{
		topic: topic,
	}
}

func (c *testRunnerClient) Run(ctx context.Context, handler Handler) error {
	for i := int64(1); i < 100; i++ {
		msg := c.newMessage(i)
		handler(ctx, msg)
	}
	return nil
}

func (c *testRunnerClient) newMessage(i int64) Message {
	msg := Message{}

	msg.Topic = c.topic
	msg.Attempt = 1
	msg.Offset = i
	msg.ConsumerID = uuid.New().String()
	msg.Partition = rand.Intn(20-1) + 1 //nolint:gosec
	msg.Value = []byte(uuid.New().String())
	return msg
}
