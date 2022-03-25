package lambdafunction

import (
	"context"
)

// HandlerOf[TIn] is a lambda handler that models a Lambda handler
// function that expects a payload of TIn and returns an error.
type HandlerOf[TIn any] func(context.Context, TIn) error

// HandlerWithOutputOf[TIn] is a lambda handler that models a Lambda
// handler function that expects a payload of TIn and returns a tuple of an
// output type (TOut) and an error.
type HandlerWithOutputOf[TIn any, TOut any] func(context.Context, TIn) (TOut, error)
