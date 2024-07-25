package consumer

import "github.com/go-errors/errors"

var (
	errClosedMessageChannel = errors.Errorf("consumer: message channel closed or in error state")
	errDoneMessageChannel   = errors.Errorf("consumer: message channel done")
)

type dispatchHandlerError struct {
	err error
}

func newDispatchHandlerError(topic string, reason error) dispatchHandlerError {
	return dispatchHandlerError{
		err: errors.Errorf("consumer[%s]: handler dispatch failed: err=%w", topic, reason),
	}
}

func (e dispatchHandlerError) Error() string {
	return e.err.Error()
}
