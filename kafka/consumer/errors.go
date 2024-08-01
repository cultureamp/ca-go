package consumer

import "github.com/go-errors/errors"

var errClosedMessageChannel = errors.Errorf("consumer: message channel closed")

type messageHandlerError struct {
	err error
}

func newMessageHandlerError(topic string, reason error) messageHandlerError {
	return messageHandlerError{
		err: errors.Errorf("consumer[%s]: handler dispatch failed: err=%w", topic, reason),
	}
}

// Error implements the error interface.
func (e messageHandlerError) Error() string {
	return e.err.Error()
}
