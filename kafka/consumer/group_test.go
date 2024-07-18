package consumer

import (
	"context"
	"testing"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupWithNewConsumerGroup(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockDecoder := newMockArvoDecoder()

	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.Errorf("failed to create group"))

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	handler := func(ctx context.Context, msg *ReceivedMessage) error {
		return errors.Errorf("test error")
	}

	c := testConsumer(t, client(mockClient), mockDecoder, Receiver(handler), int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.Subscribe(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to create group")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
}

func TestGroupWithConsumeError(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockConsumer := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockConsumer)
	mockDecoder := newMockArvoDecoder()

	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(sarama.ErrClosedConsumerGroup)
	mockGroup.On("Close").Return(nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)

	handler := func(ctx context.Context, msg *ReceivedMessage) error {
		return errors.Errorf("test error")
	}

	c := testConsumer(t, client(mockClient), mockDecoder, Receiver(handler), int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.Subscribe(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "kafka: tried to use a consumer group that was closed")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockConsumer.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}
