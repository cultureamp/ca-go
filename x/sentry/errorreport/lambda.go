package errorreport

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

// TODO: we should try to upgrade to Go 1.18 and use Generics here instead of interface{}
type EventHandler func(context.Context, interface{}) error

// A LambdaHandler is a middleware factory that provides integration with Sentry.
type LambdaHandler struct {
	repanic bool
	timeout time.Duration
}

// PanicOptions configure a LambdaHandler.
type PanicOptions struct {
	// Repanic configures whether to panic again after recovering from a panic.
	// Use this option if you have other panic handlers or want the default
	// behavior from AWS lambda runtime.
	Repanic bool
	// Timeout for the delivery of panic events. Defaults to 2s.
	//
	// If the timeout is reached, the current goroutine is no longer blocked
	// waiting, but the delivery is not canceled.
	Timeout time.Duration
}

// New returns a new LambdaHandler. Use the Handle and HandleFunc methods to wrap
// existing handlers.
func New(options PanicOptions) *LambdaHandler {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	return &LambdaHandler{
		repanic: options.Repanic,
		timeout: timeout,
	}
}

// Handle works as a middleware that wraps an existing KinesisEventHandler. A wrapped
// handler will recover from and report panics to Sentry, and provide access to
// a request-specific hub to report messages and errors.
func (h *LambdaHandler) Handle(handler EventHandler) EventHandler {
	return h.handle(handler)
}

func (h *LambdaHandler) handle(handler EventHandler) EventHandler {
	return func(ctx context.Context, event interface{}) error {
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}
		defer h.recoverWithSentry(ctx, hub)

		err := handler(ctx, event)
		if err != nil {
			ReportError(ctx, err)
		}
		return err
	}
}

func (h *LambdaHandler) recoverWithSentry(ctx context.Context, hub *sentry.Hub) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(ctx, err)
		if eventID != nil {
			hub.Flush(h.timeout)
		}
		if h.repanic {
			panic(err)
		}
	}
}
