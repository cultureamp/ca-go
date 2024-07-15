package request

import "context"

type contextValueKey string

const httpFieldIDsKey = contextValueKey("fields")

// HTTPFieldIDs represent the set of unique identifiers for a request.
type HTTPFieldIDs struct {
	RequestID     string
	CorrelationID string
}

// ContextWithHTTPFieldIDs returns a new context with the given RequestIDs
// embedded as a value.
func ContextWithHTTPFieldIDs(ctx context.Context, fields HTTPFieldIDs) context.Context {
	return context.WithValue(ctx, httpFieldIDsKey, fields)
}

// HTTPFieldIDsFromContext attempts to retrieve a RequestIDs struct from the given
// context, returning a RequestIDs struct along with a boolean signalling
// whether the retrieval was successful.
func HTTPFieldIDsFromContext(ctx context.Context) (HTTPFieldIDs, bool) {
	ids, ok := ctx.Value(httpFieldIDsKey).(HTTPFieldIDs)
	return ids, ok
}

// ContextHasHTTPFieldIDs returns whether the given context contains a RequestIDs
// value.
func ContextHasHTTPFieldIDs(ctx context.Context) bool {
	_, ok := HTTPFieldIDsFromContext(ctx)
	return ok
}
