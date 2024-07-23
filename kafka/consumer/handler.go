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

// ReceivedMessage contains the decoded message from the raw kafka message.Value.
type ReceivedMessage struct {
	Timestamp   time.Time
	Topic       string
	Offset      int64
	Key, Value  []byte
	DecodedText string // typically json, client needs to json.Unmarshal to specific domain struct
}

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

func (h *dispatchHandler) Dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// add retries, etc.
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.consumer.handle", tracer.ResourceName(msg.Topic))
	defer span.Finish()

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
	if err := h.receiver(ctx, message); err != nil {
		return err
	}

	return nil
}
