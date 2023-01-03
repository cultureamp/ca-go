//go:build integration

// All integration tests within this file require a local Kafka environment to be
// running with docker using the command `make up`. Once testing is done, the
// docker environment can be spun down with `make down`

package consumer

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cultureamp/ca-go/x/kafka/kafkatest"
)

var (
	brokerHostPort         = "localhost:9093"
	schemaRegistryHostPort = "localhost:8081"
	timeout                = 15 * time.Second
)

type TestEvent struct {
	Int int
	Str string
}

func TestConsumerGroup_Run_integration(t *testing.T) {
	tests := []struct {
		name          string
		partitions    int
		numMessages   int
		consumerCount int
		opts          []Option
	}{
		{
			name:          "1 consumer in group (1 partitions, 100 messages)",
			partitions:    1,
			numMessages:   100,
			consumerCount: 1,
		},
		{
			name:          "1 consumer in group (10 partitions, 100 messages)",
			partitions:    10,
			numMessages:   100,
			consumerCount: 1,
		},
		{
			name:          "10 consumers in group (10 partitions, 1000 messages)",
			partitions:    10,
			numMessages:   1000,
			consumerCount: 10,
		},
		{
			name:          "10 consumers in group with explicit commit (10 partitions, 1000 messages)",
			partitions:    10,
			numMessages:   1000,
			consumerCount: 10,
			opts:          []Option{WithExplicitCommit()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
			defer timeoutCancel()

			helper := kafkatest.NewTestClient[TestEvent](t, ctx, kafkatest.TestClientConfig{
				KafkaBrokerHostPort:    brokerHostPort,
				SchemaRegistryHostPort: schemaRegistryHostPort,
				NumTopicPartitions:     tt.partitions,
			})

			cfg := GroupConfig{
				Count:   tt.consumerCount,
				Brokers: []string{brokerHostPort},
				Topic:   helper.Topic,
				GroupID: uuid.New().String(),
			}
			c := NewGroup(kafka.DefaultDialer, cfg)

			numPublish := tt.numMessages
			publishDummyEvents(t, ctx, helper, numPublish)

			consumerCtx, consumerCancel := context.WithCancel(ctx)

			var numConsumed int64
			handler := func(ctx context.Context, msg Message) error {
				assert.NotEmpty(t, msg.Value)
				atomic.AddInt64(&numConsumed, 1)
				if numConsumed == int64(numPublish) && numConsumed == helper.TopicMessageCount(t, ctx) {
					consumerCancel()
				}
				return nil
			}

			errCh := c.Run(consumerCtx, handler)
			for err := range errCh {
				if !errors.Is(err, context.Canceled) {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestConsumer_Run_integration(t *testing.T) {
	ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
	defer timeoutCancel()

	helper := kafkatest.NewTestClient[TestEvent](t, ctx, kafkatest.TestClientConfig{
		KafkaBrokerHostPort:    brokerHostPort,
		SchemaRegistryHostPort: schemaRegistryHostPort,
		NumTopicPartitions:     1,
	})

	cfg := Config{
		Brokers: []string{brokerHostPort},
		Topic:   helper.Topic,
	}
	c := NewConsumer(kafka.DefaultDialer, cfg)

	numPublish := 100
	publishDummyEvents(t, ctx, helper, numPublish)
	helper.TopicMessageCount(t, ctx)

	consumerCtx, consumerCancel := context.WithCancel(ctx)

	var numConsumed int64
	handler := func(ctx context.Context, msg Message) error {
		assert.NotEmpty(t, msg.Value)
		atomic.AddInt64(&numConsumed, 1)
		if numConsumed == int64(numPublish) && numConsumed == helper.TopicMessageCount(t, ctx) {
			consumerCancel()
		}
		return nil
	}

	err := c.Run(consumerCtx, handler)
	if !errors.Is(err, context.Canceled) {
		require.NoError(t, err)
	}
}

func publishDummyEvents(t *testing.T, ctx context.Context, helper *kafkatest.TestClient[TestEvent], numPublish int) {
	rand.Seed(time.Now().UnixNano())
	var events []TestEvent
	for i := 0; i < numPublish; i++ {
		events = append(events, TestEvent{
			Int: i,
			Str: strconv.Itoa(i),
		})
	}
	helper.PublishEvents(t, ctx, events...)
}
