package consumer

import (
	"context"

	"github.com/IBM/sarama"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Message struct {
	*sarama.ConsumerMessage
}

type Handler func(ctx context.Context, msg *Message) error

type messageHandler struct {
	dispatchMessage Handler
}

func newMessageHandler(handler Handler) *messageHandler {
	return &messageHandler{
		dispatchMessage: handler,
	}
}

func (h *messageHandler) dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// add retries, etc.
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.consumer.handle", tracer.ResourceName(msg.Topic))
	defer span.Finish()

	// add generics here to convert message to type V
	// github.com/hamba/avro
	// https://github.com/confluentinc/schema-registry/blob/master/avro-serializer/src/main/java/io/confluent/
	// kafka/serializers/AbstractKafkaAvroDeserializer.java#L203
	message := &Message{msg}
	if err := h.dispatchMessage(ctx, message); err != nil {
		return err
	}

	return nil
}
