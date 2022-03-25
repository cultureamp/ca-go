package errorreport

import (
	"context"

	"github.com/cultureamp/ca-go/x/lambdafunction"
	"github.com/getsentry/sentry-go"
)

// LambdaOption is a function type that can be supplied to alter the behaviour of the
// LambdaMiddleware functions.
type LambdaOption func(o *lambdaOptions)

// lambdaOptions configures the way Sentry is used in the context of a
// Lambda handler wrapper.
type lambdaOptions struct {
	// Repanic configures whether to panic again after recovering from a panic.
	// Use this option if you have other panic handlers or want the default
	// behavior from AWS lambda runtime. Defaults to true.
	Repanic bool
}

// WithRepanic configures whether to panic again after reporting an error to
// Sentry. This setting defaults to true, as typically the function invocation
// should be allowed to fail by the standard Lambda mechanisms once a panic
// occurs.
func WithRepanic(repanic bool) LambdaOption {
	return func(o *lambdaOptions) {
		o.Repanic = repanic
	}
}

// LambdaMiddleware[TIn] provides error-handling middleware for a Lambda
// function that has a payload type of TIn. This suits Lambda functions like
// event processors, where the return has no payload.
func LambdaMiddleware[TIn any](nextHandler lambdafunction.HandlerOf[TIn], config ...LambdaOption) lambdafunction.HandlerOf[TIn] {
	options := configure(config)

	return func(ctx context.Context, event TIn) error {
		defer beforeHandler(ctx, options)()

		err := nextHandler(ctx, event)

		afterHandler(ctx, err)

		return err
	}
}

// LambdaWithOutputMiddleware[TIn, TOut] provides error-handling middleware for
// a Lambda function that has a payload type of TIn and returns the tuple TOut,error.
func LambdaWithOutputMiddleware[TIn any, TOut any](nextHandler lambdafunction.HandlerWithOutputOf[TIn, TOut], config ...LambdaOption) lambdafunction.HandlerWithOutputOf[TIn, TOut] {
	options := configure(config)

	return func(ctx context.Context, event TIn) (TOut, error) {
		defer beforeHandler(ctx, options)()

		out, err := nextHandler(ctx, event)

		afterHandler(ctx, err)

		return out, err
	}
}

func configure(config []LambdaOption) lambdaOptions {
	options := lambdaOptions{
		Repanic: true,
	}
	for _, c := range config {
		c(&options)
	}
	return options
}

func beforeHandler(ctx context.Context, options lambdaOptions) func() {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}
	return func() {
		if err := recover(); err != nil {
			_ = hub.RecoverWithContext(ctx, err)

			if options.Repanic {
				panic(err)
			}
		}
	}
}

func afterHandler(ctx context.Context, err error) {
	if err != nil {
		ReportError(ctx, err)
	}
}
