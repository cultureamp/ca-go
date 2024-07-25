package consumer

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type handler interface {
	Dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error
}

// ReceivedMessage contains the underlying kafka message,
// as well as the Avro DecodedText from the raw kafka message.Value.
type ReceivedMessage struct {
	Timestamp   time.Time
	Topic       string
	Offset      int64
	Key, Value  []byte
	DecodedText string // typically json, client needs to json.Unmarshal to specific domain struct
}

// Receiver is the client's message handler that processes the ReceivedMessage.
// Returning an error will cause the consumer to stop consuming messages.
type Receiver func(ctx context.Context, msg *ReceivedMessage) error

type dispatchHandler struct {
	receiver Receiver
	decoder  decoder
}

func newHandler(receiver Receiver, decoder decoder) *dispatchHandler {
	return &dispatchHandler{
		receiver: receiver,
		decoder:  decoder,
	}
}

// Dispatch handles the kafka message by decoding the message and calling the client's Receiver.
// Returning an error will cause the consumer to stop consuming messages.
func (h *dispatchHandler) Dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error {
	text, err := h.decoder.Decode(msg)
	if err != nil {
		return err
	}

	message := &ReceivedMessage{
		Timestamp:   msg.Timestamp,
		Topic:       msg.Topic,
		Offset:      msg.Offset,
		Key:         msg.Key,
		Value:       msg.Value,
		DecodedText: text,
	}
	if err := h.dispatchToClient(ctx, message); err != nil {
		return err
	}

	return nil
}

func (h *dispatchHandler) dispatchToClient(ctx context.Context, msg *ReceivedMessage) error {
	// add retries, etc.
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.consumer.handle", tracer.ResourceName(msg.Topic))
	defer span.Finish()

	if err := h.receiver(ctx, msg); err != nil {
		return err
	}

	return nil
}
