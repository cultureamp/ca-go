package consumer

import "github.com/go-errors/errors"

var errClosedMessageChannel = errors.Errorf("consumer: message channel closed")

type dispatchHandlerError struct {
	err error
}

func newDispatchHandlerError(topic string, reason error) dispatchHandlerError {
	return dispatchHandlerError{
		err: errors.Errorf("consumer[%s]: handler dispatch failed: err=%w", topic, reason),
	}
}

// Error implements the error interface.
func (e dispatchHandlerError) Error() string {
	return e.err.Error()
}
