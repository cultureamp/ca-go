package consumer

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceWithCancelledContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var calls atomic.Int32
	handler := func(ctx context.Context, msg *Message) error {
		calls.Add(1)
		return nil
	}

	s := testService(t, ctx, handler, 4)
	assert.NotNil(t, s)

	// non blocking
	s.Start(ctx)
	defer s.Stop()

	// sleep this thread for a bit to let the service do some work
	time.Sleep(1 * time.Second)

	// cancel the context to signal the service to stop
	cancel()
	assert.Equal(t, int32(4), calls.Load())
}

func TestServiceWithStop(t *testing.T) {
	ctx := context.Background()

	var calls atomic.Int32
	handler := func(ctx context.Context, msg *Message) error {
		calls.Add(1)
		return nil
	}

	s := testService(t, ctx, handler, 2)
	assert.NotNil(t, s)

	// non blocking
	s.Start(ctx)

	// sleep this thread for a bit to let the service do some work
	time.Sleep(1 * time.Second)

	// stop the Service
	err := s.Stop()
	assert.Nil(t, err)
	assert.Equal(t, int32(2), calls.Load())
}

func TestServiceWithClientError(t *testing.T) {
	ctx := context.Background()

	var calls atomic.Int32
	handler := func(ctx context.Context, msg *Message) error {
		calls.Add(1)
		return errors.Errorf("test error")
	}

	s := testService(t, ctx, handler, 3)
	assert.NotNil(t, s)

	// non blocking
	s.Start(ctx)

	// sleep this thread for a bit to let the service do some work
	time.Sleep(1 * time.Second)

	// stop the Service
	err := s.Stop()
	assert.Nil(t, err)
	assert.Equal(t, int32(3), calls.Load())
}

func TestServiceWithDoubleStartDoubleStop(t *testing.T) {
	ctx := context.Background()

	var calls atomic.Int32
	handler := func(ctx context.Context, msg *Message) error {
		calls.Add(1)
		return errors.Errorf("test error")
	}

	s := testService(t, ctx, handler, 3)
	assert.NotNil(t, s)

	// non blocking
	s.Start(ctx)
	s.Start(ctx)

	// sleep this thread for a bit to let the service do some work
	time.Sleep(1 * time.Second)

	// stop the Service
	err := s.Stop()
	assert.Nil(t, err)
	assert.Equal(t, int32(3), calls.Load())
	err = s.Stop()
	assert.Nil(t, err)
}

func testService(t *testing.T, ctx context.Context, handler Handler, numMessages int64) *Service {
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

	s, err := NewService(
		WithKafkaClient(mockClient),
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(handler),
		WithLogging(newTestLogger()),
		WithReturnOnClientDispathError(false),
	)
	assert.Nil(t, err)

	return s
}
