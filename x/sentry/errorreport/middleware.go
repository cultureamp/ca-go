package errorreport

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"goa.design/goa"
)

type OnRequestPanicHandler func(context.Context, http.ResponseWriter, error)

func NewHTTPMiddleware(onRequestPanic OnRequestPanicHandler) func(http.Handler) http.Handler {
	sentryWrapper := sentryhttp.New(sentryhttp.Options{
		// Repanic to propagate the error to the onRequestPanic handler.
		Repanic: true,
	})

	return func(next http.Handler) http.Handler {
		// Wrap downstream HTTP handlers with the repanic capability.
		sentryHandler := sentryWrapper.Handle(next)

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer recoverRequestPanic(req.Context(), w, onRequestPanic)

			sentryHandler.ServeHTTP(w, req)
		})
	}
}

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
