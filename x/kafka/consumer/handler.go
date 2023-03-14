package consumer

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
	kafkatrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/segmentio/kafka.go.v0"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type handlerExecutor struct {
	ConsumerID            string
	GroupID               string
	DataDogTracingEnabled bool
	BackOffConstructor    HandlerRetryBackOffConstructor
	NotifyErr             NotifyError
}

func (h *handlerExecutor) execute(ctx context.Context, msg kafka.Message, handler Handler) error {
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
			return fmt.Errorf("unable to extract data dog span context from kafka message: %w", err)
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
			return fmt.Errorf("consumer handler error: %w", ctx.Err())
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
			if h.NotifyErr != nil {
				h.NotifyErr(ctx, err, consumerMsg)
			}
			continue
		}

		return nil
	}
}
