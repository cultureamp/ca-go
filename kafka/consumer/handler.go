package consumer

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-errors/errors"
	"github.com/segmentio/kafka-go"

	kafkatrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/segmentio/kafka.go.v0"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type messageHandler struct {
	ConsumerID            string
	GroupID               string
	DataDogTracingEnabled bool
	BackOffConstructor    HandlerRetryBackOffConstructor
	clientNotify          NotifyError
}

func (h *messageHandler) execute(ctx context.Context, msg kafka.Message, handler Handler) error {
	var err error
	var backOff backoff.BackOff

	if h.BackOffConstructor == nil {
		backOff = &backoff.StopBackOff{}
	} else {
		backOff = h.BackOffConstructor()
	}

	attempt := 0
	ticker := backoff.NewTicker(backOff)
	defer ticker.Stop()

	if h.DataDogTracingEnabled {
		spanCtx, err := kafkatrace.ExtractSpanContext(msg)
		if err != nil {
			return errors.Errorf("unable to extract data dog span context from kafka message: %w", err)
		}
		span := tracer.StartSpan("consumer.handle", tracer.ChildOf(spanCtx))
		defer span.Finish()
		ctx = tracer.ContextWithSpan(ctx, span)
	}

	for {
		select {
		case <-ctx.Done():
			if err == nil {
				return ctx.Err()
			}
			return errors.Errorf("consumer handler error: %w", err)
		case _, ok := <-ticker.C:
			if !ok {
				return err
			}
		}

		attempt++

		consumerMsg := Message{
			Message: msg,
			Metadata: Metadata{
				GroupID:    h.GroupID,
				ConsumerID: h.ConsumerID,
				Attempt:    attempt,
			},
		}

		if err = handler(ctx, consumerMsg); err != nil {
			h.clientNotify(ctx, err, consumerMsg)
			continue
		}

		return nil
	}
}
