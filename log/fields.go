package log

import (
	"net/http"
)

// Deprecated: Fields is a drop in replacement for glamplify log.Fields in logging statements.
type Fields map[string]interface{}

// Deprecated: WithRequestTracing added a "tracing" subdocument to the log that
// include important trace, request and correlation headers.
func (f Fields) WithRequestTracing(req *http.Request) Fields {
	if req == nil {
		return f
	}

	traceID := req.Header.Get(TraceIDHeader)
	requestID := req.Header.Get(RequestIDHeader)
	correlationID := req.Header.Get(CorrelationIDHeader)

	f["system"] = Fields{
		"trace_id":       traceID,
		"request_id":     requestID,
		"correlation_id": correlationID,
	}
	return f
}

// Deprecated: WithAuthenticatedUserTracing added a "authentication" subdocument to the log that
// include important account, user and realuser fields.
func (f Fields) WithAuthenticatedUserTracing(auth *AuthPayload) Fields {
	if auth == nil {
		return f
	}

	f["authentication"] = Fields{
		"account_id":  auth.CustomerAccountID,
		"realuser_id": auth.RealUserID,
		"user_id":     auth.UserID,
	}
	return f
}

// Deprecated: Merge combines multiple legacy log.Fields together.
func (fields Fields) Merge(other ...Fields) Fields {
	merged := Fields{}

	for k, v := range fields {
		merged[k] = v
	}

	for _, f := range other {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}
