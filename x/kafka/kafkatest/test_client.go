package kafkatest

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const timeout = 30 * time.Second

// TestClientConfig is a configuration object used to create a new TestClient.
type TestClientConfig struct {
	KafkaBrokerHostPort    string
	SchemaRegistryHostPort string
	NumTopicPartitions     int
}

// TestClient is a Kafka client that allows you to easily setup up a Kafka topic
// and interact with it for testing.
//
// TestClient client accepts a type parameter EventType, which is the raw type of
// the event being published/consumed to and from the topic. It is used to within
// TestClient to simplify encoding and decoding processes.
type TestClient[EventType any] struct {
	// Topic is the name of the created test topic
	Topic    string
	config   TestClientConfig
	client   *kafka.Client
	registry *TestRegistry[EventType]
	writer   *kafka.Writer
	reader   *kafka.Reader
}

// NewTestClient creates a new test topic in Kafka and returns a TestClient.
func NewTestClient[EventType any](t *testing.T, ctx context.Context, cfg TestClientConfig) *TestClient[EventType] {
	client := &kafka.Client{
		Addr:    kafka.TCP(cfg.KafkaBrokerHostPort),
		Timeout: timeout,
	}

	topicConfig := kafka.TopicConfig{
		Topic:             uuid.New().String(),
		NumPartitions:     cfg.NumTopicPartitions,
		ReplicationFactor: 1,
	}

	_, err := client.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Addr:         kafka.TCP(cfg.KafkaBrokerHostPort),
		Topics:       []kafka.TopicConfig{topicConfig},
		ValidateOnly: false,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, deleteErr := client.DeleteTopics(context.Background(), &kafka.DeleteTopicsRequest{
			Addr:   kafka.TCP(cfg.KafkaBrokerHostPort),
			Topics: []string{topicConfig.Topic},
		})
		assert.NoError(t, deleteErr, "error deleting topic %s", topicConfig.Topic)
	})

	subject := fmt.Sprintf("testsubject-%s", uuid.New().String())

	return &TestClient[EventType]{
		Topic:    topicConfig.Topic,
		config:   cfg,
		client:   client,
		registry: NewTestRegistry[EventType](t, ctx, cfg.SchemaRegistryHostPort, subject),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaBrokerHostPort),
			Topic:    topicConfig.Topic,
			Balancer: &kafka.RoundRobin{},
		},
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.KafkaBrokerHostPort},
			Topic:   topicConfig.Topic,
			Dialer:  kafka.DefaultDialer,
		}),
	}
}

// PublishEvents generates kafka messages from the passed in raw event types and
// publishes them to the test topic.
func (c *TestClient[EventType]) PublishEvents(t *testing.T, ctx context.Context, events ...EventType) {
	var msgs []kafka.Message

	for i, event := range events {
		msg := kafka.Message{
			Value: c.registry.Encode(t, ctx, event),
			Time:  time.Now(),
			Key:   []byte(strconv.Itoa(i)),
		}
		msgs = append(msgs, msg)
	}

	c.PublishMessages(t, ctx, msgs...)
}

// PublishMessages publishes the passed in kafka messages to the test topic.
func (c *TestClient[EventType]) PublishMessages(t *testing.T, ctx context.Context, msgs ...kafka.Message) {
	require.NoError(t, c.writer.WriteMessages(ctx, msgs...))
}

// ConsumeEvent reads the next message from the topic and commits the offset.
// The message is decoded into a new declaration of EventType and returned.
func (c *TestClient[EventType]) ConsumeEvent(t *testing.T, ctx context.Context) EventType {
	msg, err := c.reader.ReadMessage(ctx)
	require.NoError(t, err)
	return c.registry.Decode(t, ctx, msg.Value)
}

// TopicMessageCount returns the total number of messages that are currently in
// the test topic.
func (c *TestClient[EventType]) TopicMessageCount(t *testing.T, ctx context.Context) int64 {
	metadataInput := kafka.MetadataRequest{Topics: []string{c.Topic}}
	metadataOutput, err := c.client.Metadata(ctx, &metadataInput)
	require.NoError(t, err)

	offsetsInput := kafka.ListOffsetsRequest{Topics: make(map[string][]kafka.OffsetRequest)}
	for _, topic := range metadataOutput.Topics {
		require.NoError(t, err)

		var topicPartitions []kafka.OffsetRequest
		for _, partition := range topic.Partitions {
			topicPartitions = append(topicPartitions, kafka.OffsetRequest{
				Partition: partition.ID,
				Timestamp: kafka.LastOffset,
			})
		}
		offsetsInput.Topics[topic.Name] = topicPartitions
	}
	offsetsOutput, err := c.client.ListOffsets(ctx, &offsetsInput)
	require.NoError(t, err)

	var numMessages int64 = 0
	for _, partitions := range offsetsOutput.Topics {
		for _, partition := range partitions {
			require.NoError(t, err)
			numMessages += partition.LastOffset
		}
	}
	return numMessages
}

// Client returns the internal kafka client.
func (c *TestClient[EventType]) Client() *kafka.Client {
	return c.client
}

// Writer returns the internal kafka writer.
func (c *TestClient[EventType]) Writer() *kafka.Writer {
	return c.writer
}

// Reader returns the internal kafka reader.
func (c *TestClient[EventType]) Reader() *kafka.Reader {
	return c.reader
}

// Registry returns the internal kafka test registry.
func (c *TestClient[EventType]) Registry() *TestRegistry[EventType] {
	return c.registry
}
