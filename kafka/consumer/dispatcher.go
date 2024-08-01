package consumer

import (
	"context"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type dispatcher interface {
	Dispatch(ctx context.Context, msg *ReceivedMessage) error
}

// ReceivedMessage contains the underlying kafka message,
// as well as the Avro DecodedText from the raw kafka message.Value.
type ReceivedMessage struct {
	Timestamp time.Time
	Topic     string
	Offset    int64
	Key       []byte
	Value     string // typically json, client needs to json.Unmarshal to specific domain struct
}

// Receiver is the client's message handler that processes the ReceivedMessage.
// Returning an error will cause the consumer to stop consuming messages.
type Receiver func(ctx context.Context, msg *ReceivedMessage) error

type dispatchHandler struct {
	receiver Receiver
}

func newDispatcher(receiver Receiver) *dispatchHandler {
	return &dispatchHandler{
		receiver: receiver,
	}
}

// Dispatch handles the kafka message by decoding the message and calling the client's Receiver.
// Returning an error will cause the consumer to stop consuming messages.
func (h *dispatchHandler) Dispatch(ctx context.Context, msg *ReceivedMessage) error {
	// add retries, etc.
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.consumer.dispatch", tracer.ResourceName(msg.Topic))

	// Set tags
	span.SetTag("kafka.consumer.dispatch.message.key", string(msg.Key))
	span.SetTag("kafka.consumer.dispatch.message.offset", msg.Offset)
	span.SetTag("kafka.consumer.dispatch.message.timestamp", msg.Timestamp)
	span.SetTag("kafka.consumer.dispatch.message.topic", msg.Topic)

	if err := h.receiver(ctx, msg); err != nil {
		span.Finish(tracer.WithError(err))
		return err
	}

	span.Finish()
	return nil
}
