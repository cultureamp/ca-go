package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConsumer(t *testing.T) {
	c, err := NewSubscriber()
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing brokers")

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing topics")

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing group")

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing schema registry URL")

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithSchemaRegistryURL("http://localhost:8081"),
	)

	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing message handler")

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithSchemaRegistryURL("http://localhost:8081"),
		WithHandler(func(ctx context.Context, msg *ReceivedMessage) error { return nil }),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithHandler(func(ctx context.Context, msg *ReceivedMessage) error { return nil }),
		WithSchemaRegistryURL("http://localhost:8081"),
		WithAssignor("abc"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unrecognized consumer group partition assignor")

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithHandler(func(ctx context.Context, msg *ReceivedMessage) error { return nil }),
		WithSchemaRegistryURL("http://localhost:8081"),
		WithAssignor("roundrobin"),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithHandler(func(ctx context.Context, msg *ReceivedMessage) error { return nil }),
		WithSchemaRegistryURL("http://localhost:8081"),
		WithAssignor("roundrobin"),
		WithVersion("abc"),
	)
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid kafka version")

	c, err = NewSubscriber(
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithHandler(func(ctx context.Context, msg *ReceivedMessage) error { return nil }),
		WithSchemaRegistryURL("http://localhost:8081"),
		WithAssignor("roundrobin"),
		WithVersion("1.0.0"),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	// full list for coverage purposes
	c, err = NewSubscriber(
		WithConsumerID("abc.123.uuid"),
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithHandler(func(ctx context.Context, msg *ReceivedMessage) error { return nil }),
		WithSchemaRegistryURL("http://localhost:8081"),
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
	mockConsumerClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockConsumerClaim)
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	mockDecoder := newMockArvoDecoder(mockSchemaRegistryClient)

	schema := testSubscriberSchema(t)
	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockClient.On("MarkMessageConsumed", mock.Anything, mock.Anything, mock.Anything)
	mockSession.On("Context").Return(ctx)
	mockConsumerClaim.On("Topic").Return("test-consumer-ctx-deadline-topic")
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)
	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(schema, nil)
	mockDecoder.On("Decode", mock.Anything).Return(`{"id": 123,"name": "test"}`, nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockConsumerClaim.On("Messages").Return(receiverChannel)

	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		assert.Equal(t, `{"id": 123,"name": "test"}`, msg.DecodedText)
		return nil
	}

	c := testConsumer(t, kafkaClient(mockClient), mockDecoder, mockReceiver, int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.ConsumeAll(ctx)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockConsumerClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}

func TestConsumerWithDecodeError(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockConsumerClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockConsumerClaim)
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	mockDecoder := newMockArvoDecoder(mockSchemaRegistryClient)

	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockSession.On("Context").Return(ctx)
	mockConsumerClaim.On("Topic").Return("test-consumer-decode-error-topic")
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)
	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(nil, errors.Errorf("test schema error"))
	mockDecoder.On("Decode", mock.Anything).Return("", errors.Errorf("test decode error"))

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockConsumerClaim.On("Messages").Return(receiverChannel)

	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		return nil
	}

	c := testConsumer(t, kafkaClient(mockClient), mockDecoder, mockReceiver, int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.ConsumeAll(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "test decode error")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockConsumerClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}

func TestConsumerWithHandlerError(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockConsumerClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockConsumerClaim)
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	mockDecoder := newMockArvoDecoder(mockSchemaRegistryClient)

	schema := testSubscriberSchema(t)
	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockSession.On("Context").Return(ctx)
	mockConsumerClaim.On("Topic").Return("test-consumer-handle-error-topic")
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)
	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(schema, nil)
	mockDecoder.On("Decode", mock.Anything).Return(`{"id": 123,"name": "test"}`, nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockConsumerClaim.On("Messages").Return(receiverChannel)

	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		return errors.Errorf("test handler error")
	}

	c := testConsumer(t, kafkaClient(mockClient), mockDecoder, mockReceiver, int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.ConsumeAll(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "test handler error")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockConsumerClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}

