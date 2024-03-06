package log

import (
	"context"
)

// EventCtxKey type
type EventCtxKey int

const (
	// RequestFieldsCtx EventCtxKey = iota
	RequestFieldsCtx EventCtxKey = iota
)

type RequestScopedFields struct {
	TraceID             string `json:"trace_id"`       // AWS XRAY trace id. Format of this is controlled by AWS. Do not rely on it, some services may not use XRAY.
	RequestID           string `json:"request_id"`     // Client generated RANDOM string. Most of the time this will be empty. Clients can set this to help us diagnose issues.
	CorrelationID       string `json:"correlation_id"` // Set ALWAYS by the web-gateway as a UUID v4.
	UserAggregateID     string `json:"user"`           // If JWT and correct key present, then this will be set to the Effective User UUID
	CustomerAggregateID string `json:"customer"`       // If JWT and correct key present, then this will be set to the Customer UUID (aka Account)
}

// AddRequestFields adds a RequestScopedFields to the context
func AddRequestFields(ctx context.Context, rsFields RequestScopedFields) context.Context {
	return context.WithValue(ctx, RequestFieldsCtx, rsFields)
}

// GetRequestScopedFields gets the RequestScopedFields from the context
func GetRequestScopedFields(ctx context.Context) (RequestScopedFields, bool) {
	if ctx == nil {
		return RequestScopedFields{}, false
	}

	rsFields, ok := ctx.Value(RequestFieldsCtx).(RequestScopedFields)
	return rsFields, ok
}

func (lf *Property) WithGlamplifyRequestFieldsFromCtx(ctx context.Context) *Property {
	rsFields, ok := GetRequestScopedFields(ctx)
	if !ok {
		return lf
	}

	return lf.doc("authentication", SubDoc().
		Str("account_id", rsFields.CustomerAggregateID).
		Str("user_id", rsFields.UserAggregateID),
	).doc("tracing", SubDoc().
		Str("trace_id", rsFields.TraceID).
		Str("request_id", rsFields.RequestID).
		Str("correlation_id", rsFields.CorrelationID),
	)
}
