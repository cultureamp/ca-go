package consumer

import (
	"context"

	"github.com/IBM/sarama"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ConsumerMessage struct {
	*sarama.ConsumerMessage
}

type Handler func(ctx context.Context, msg *ConsumerMessage) error

type messageHandler struct {
	dispatchMessage Handler
}

func newMessageHandler(handler Handler) *messageHandler {
	return &messageHandler{
		dispatchMessage: handler,
	}
}

func (h *messageHandler) dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// todo: add retries, etc.

	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.consumer.handle", tracer.ResourceName(msg.Topic))
	defer span.Finish()

	message := &ConsumerMessage{msg}
	if err := h.dispatchMessage(ctx, message); err != nil {
		return err
	}

	return nil
}
