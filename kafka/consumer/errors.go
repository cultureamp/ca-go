package consumer

import "github.com/go-errors/errors"

var errClosedMessageChannel = errors.Errorf("consumer: message channel closed or in error state")

type dispatchHandlerError struct {
	err error
}

func newDispatchHandlerError(reason error) dispatchHandlerError {
	return dispatchHandlerError{
		err: errors.Errorf("consumer: handler dispatch failed: err=%w", reason),
	}
}

func (e dispatchHandlerError) Error() string {
	return e.err.Error()
}
