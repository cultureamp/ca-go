package consumer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestNewConsumerWithMockClient(t *testing.T) {
	ctx := context.Background()
	mockClient := newMockKafkaClient()

	c, err := NewConsumer(
		WithBrokers([]string{"localhost:9001"}),
		WithTopics([]string{"test-topic"}),
		WithGroupId("group_id"),
		WithAssignor("roundrobin"),
		WithHandler(func(ctx context.Context, msg *ConsumerMessage) error { return nil }),
		WithKafkaClient(mockClient),
	)
	assert.NotNil(t, c)
	assert.Nil(t, err)

	err = c.Consume(ctx)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing mock implementation")
}
