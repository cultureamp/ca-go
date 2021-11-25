package request

import "context"

type contextValueKey string

const requestIDsKey = contextValueKey("fields")

// RequestIDs represent the set of unique identifiers for a request.
type RequestIDs struct {
	RequestID     string
	CorrelationID string
}

// AddRequestIDsToContext returns a new context with the given RequestIDs
// embedded as a value.
func ContextWithRequestIDs(ctx context.Context, fields RequestIDs) context.Context {
	return context.WithValue(ctx, requestIDsKey, fields)
}

func RequestIDsFromContext(ctx context.Context) (RequestIDs, bool) {
	ids, ok := ctx.Value(requestIDsKey).(RequestIDs)
	return ids, ok
}

func ContextHasRequestIDs(ctx context.Context) bool {
	_, ok := RequestIDsFromContext(ctx)
	return ok
}
