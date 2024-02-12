package consumer

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
)

type ConsumerClient interface {
	Run(ctx context.Context, handler Handler) error
	Stop() error
}

type testRunnerClient struct {
	topic  string
	stopCh chan struct{}
}

func newTestRunnerClient(topic string) *testRunnerClient {
	return &testRunnerClient{
		topic:  topic,
		stopCh: make(chan struct{}),
	}
}

func (c *testRunnerClient) Run(ctx context.Context, handler Handler) error {
	i := int64(1)
	for {
		select {
		case <-c.stopCh:
			return nil
		default:
			i++
			msg := c.newMessage(i)
			err := handler(ctx, msg)
			if err != nil {
				return err
			}
		}
	}
}

func (c *testRunnerClient) Stop() error {
	close(c.stopCh)
	return nil
}

func (c *testRunnerClient) newMessage(i int64) Message {
	msg := Message{}

	msg.Topic = c.topic
	msg.Attempt = 1
	msg.Offset = i
	msg.ConsumerID = uuid.New().String()
	msg.Partition = rand.Intn(20-1) + 1 //#nosec G404 -- This for the test runner client, not used in production.
	msg.Value = []byte(uuid.New().String())
	return msg
}
