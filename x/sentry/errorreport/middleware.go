package errorreport

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	goa "goa.design/goa/v3/pkg"
)

// OnRequestPanicHandler is a function that can be supplied to HTTP middleware
// to perform further processing of an HTTP request after a panic has occurred.
type OnRequestPanicHandler func(context.Context, http.ResponseWriter, error)

// defaultRequestPanicHandler writes a JSON:API style error response with a
// 500 status code.
func defaultRequestPanicHandler(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(`{"errors":[{"status":"500","title":"Internal Server Error"}]}`))
}

// NewHTTPMiddleware returns an http.Handler that reports panics to Sentry, recovers
// from the panic, and calls the OnRequestPanicHandler if provided. If a handler is
// not provided, returns a JSON:API structured body with status 500.
func NewHTTPMiddleware(onRequestPanic OnRequestPanicHandler) func(http.Handler) http.Handler {
	panicHandler := defaultRequestPanicHandler
	if onRequestPanic != nil {
		panicHandler = onRequestPanic
	}

	sentryWrapper := sentryhttp.New(sentryhttp.Options{
		// Repanic to propagate the error to the onRequestPanic handler.
		Repanic: true,
	})

	// Returns a handler that configures Sentry to report the panic, and then re-panic.
	//
	// The new panic is recovered to enable the onRequestPanic function to be executed so
	// it can perform further manipulation of the HTTP response before it's sent
	// to the client.
	//
	// The most common action on panic will be to change the HTTP response code
	// and request body to reflect an error occurring.
	return func(next http.Handler) http.Handler {
		sentryHandler := sentryWrapper.Handle(next)

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer recoverRequestPanic(req.Context(), w, panicHandler)

			scope := sentry.CurrentHub().PushScope()
			addRequestFieldsToScope(req.Context(), scope)
			defer sentry.PopScope()

			sentryHandler.ServeHTTP(w, req)
		})
	}
}

// NewGoaEndpointMiddleware returns Goa middleware to detect and report
// errors to Sentry.
func NewGoaEndpointMiddleware() func(goa.Endpoint) goa.Endpoint {
	return func(next goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			res, err := next(ctx, request)
			if err != nil {
				ReportError(ctx, err)
			}

			return res, err
		}
	}
}

func recoverRequestPanic(ctx context.Context, w http.ResponseWriter, errorHandler func(context.Context, http.ResponseWriter, error)) {
	if r := recover(); r != nil {
		// convert to an error if it's not one already
		err, ok := r.(error)
		if !ok {
			err = errors.New(fmt.Sprint(r))
		}

		errorHandler(ctx, w, err)
	}
}
