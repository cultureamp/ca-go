package request

import "context"

type contextValueKey string

const requestIDsKey = contextValueKey("fields")

// RequestIDs represent the set of unique identifiers for a request.
type RequestIDs struct {
	RequestID     string
	CorrelationID string
}

// ContextWithRequestIDs returns a new context with the given RequestIDs
// embedded as a value.
func ContextWithRequestIDs(ctx context.Context, fields RequestIDs) context.Context {
	return context.WithValue(ctx, requestIDsKey, fields)
}

// RequestIDsFromContext attempts to retrieve a RequestIDs struct from the given
// context, returning a RequestIDs struct along with a boolean signalling
// whether the retrieval was successful.
func RequestIDsFromContext(ctx context.Context) (RequestIDs, bool) {
	ids, ok := ctx.Value(requestIDsKey).(RequestIDs)
	return ids, ok
}

// ContextHasRequestIDs returns whether the given context contains a RequestIDs
// value.
func ContextHasRequestIDs(ctx context.Context) bool {
	_, ok := RequestIDsFromContext(ctx)
	return ok
}
