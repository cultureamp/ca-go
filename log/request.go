package log

import (
	"context"
	"net/http"
)

const (
	TraceHeader       = "X-Amzn-Trace-Id"
	RequestHeader     = "X-Request-ID"
	CorrelationHeader = "X-Correlation-ID"
)

type contextRequestIDKey string

const requestIDsKey = contextRequestIDKey("request_ids")

// RequestIDs contains the request header values for X-Amzn-Trace-Id, X-Request-ID and X-Correlation-ID.
type RequestIDs struct {
	TraceID       string
	RequestID     string
	CorrelationID string
}

// ContextWithRequestIDs returns a new context with the given a http.Request
func ContextWithRequest(ctx context.Context, req *http.Request) context.Context {
	ids := getRequestIDs(req)
	return ContextWithRequestIDs(ctx, ids)
}

// ContextWithRequestIDs returns a new context with the given RequestIDs
// embedded as a value.
func ContextWithRequestIDs(ctx context.Context, ids RequestIDs) context.Context {
	return context.WithValue(ctx, requestIDsKey, ids)
}

// RequestIDsFromContext attempts to retrieve a RequestIDs struct from the given
// context, returning a RequestIDs struct along with a boolean signalling
// whether the retrieval was successful.
func RequestIDsFromContext(ctx context.Context) (RequestIDs, bool) {
	ids, ok := ctx.Value(requestIDsKey).(RequestIDs)
	return ids, ok
}

func getRequestIDs(req *http.Request) RequestIDs {
	ids := RequestIDs{}

	if req == nil {
		return ids
	}

	ids.TraceID = req.Header.Get(TraceHeader)
	ids.RequestID = req.Header.Get(RequestHeader)
	ids.CorrelationID = req.Header.Get(CorrelationHeader)
	return ids
}
