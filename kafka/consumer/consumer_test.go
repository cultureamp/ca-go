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
		WithBrokers([]string{"localhost:9001"}),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing topics")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing group")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unrecognized consumer group partition assignor")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("abc"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unrecognized consumer group partition assignor")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing message handler")

	c, err = NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(func(ctx context.Context, msg *ConsumerMessage) error { return nil }),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestNewConsumerCtxDeadLine(t *testing.T) {
	ctx := context.Background()

	// ctx will expire in 1 second
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancelCtx := context.WithDeadline(ctx, deadline)
	defer cancelCtx()

	mockClient := newMockKafkaClient()
	mockConsumerGroupSession := newMockConsumerGroupSession()
	mockConsumerGroupClaim := newMockConsumerGroupClaim()
	mockConsumerGroup := newMockConsumerGroup(mockConsumerGroupSession, mockConsumerGroupClaim)

	mockClient.On("newConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockConsumerGroup, nil)
	mockClient.On("commitMessage", mock.Anything, mock.Anything)
	mockConsumerGroupSession.On("Context").Return(ctx)
	mockConsumerGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockConsumerGroup.On("Close").Return(nil)

	saramaMessage := &sarama.ConsumerMessage{
		Topic:     "test",
		Partition: 1,
		Key:       []byte("key"),
		Value:     []byte("value"),
		Offset:    1,
		Timestamp: time.Now(),
		Headers:   nil,
	}

	// push one message into the channel
	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	mockChannel <- saramaMessage
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockConsumerGroupClaim.On("Messages").Return(receiverChannel)

	c, err := NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(func(ctx context.Context, msg *ConsumerMessage) error { return nil }),
		WithKafkaClient(mockClient),
		WithLogging(newTestLogger()),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	err = c.Consume(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "context deadline exceeded")

	mockClient.AssertExpectations(t)
	mockConsumerGroup.AssertExpectations(t)
}

func TestNewConsumerWithReceiverError(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockConsumerGroupSession := newMockConsumerGroupSession()
	mockConsumerGroupClaim := newMockConsumerGroupClaim()
	mockConsumerGroup := newMockConsumerGroup(mockConsumerGroupSession, mockConsumerGroupClaim)

	mockClient.On("newConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockConsumerGroup, nil)
	mockConsumerGroupSession.On("Context").Return(ctx)
	mockConsumerGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockConsumerGroup.On("Close").Return(nil)

	saramaMessage := &sarama.ConsumerMessage{
		Topic:     "test",
		Partition: 1,
		Key:       []byte("key"),
		Value:     []byte("value"),
		Offset:    1,
		Timestamp: time.Now(),
		Headers:   nil,
	}

	// push one message into the channel
	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	mockChannel <- saramaMessage
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockConsumerGroupClaim.On("Messages").Return(receiverChannel)

	c, err := NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(func(ctx context.Context, msg *ConsumerMessage) error { return errors.Errorf("test error") }),
		WithKafkaClient(mockClient),
		WithLogging(newTestLogger()),
		WithReturnOnError(true),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	err = c.Consume(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "test error")

	mockClient.AssertExpectations(t)
	mockConsumerGroup.AssertExpectations(t)
}
