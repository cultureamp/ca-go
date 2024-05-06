package log

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/rs/zerolog"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type LoggerOption func(zerolog.Context) zerolog.Context

func WithProperties(fields *Field) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Dict("properties", fields.impl)
	}
}

// WithRequestTracing adds a "tracing" subdocument to the log that
// includes important trace, request and correlation fields.
func WithRequestTracing(req *http.Request) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		if req == nil {
			return lc
		}

		props := requestTracingFields(req)
		return lc.Dict("tracing", props.impl)
	}
}

// WithRequestDiagnostics adds a "request" subdocument to the log that
// includes important request fields.
func WithRequestDiagnostics(req *http.Request) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		if req == nil {
			return lc
		}

		props := requestDiagnosticsFields(req)
		return lc.Dict("request", props.impl)
	}
}

// WithAuthenticatedUserTracing adds an "authentication" subdocument to the log that
// includes important account, user and realuser fields.
func WithAuthenticatedUserTracing(auth *AuthPayload) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		if auth == nil {
			return lc
		}

		props := authenticatedUserTracingFields(auth)
		return lc.Dict("authentication", props.impl)
	}
}

// WithAuthorizationTracing adds an "authorization" subdocument to the log that
// includes important authorization headers that are automatically redacted.
func WithAuthorizationTracing(req *http.Request) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		if req == nil {
			return lc
		}

		props := authorizationTracingFields(req)
		return lc.Dict("authorization", props.impl)
	}
}

// WithDatadogTracing adds a "datadog" subdocument to the log that
// includes the fields dd.trace_id and dd.span_id. If Xray is configured it also
// adds xray.trace_id and xray.seg_id fields.
func WithDatadogTracing(ctx context.Context) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		if ctx == nil {
			return lc
		}

		span, ok := tracer.SpanFromContext(ctx)
		if ok {
			lc = lc.
				Uint64("dd.trace_id", span.Context().TraceID()).
				Uint64("dd.span_id", span.Context().SpanID())
		}

		seg := xray.GetSegment(ctx)
		if seg != nil {
			lc = lc.
				Str("xray.trace_id", seg.TraceID).
				Str("xray.seg_id", seg.ID)
		}

		return lc
	}
}

// WithSystemTracing adds a "system" subdocument to the log that
// includes important host, runtime, cpu and loc fields.
func WithSystemTracing() LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		props := systemTracingFields()
		return lc.Dict("system", props.impl)
	}
}
