package kafkatest

import (
	"context"
	"fmt"
	"testing"

	"github.com/heetch/avro"
	"github.com/heetch/avro/avroregistry"
	"github.com/stretchr/testify/require"
)

// TestRegistry is a kafka registry used to easily encode/decode messages for
// testing.
//
// TestRegistry client accepts a type parameter EventType, which is the raw type
// of the event to encode and decode into.
type TestRegistry[EventType any] struct {
	subject  string
	registry *avroregistry.Registry
	encoder  *avro.SingleEncoder
	decoder  *avro.SingleDecoder
}

// NewTestRegistry returns a new TestRegistry. During creation the schema for
// EventType is automatically generated and registered under a test subject.
func NewTestRegistry[EventType any](t *testing.T, ctx context.Context, hostPort string, subject string) *TestRegistry[EventType] {
	r, err := avroregistry.New(avroregistry.Params{ServerURL: fmt.Sprint("http://", hostPort)})
	require.NoError(t, err)

	var eventType EventType

	// Register schema for subject
	avroType, err := avro.TypeOf(eventType)
	require.NoError(t, err, "cannot generate avro schema for %T: %w", eventType, err)
	_, err = r.Register(ctx, subject, avroType)
	require.NoError(t, err)
	t.Cleanup(func() {
		deleteErr := r.DeleteSubject(context.Background(), subject)
		require.NoError(t, deleteErr, "error deleting subject")
	})

	return &TestRegistry[EventType]{
		subject:  subject,
		registry: r,
		encoder:  avro.NewSingleEncoder(r.Encoder(subject), nil),
		decoder:  avro.NewSingleDecoder(r.Decoder(), nil),
	}
}

// Encode returns the event marshaled using the Avro binary encoding.
func (r *TestRegistry[EventType]) Encode(t *testing.T, ctx context.Context, event EventType) []byte {
	b, err := r.encoder.Marshal(ctx, event)
	require.NoError(t, err)
	return b
}

// Decode unmarshals the given message into a new declaration of EventType and
// returns it.
func (r *TestRegistry[EventType]) Decode(t *testing.T, ctx context.Context, data []byte) EventType {
	var event EventType
	_, err := r.decoder.Unmarshal(ctx, data, &event)
	require.NoError(t, err)
	return event
}
