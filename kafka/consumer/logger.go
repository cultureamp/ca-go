package consumer

import (
	"fmt"

	"github.com/cultureamp/ca-go/log"
	"github.com/go-errors/errors"
)

type ClientLogger interface {
	// Kafka-go Logger interface
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

type noopClientLogger struct{}

func (noop *noopClientLogger) Infof(string, ...interface{})  {}
func (noop *noopClientLogger) Errorf(string, ...interface{}) {}

type autoClientLogger struct{}

func (l *autoClientLogger) Infof(msg string, args ...interface{}) {
	details := fmt.Sprintf(msg, args...)
	log.Info("kafka_reader_info").Details(details)
}

func (l *autoClientLogger) Errorf(msg string, args ...interface{}) {
	details := fmt.Sprintf(msg, args...)
	err := errors.Errorf("kafka_reader_error")
	log.Error("kafka_reader_error", err).Details(details)
}
