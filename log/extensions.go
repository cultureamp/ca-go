package log

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/segmentio/kafka-go"

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

	traceID := req.Header.Get(TraceIDHeader)
	requestID := req.Header.Get(RequestIDHeader)
	correlationID := req.Header.Get(CorrelationIDHeader)

	return lf.doc("tracing", Add().
		Str("trace_id", traceID).
		Str("request_id", requestID).
		Str("correlation_id", correlationID),
	)
}

// WithRequestDiagnostics adds a "request" subdocument to the log that
// includes important request fields.
func (lf *Property) WithRequestDiagnostics(req *http.Request) *Property {
	if req == nil {
		return lf
	}

	url := req.URL

	return lf.doc("request", Add().
		Str("method", req.Method).
		Str("proto", req.Proto).
		Str("host", req.Host).
		Str("scheme", url.Scheme).
		Str("path", url.Path).
		Str("query", url.RawQuery).
		Str("fragment", url.Fragment),
	)
}

// WithAuthenticatedUserTracing adds an "authentication" subdocument to the log that
// includes important account, user and realuser fields.
func (lf *Property) WithAuthenticatedUserTracing(auth *AuthPayload) *Property {
	if auth == nil {
		return lf
	}

	return lf.doc("authentication", Add().
		Str("account_id", auth.CustomerAccountID).
		Str("realuser_id", auth.RealUserID).
		Str("user_id", auth.UserID),
	)
}

// WithAuthorizationTracing adds an "authorization" subdocument to the log that
// includes important authorization headers that are automatically redacted.
func (lf *Property) WithAuthorizationTracing(req *http.Request) *Property {
	if req == nil {
		return lf
	}

	auth_token := req.Header.Get(AuthorizationHeader)
	xca_auth_token := req.Header.Get(XCAServiceGatewayAuthorizationHeader)
	user_agent := req.Header.Get(UserAgentHeader)
	forward_for := req.Header.Get(XForwardedForHeader)

	return lf.doc("authorization", Add().
		Str("authorization_token", redactString(auth_token)).
		Str("xca_service_authorization_token", redactString(xca_auth_token)).
		Str("user_agent", user_agent).
		Str("x_forwarded_for", forward_for),
	)
}

// WithDatadogTracing adds a "datadog" subdocument to the log that
// includes the spans dd.trace_id and dd.span_id. If Xray is configured it also
// adds xray.trace_id and xray.seg_id fields.
func (lf *Property) WithDatadogTracing(ctx context.Context) *Property {
	if ctx == nil {
		return lf
	}

	props := Add()

	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		props = props.
			UInt64("dd.trace_id", span.Context().TraceID()).
			UInt64("dd.span_id", span.Context().SpanID())
	}

	seg := xray.GetSegment(ctx)
	if seg != nil {
		props = props.
			Str("xray.trace_id", seg.TraceID).
			Str("xray.seg_id", seg.ID)
	}

	if !ok && seg == nil {
		return lf
	}

	return lf.doc("datadog", props)
}

// WithKafkaTracing adds a "kafka" subdocument to the log that
// includes topic, partition and offset.
func (lf *Property) WithKafkaTracing(msg *kafka.Message) *Property {
	if msg == nil {
		return lf
	}

	return lf.doc("kafka", Add().
		Str("topic", msg.Topic).
		Int("partition", msg.Partition).
		Int64("offset", msg.Offset),
	)
}

// WithSystemTracing adds a "system" subdocument to the log that
// includes important host, runtime, cpu and loc fields.
func (lf *Property) WithSystemTracing() *Property {
	host, _ := os.Hostname()
	_, path, line, ok := runtime.Caller(1)
	file := "unknown"
	if ok {
		file = filepath.Base(path)
	}
	buildInfo, _ := debug.ReadBuildInfo()

	return lf.doc("system", Add().
		Str("os", runtime.GOOS).
		Int("num_cpu", runtime.NumCPU()).
		Str("host", host).
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Str("loc", file+":"+strconv.Itoa(line)),
	)
}
