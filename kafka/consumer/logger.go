package consumer

import (
	"context"
	"fmt"

	"github.com/cultureamp/ca-go/log"
	"github.com/go-errors/errors"
)

// NotifyError is a notify-on-error function used to report consumer handler errors.
type NotifyError func(ctx context.Context, err error, msg Message)

type ClientLogger interface {
	// Kafka-go Logger interface
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

type noopClientLogger struct{}

func (noop *noopClientLogger) Infof(string, ...interface{})  {}
func (noop *noopClientLogger) Errorf(string, ...interface{}) {}

type autoKafkaLogger struct{}

func (l *autoKafkaLogger) Infof(msg string, args ...interface{}) {
	details := fmt.Sprintf(msg, args...)
	log.Info("kafka_reader_info").Details(details)
}

func (l *autoKafkaLogger) Errorf(msg string, args ...interface{}) {
	details := fmt.Sprintf(msg, args...)
	err := errors.Errorf("kafka_reader_error")
	log.Error("kafka_reader_error", err).Details(details)
}

func autoClientNotifyError(_ context.Context, err error, msg Message) {
	log.Error("auto_consumer_notify_error", err).
		WithSystemTracing().
		Properties(log.SubDoc().
			Str("topic", msg.Topic).
			Str("key", string(msg.Key)).
			Str("value", string(msg.Value)),
		).Details("error consuming message")
}
