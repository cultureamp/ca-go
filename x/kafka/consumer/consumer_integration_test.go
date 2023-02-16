//go:build integration

// All integration tests within this file require a local Kafka environment to be
// running with docker using the command `make up`. Once testing is done, the
// docker environment can be spun down with `make down`

package consumer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	brokerHostPort         = "localhost:9093"
	schemaRegistryHostPort = "localhost:8081"
	timeout                = 60 * time.Second
)

type TestEvent struct {
	Int int
	Str string
}

// TestConsumerGroup_Run_integration tests consumers with a local Kafka environment.
// Each test case configures the consumer and topic with unique settings in order
// to test performance across a range of scenarios. Total time taken and average
// messages processed per second are also logged at the end of each test run.
func TestConsumerGroup_Run_integration(t *testing.T) {
	// Simulate duration of an average http request
	const handlerSleepDuration = 50 * time.Millisecond

	tests := []struct {
		name          string
		partitions    int
		numMessages   int
		consumerCount int
		opts          []Option
	}{
		{
			name:          "1 consumer, 1 partition, 25 messages)",
			partitions:    1,
			numMessages:   25,
			consumerCount: 1,
		},
		{
			name:          "10 consumers, 10 partitions, 250 messages)",
			partitions:    10,
			numMessages:   250,
			consumerCount: 10,
		},
		{
			name:          "10 consumers (explicit commit), 10 partitions, 250 messages)",
			partitions:    10,
			numMessages:   250,
			consumerCount: 10,
			opts:          []Option{WithExplicitCommit()},
		},
		{
			name:          "1 consumer (batching), 1 partitions, 500 messages)",
			partitions:    1,
			numMessages:   500,
			consumerCount: 1,
			opts:          []Option{WithMessageBatching(20, newGetOrderingKey(t, 20))},
		},
		{
			name:          "1 consumer (batching), 10 partitions, 100 messages)",
			partitions:    10,
			numMessages:   100,
			consumerCount: 1,
			opts:          []Option{WithMessageBatching(20, newGetOrderingKey(t, 20))},
		},
		{
			name:          "10 consumers (batching), 10 partitions, 500 messages",
			partitions:    10,
			numMessages:   500,
			consumerCount: 10,
			opts:          []Option{WithMessageBatching(20, newGetOrderingKey(t, 20))},
		},
		{
			name:          "10 consumers (batching), 100 partitions, 500 messages",
			partitions:    10,
			numMessages:   500,
			consumerCount: 10,
			opts:          []Option{WithMessageBatching(20, newGetOrderingKey(t, 20))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
			defer timeoutCancel()

			tc := kafkatest.NewTestClient[TestEvent](t, ctx, kafkatest.TestClientConfig{
				KafkaBrokerHostPort:    brokerHostPort,
				SchemaRegistryHostPort: schemaRegistryHostPort,
				NumTopicPartitions:     tt.partitions,
			})

			cfg := GroupConfig{
				Count:   tt.consumerCount,
				Brokers: []string{brokerHostPort},
				Topic:   tc.Topic,
				GroupID: uuid.New().String(),
			}
			c := NewGroup(kafka.DefaultDialer, cfg, tt.opts...)

			numPublish := tt.numMessages
			publishDummyEvents(t, ctx, tc, numPublish)

			stopCh := make(chan bool)
			numConsumed := new(safeCounter)
			handler := func(ctx context.Context, msg Message) error {
				time.Sleep(handlerSleepDuration)
				assert.NotEmpty(t, msg.Value)
				numConsumed.inc()
				i := numConsumed.val()
				if i == numPublish && i == int(tc.TopicMessageCount(t, ctx)) {
					stopCh <- true
				}
				return nil
			}

			start := time.Now()
			errCh := c.Run(ctx, handler)

		SelectLoop:
			for {
				select {
				case <-stopCh:
					c.Stop()
					break SelectLoop
				case err, ok := <-errCh:
					if ok {
						require.NoError(t, err)
					}
					break SelectLoop
				}
			}

			duration := time.Now().Sub(start)

			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("\nResults - %s\n", tt.name))
			buf.WriteString(fmt.Sprintf("- Processing time: %s\n", duration.String()))
			buf.WriteString(fmt.Sprintf("- Average messages per second: %.1f\n", float64(numConsumed.val())/duration.Seconds()))
			t.Log(buf.String())
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

func newGetOrderingKey(t *testing.T, mod int) GetOrderingKey {
	return func(ctx context.Context, message kafka.Message) string {
		i, err := strconv.Atoi(string(message.Key))
		require.NoError(t, err)
		return strconv.Itoa(i % mod)
	}
}
