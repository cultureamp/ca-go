package consumer

import (
	"context"

	"github.com/IBM/sarama"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Message struct {
	json string // json
}

type Receiver func(ctx context.Context, msg *Message) error

type messageHandler struct {
	receiver Receiver
	decoder  decoder
}

func newMessageHandler(receiver Receiver, decoder decoder) *messageHandler {
	return &messageHandler{
		receiver: receiver,
		decoder:  decoder,
	}
}

func (h *messageHandler) dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// add retries, etc.
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.consumer.handle", tracer.ResourceName(msg.Topic))
	defer span.Finish()

	// add generics here to convert message to type V ???

	json, err := h.decoder.DecodeAsString(msg)
	if err != nil {
		return err
	}

	message := &Message{json}
	if err := h.receiver(ctx, message); err != nil {
		return err
	}

	return nil
}
