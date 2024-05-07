package log

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	TraceIDHeader                        = "X-amzn-Trace-ID"
	RequestIDHeader                      = "X-Request-ID"
	CorrelationIDHeader                  = "X-Correlation-ID"
	ErrorUUID                            = "00000000-0000-0000-0000-000000000000"
	AuthorizationHeader                  = "Authorization"
	XCAServiceGatewayAuthorizationHeader = "X-CA-SGW-Authorization"
	UserAgentHeader                      = "User-Agent"
	XForwardedForHeader                  = "X-Forwarded-For"
)

// AuthPayload contains the customer account_id, user_id and realuser_id uuids.
type AuthPayload struct {
	// CustomerAccountID is the ID of the currently logged in user's parent
	// account/organization, sometimes known as the "account_aggregate_id".
	CustomerAccountID string
	// UserID is the ID of the currently authenticated user, and will
	// generally be a "user_aggregate_id".
	UserID string
	// RealUserID, when supplied, is the ID of the user who is currently
	// impersonating the current "UserID". This value is optional.
	RealUserID string
}

// WithRequestTracing adds a "tracing" subdocument to the log that
// includes important trace, request and correlation fields.
func (lf *Property) WithRequestTracing(req *http.Request) *Property {
	if req == nil {
		return lf
	}

	props := requestTracingFields(req)
	return lf.doc("tracing", props)
}

// WithRequestDiagnostics adds a "request" subdocument to the log that
// includes important request fields.
func (lf *Property) WithRequestDiagnostics(req *http.Request) *Property {
	if req == nil {
		return lf
	}

	props := requestDiagnosticsFields(req)
	return lf.doc("request", props)
}

// WithAuthenticatedUserTracing adds an "authentication" subdocument to the log that
// includes important account, user and realuser fields.
func (lf *Property) WithAuthenticatedUserTracing(auth *AuthPayload) *Property {
	if auth == nil {
		return lf
	}

	props := authenticatedUserTracingFields(auth)
	return lf.doc("authentication", props)
}

// WithAuthorizationTracing adds an "authorization" subdocument to the log that
// includes important authorization headers that are automatically redacted.
func (lf *Property) WithAuthorizationTracing(req *http.Request) *Property {
	if req == nil {
		return lf
	}

	props := authorizationTracingFields(req)
	return lf.doc("authorization", props)
}

// WithDatadogTracing adds a "datadog" subdocument to the log that
// includes the fields dd.trace_id and dd.span_id. If Xray is configured it also
// adds xray.trace_id and xray.seg_id fields.
func (lf *Property) WithDatadogTracing(ctx context.Context) *Property {
	if ctx == nil {
		return lf
	}

	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		lf.impl = lf.impl.
			Uint64("dd.trace_id", span.Context().TraceID()).
			Uint64("dd.span_id", span.Context().SpanID())
	}

	seg := xray.GetSegment(ctx)
	if seg != nil {
		lf.impl = lf.impl.
			Str("xray.trace_id", seg.TraceID).
			Str("xray.seg_id", seg.ID)
	}

	return lf
}

// WithSystemTracing adds a "system" subdocument to the log that
// includes important host, runtime, cpu and loc fields.
func (lf *Property) WithSystemTracing() *Property {
	props := systemTracingFields()
	return lf.doc("system", props)
}
