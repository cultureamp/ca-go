package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConsumer(t *testing.T) {
	c, err := NewConsumer()
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing brokers")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing topics")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing group")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing message handler")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithHandler(func(ctx context.Context, msg *Message) error { return nil }),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithHandler(func(ctx context.Context, msg *Message) error { return nil }),
		WithAssignor("abc"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unrecognized consumer group partition assignor")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithHandler(func(ctx context.Context, msg *Message) error { return nil }),
		WithAssignor("roundrobin"),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithHandler(func(ctx context.Context, msg *Message) error { return nil }),
		WithAssignor("roundrobin"),
		WithVersion("abc"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid kafka version")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithHandler(func(ctx context.Context, msg *Message) error { return nil }),
		WithAssignor("roundrobin"),
		WithVersion("1.0.0"),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	// full list for coverage purposes
	c, err = NewConsumer(
		WithConsumerID("abc.123.uuid"),
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithHandler(func(ctx context.Context, msg *Message) error { return nil }),
		WithAssignor("roundrobin"),
		WithOldest(false),
		WithLogging(newTestLogger()),
		WithDebugLogger(newTestLogger()),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestConsumerCtxDeadLine(t *testing.T) {
	ctx := context.Background()
	// ctx will expire in 1 second
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancelCtx := context.WithDeadline(ctx, deadline)
	defer cancelCtx()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockClaim)

	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockClient.On("CommitMessage", mock.Anything, mock.Anything)
	mockSession.On("Context").Return(ctx)
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockClaim.On("Messages").Return(receiverChannel)

	handler := func(ctx context.Context, msg *Message) error {
		return nil
	}

	c := testConsumer(t, kafkaClient(mockClient), Handler(handler), int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.Consume(ctx)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}

func TestConsumerWithReceiverError(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockClaim)

	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockSession.On("Context").Return(ctx)
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockClaim.On("Messages").Return(receiverChannel)

	handler := func(ctx context.Context, msg *Message) error {
		return errors.Errorf("test error")
	}

	c := testConsumer(t, kafkaClient(mockClient), Handler(handler), int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.Consume(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "test error")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}

func TestConsumerWithDoubleConsumeAndStop(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockClaim)

	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockSession.On("Context").Return(ctx)
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockClaim.On("Messages").Return(receiverChannel)

	handler := func(ctx context.Context, msg *Message) error {
		return errors.Errorf("test error")
	}

	c := testConsumer(t, kafkaClient(mockClient), Handler(handler), int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.Consume(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "test error")

	err = c.Consume(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "consumer group already running!")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}


func testConsumer(t *testing.T, client kafkaClient, handler Handler, numMessages int64, ch chan *sarama.ConsumerMessage) *Consumer {
	// push a few messages into the channel
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
		ch <- saramaMessage
	}

	c, err := NewConsumer(
		WithKafkaClient(client),
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(handler),
		WithLogging(newTestLogger()),
		WithReturnOnClientDispathError(true),
	)
	assert.Nil(t, err)

	return c
}
