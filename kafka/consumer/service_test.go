package consumer

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceWithCancelledContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var calls atomic.Int32
	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		calls.Add(1)
		return nil
	}

	s := testService(t, ctx, mockReceiver, 4)
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
	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		calls.Add(1)
		return nil
	}

	s := testService(t, ctx, mockReceiver, 2)
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

func TestServiceWithHandlerError(t *testing.T) {
	ctx := context.Background()

	var calls atomic.Int32
	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		calls.Add(1)
		return errors.Errorf("test error")
	}

	s := testService(t, ctx, mockReceiver, 3)
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
	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		calls.Add(1)
		return errors.Errorf("test error")
	}

	s := testService(t, ctx, mockReceiver, 3)
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

func testService(t *testing.T, ctx context.Context, receiver Receiver, numMessages int64) *Service {
	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockConsumer := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockConsumer)
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	mockDecoder := newMockArvoDecoder(mockSchemaRegistryClient)

	schema := testServiceSchema(t)
	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockClient.On("CommitMessage", mock.Anything, mock.Anything)
	mockSession.On("Context").Return(ctx)
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)
	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(schema, nil)
	mockDecoder.On("Decode", mock.Anything).Return(`
	{
		"id": 123,
		"name": "test"
	}`, nil)

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
	mockConsumer.On("Messages").Return(receiverChannel)

	s, err := NewService(
		WithKafkaClient(mockClient),
		WithAvroDecoder(mockDecoder),
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(receiver),
		WithSchemaRegistryURL("http://localhost:8081"),
		WithLogging(newTestLogger()),
		WithReturnOnClientDispathError(false),
	)
	assert.Nil(t, err)

	return s
}

func testServiceSchema(t *testing.T) *srclient.Schema {
	codec, err := goavro.NewCodec(`
	{
		"type": "record",
		"name": "TestObject",
		"namespace": "ca.dataedu",
		"fields": [
			{
			"name": "id",
			"type": [
				"null",
				"int"
			],
			"default": null
			},
			{
			"name": "name",
			"type": [
				"null",
				"string"
			],
			"default": null
			}
		]
		}`)
	assert.Nil(t, err)

	schema, err := srclient.NewSchema(1, "test", srclient.Avro, 1, nil, codec, nil)
	assert.Nil(t, err)

	return schema
}
