package errorreport

import (
	"context"
	"fmt"

	"github.com/cultureamp/ca-go/x/lambdafunction"
	"github.com/getsentry/sentry-go"
)

// LambdaErrorOptions configures the way Sentry is used in the context of a
// Lambda handler wrapper.
type LambdaErrorOptions struct {
	// Repanic configures whether to panic again after recovering from a panic.
	// Use this option if you have other panic handlers or want the default
	// behavior from AWS lambda runtime. Defaults to true.
	Repanic *bool
}

// LambdaMiddleware[TIn] provides error-handling middleware for a Lambda
// function that has a payload type of TIn. This suits Lambda functions like
// event processors, where the return has no payload.
func LambdaMiddleware[TIn any](nextHandler lambdafunction.HandlerOf[TIn], options LambdaErrorOptions) lambdafunction.HandlerOf[TIn] {
	return func(ctx context.Context, event TIn) error {
		defer beforeHandler(ctx, options)()

		err := nextHandler(ctx, event)

		afterHandler(ctx, err)

		return err
	}
}

// LambdaWithOutputMiddleware[TIn, TOut] provides error-handling middleware for
// a Lambda function that has a payload type of TIn and returns the tuple TOut,error.
func LambdaWithOutputMiddleware[TIn any, TOut any](nextHandler lambdafunction.HandlerWithOutputOf[TIn, TOut], options LambdaErrorOptions) lambdafunction.HandlerWithOutputOf[TIn, TOut] {
	return func(ctx context.Context, event TIn) (TOut, error) {
		defer beforeHandler(ctx, options)()
		fmt.Println("afterbefore")

		out, err := nextHandler(ctx, event)

		fmt.Println("beforeafter")
		afterHandler(ctx, err)
		fmt.Println("afterafter")

		return out, err
	}
}

func beforeHandler(ctx context.Context, options LambdaErrorOptions) func() {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}
	return func() {
		if err := recover(); err != nil {
			fmt.Println("panic found!")
			_ = hub.RecoverWithContext(ctx, err)

			if options.Repanic == nil || *options.Repanic {
				fmt.Println("repanic!")
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
