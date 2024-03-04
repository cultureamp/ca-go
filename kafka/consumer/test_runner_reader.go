package consumer

import (
	"context"
	"math/rand"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type testRunnerReader struct {
	topic string
}

func newTestRunnerReader(topic string) *testRunnerReader {
	return &testRunnerReader{
		topic: topic,
	}
}

func (trr *testRunnerReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg := trr.newMessage()
	return msg, nil
}

func (trr *testRunnerReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	msg := trr.newMessage()
	return msg, nil
}

func (trr *testRunnerReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	return nil
}

func (trr *testRunnerReader) Close() error {
	return nil
}

func (trr *testRunnerReader) newMessage() kafka.Message {
	msg := kafka.Message{}

	msg.Topic = trr.topic
	msg.Offset = rand.Int63()           //#nosec G404 -- This for the test runner reader, not used in production.
	msg.Partition = rand.Intn(20-1) + 1 //#nosec G404 -- This for the test runner reader, not used in production.
	msg.Value = []byte(uuid.New().String())
	return msg
}
