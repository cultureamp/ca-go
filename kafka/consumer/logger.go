package consumer

import (
	"context"
)

// NotifyError is a notify-on-error function used to report consumer handler errors.
type NotifyError func(ctx context.Context, err error, msg Message)

type ClientLogger interface {
	// Kafka-go Logger interface
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}
