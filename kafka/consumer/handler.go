package consumer

import (
	"context"

	"github.com/IBM/sarama"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Handler interface {
	Dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error
}

type ReceivedMessage struct {
	json string // json, client needs to json.Unmarshal to struct
}

type Receiver func(ctx context.Context, msg *ReceivedMessage) error

type handler struct {
	receiver Receiver
	decoder  decoder
}

func newHandler(receiver Receiver, decoder decoder) *handler {
	return &handler{
		receiver: receiver,
		decoder:  decoder,
	}
}

func (h *handler) Dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// add retries, etc.
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.consumer.handle", tracer.ResourceName(msg.Topic))
	defer span.Finish()

	json, err := h.decoder.Decode(msg)
	if err != nil {
		return err
	}

	message := &ReceivedMessage{json}
	if err := h.receiver(ctx, message); err != nil {
		return err
	}

	return nil
}
