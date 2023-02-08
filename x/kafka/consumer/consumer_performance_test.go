//go:build performance

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

const (
	kafkaHostPort    = "localhost:9093"
	registryHostPort = "localhost:8081"
)

type Event struct {
	Int int
	Str string
}

func TestConsumerGroup_Run_performance(t *testing.T) {
	ctx := context.Background()

	getBatchKeyFn := func(message kafka.Message) (string, error) {
		i, err := strconv.Atoi(string(message.Key))
		require.NoError(t, err)
		return strconv.Itoa(i % 100), nil
	}

	tests := []struct {
		name          string
		consumerCount int
		numPartitions int
		numMessages   int
		opts          []Option
	}{
		{
			name:          "25 consumer in group (100 partitions, 10,000 messages)",
			consumerCount: 25,
			numPartitions: 100,
			numMessages:   10000,
			opts:          []Option{WithExplicitCommit()},
		},
		{
			name:          "25 consumer in group (100 partitions, 10,000 messages)",
			consumerCount: 25,
			numPartitions: 100,
			numMessages:   10000,
			opts:          []Option{WithMessageBatching(100, getBatchKeyFn)},
		},
		{
			name:          "1111125 consumer in group (100 partitions, 10,000 messages)",
			consumerCount: 1,
			numPartitions: 100,
			numMessages:   400,
			opts:          []Option{WithExplicitCommit()},
		},
		{
			name:          "2225 consumer in group (100 partitions, 10,000 messages)",
			consumerCount: 1,
			numPartitions: 100,
			numMessages:   1000,
			opts:          []Option{WithMessageBatching(100, getBatchKeyFn)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := kafkatest.NewTestClient[Event](t, ctx, kafkatest.TestClientConfig{
				KafkaBrokerHostPort:    kafkaHostPort,
				SchemaRegistryHostPort: registryHostPort,
				NumTopicPartitions:     tt.numPartitions,
			})

			cfg := GroupConfig{
				Count:   tt.consumerCount,
				Brokers: []string{kafkaHostPort},
				Topic:   tc.Topic,
				GroupID: uuid.New().String(),
			}
			c := NewGroup(kafka.DefaultDialer, cfg, tt.opts...)

			duration := consumerGroupProcessDuration(t, c, tc, tt.numMessages)
			t.Logf("[Group size: %d, Partitions: %d] processing duration for %d messages: %s", cfg.Count, tt.numPartitions, tt.numMessages, duration.String())
		})
	}

}

func consumerGroupProcessDuration(t *testing.T, c *Group, tc *kafkatest.TestClient[Event], numMessages int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	ctx, timeoutCancel := context.WithTimeout(context.Background(), time.Minute)
	defer timeoutCancel()

	var events []Event
	for i := 0; i < numMessages; i++ {
		events = append(events, Event{
			Int: i,
			Str: strconv.Itoa(i),
		})
	}
	tc.PublishEvents(t, ctx, events...)

	consumerCtx, consumerCancel := context.WithCancel(ctx)

	var numConsumed int64
	handler := func(ctx context.Context, msg Message) error {
		time.Sleep(50 * time.Millisecond) // simulate http request
		assert.NotEmpty(t, msg.Value)
		atomic.AddInt64(&numConsumed, 1)
		if numConsumed == int64(numMessages) && numConsumed == tc.TopicMessageCount(t, ctx) {
			consumerCancel()
		}
		return nil
	}

	start := time.Now()
	errCh := c.Run(consumerCtx, handler)
	for err := range errCh {
		if !errors.Is(err, context.Canceled) {
			require.NoError(t, err)
		}
	}
	return time.Now().Sub(start)
}