func TestConsumerWithChannelError(t *testing.T) {
	ctx := context.Background()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockConsumerClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockConsumerClaim)
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	mockDecoder := newMockArvoDecoder(mockSchemaRegistryClient)

	schema := testSubscriberSchema(t)
	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockSession.On("Context").Return(ctx)
	mockConsumerClaim.On("Topic").Return("test-consumer-channel-error-topic")
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)
	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(schema, nil)
	mockDecoder.On("Decode", mock.Anything).Return(`{"id": 123,"name": "test"}`, nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	close(mockChannel)
	mockConsumerClaim.On("Messages").Return(receiverChannel)

	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		return nil
	}

	c := testConsumer(t, kafkaClient(mockClient), mockDecoder, mockReceiver, int64(0), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.ConsumeAll(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "message channel closed")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockConsumerClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}

func TestConsumerWithDoubleSubscribeAndSingleStop(t *testing.T) {
	ctx := context.Background()
	// ctx will expire in 1 second
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancelCtx := context.WithDeadline(ctx, deadline)
	defer cancelCtx()

	mockClient := newMockKafkaClient()
	mockSession := newMockConsumerGroupSession()
	mockConsumerClaim := newMockConsumerGroupClaim()
	mockGroup := newMockConsumerGroup(mockSession, mockConsumerClaim)
	mockSchemaRegistryClient := newMockSchemaRegistryClient()
	mockDecoder := newMockArvoDecoder(mockSchemaRegistryClient)

	schema := testSubscriberSchema(t)
	mockClient.On("NewConsumerGroup", mock.Anything, mock.Anything, mock.Anything).Return(mockGroup, nil)
	mockClient.On("MarkMessageConsumed", mock.Anything, mock.Anything, mock.Anything)
	mockSession.On("Context").Return(ctx)
	mockConsumerClaim.On("Topic").Return("test-consumer-double-subscribe-single-stop-topic")
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockGroup.On("Close").Return(nil)
	mockSchemaRegistryClient.On("GetSchemaByID", mock.Anything).Return(schema, nil)
	mockDecoder.On("Decode", mock.Anything).Return(`{"id": 123,"name": "test"}`, nil)

	mockChannel := make(chan *sarama.ConsumerMessage, 10)
	var receiverChannel (<-chan *sarama.ConsumerMessage)
	receiverChannel = mockChannel
	mockConsumerClaim.On("Messages").Return(receiverChannel)

	mockReceiver := func(ctx context.Context, msg *ReceivedMessage) error {
		return nil
	}

	c := testConsumer(t, kafkaClient(mockClient), mockDecoder, mockReceiver, int64(3), mockChannel)
	assert.NotNil(t, c)

	// blocks until Kafka rebalance, handler error or context.Done
	err := c.ConsumeAll(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "context deadline exceeded")

	err = c.ConsumeAll(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "consumer group already running!")

	// after finished, clean up
	err = c.Stop()
	assert.Nil(t, err)

	err = c.Stop()
	assert.Nil(t, err)

	mockClient.AssertExpectations(t)
	mockSession.AssertExpectations(t)
	mockConsumerClaim.AssertExpectations(t)
	mockGroup.AssertExpectations(t)
}

func testConsumer(t *testing.T, client kafkaClient, decoder decoder, receiver Receiver, numMessages int64, ch chan *sarama.ConsumerMessage) *Subscriber {
	// push a few messages into the channel
	for i := range numMessages {
		saramaMessage := &sarama.ConsumerMessage{
			Topic:     "test",
			Partition: 1,
			Key:       []byte("key"),
			Value:     []byte(`{"id": 123,"name": "test"}`),
			Offset:    i,
			Timestamp: time.Now(),
			Headers:   nil,
		}
		ch <- saramaMessage
	}

	c, err := NewSubscriber(
		WithKafkaClient(client),
		WithAvroDecoder(decoder),
		WithBrokers([]string{"localhost:9092"}),
		WithTopics([]string{"test-topic"}),
		WithGroupID("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(receiver),
		WithSchemaRegistryURL("http://localhost:8081"),
		WithLogging(newTestLogger()),
		WithReturnOnClientDispathError(true),
	)
	assert.Nil(t, err)

	return c
}

func testSubscriberSchema(t *testing.T) *srclient.Schema {
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
