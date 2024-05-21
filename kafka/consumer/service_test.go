package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceWithCancelledContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	calls := 0
	handler := func(ctx context.Context, msg *Message) error {
		calls++
		return nil
	}

	c := testConsumer(t, ctx, handler, 4)
	assert.NotNil(t, c)

	s := NewService(c)
	assert.NotNil(t, s)

	// non blocking
	s.Start(ctx)
	defer s.Stop()

	// sleep this thread for a bit to let the service do some work
	time.Sleep(1 * time.Second)

	// cancel the context to signal the service to stop
	cancel()
	assert.Equal(t, 4, calls)
}

func TestServiceWithStop(t *testing.T) {
	ctx := context.Background()

	calls := 0
	handler := func(ctx context.Context, msg *Message) error {
		calls++
		return nil
	}

	c := testConsumer(t, ctx, handler, 2)
	assert.NotNil(t, c)

	s := NewService(c)
	assert.NotNil(t, s)

	// non blocking
	s.Start(ctx)

	// sleep this thread for a bit to let the service do some work
	time.Sleep(1 * time.Second)

	s.Stop()
	assert.Equal(t, 2, calls)
}

func testConsumer(t *testing.T, ctx context.Context, handler Handler, numMessages int64) *Consumer {
	mockClient := newMockKafkaClient()
	mockConsumerGroupSession := newMockConsumerGroupSession()
	mockConsumerGroupClaim := newMockConsumerGroupClaim()
	mockConsumerGroup := newMockConsumerGroup(mockConsumerGroupSession, mockConsumerGroupClaim)

	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockConsumerGroup, nil)
	mockClient.On("CommitMessage", mock.Anything, mock.Anything)
	mockConsumerGroupSession.On("Context").Return(ctx)
	mockConsumerGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockConsumerGroup.On("Close").Return(nil)

	// push a few messages into the channel
	mockChannel := make(chan *sarama.ConsumerMessage, 10)

	for i := range numMessages {
		saramaMessage := &sarama.ConsumerMessage{
			Topic:     "test",
			Partition: 1,
			Key:       []byte("key"),
			Value:     []byte("value"),
			Offset:    i,
			Timestamp: time.Now(),
			Headers:   nil,
		}
		mockChannel <- saramaMessage
	}

	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockConsumerGroupClaim.On("Messages").Return(receiverChannel)

	c, err := NewConsumer(
		WithKafkaClient(mockClient),
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(handler),
		WithLogging(newTestLogger()),
		WithReturnOnError(true),
	)
	assert.Nil(t, err)

	return c
}
