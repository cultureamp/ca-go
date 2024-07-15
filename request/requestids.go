package request

import "context"

type contextValueKey string

const uniqueIDsKey = contextValueKey("fields")

// UniqueIDs represent the set of unique identifiers for a request.
type UniqueIDs struct {
	RequestID     string
	CorrelationID string
}

// ContextWithUniqueIDs returns a new context with the given RequestIDs
// embedded as a value.
func ContextWithUniqueIDs(ctx context.Context, fields UniqueIDs) context.Context {
	return context.WithValue(ctx, uniqueIDsKey, fields)
}

// UniqueIDsFromContext attempts to retrieve a RequestIDs struct from the given
// context, returning a RequestIDs struct along with a boolean signalling
// whether the retrieval was successful.
func UniqueIDsFromContext(ctx context.Context) (UniqueIDs, bool) {
	ids, ok := ctx.Value(uniqueIDsKey).(UniqueIDs)
	return ids, ok
}

// ContextHasUniqueIDs returns whether the given context contains a RequestIDs
// value.
func ContextHasUniqueIDs(ctx context.Context) bool {
	_, ok := UniqueIDsFromContext(ctx)
	return ok
}
