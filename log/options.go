package log

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	iso8601 "github.com/sosodev/duration"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type LoggerOption func(zerolog.Context) zerolog.Context

// Str adds the property key with val as a string to the log.
// Note: Empty string values will not be logged.
func WithStr(key string, val string) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		if val == "" {
			return lc
		}
		return lc.Str(key, val)
	}
}

// Int adds the property key with val as an int to the log.
func WithInt(key string, val int) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Int(key, val)
	}
}

// UInt adds the property key with val as an int to the log.
func WithUInt(key string, val uint) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Uint(key, val)
	}
}

// Int64 adds the property key with val as an int64 to the log.
func WithInt64(key string, val int64) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Int64(key, val)
	}
}

// UInt64 adds the property key with val as an uint64 to the log.
func WithUInt64(key string, val uint64) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Uint64(key, val)
	}
}

// Float32 adds the property key with val as an float32 to the log.
func WithFloat32(key string, val float32) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Float32(key, val)
	}
}

// Float64 adds the property key with val as an float64 to the log.
func WithFloat64(key string, val float64) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Float64(key, val)
	}
}

// Bool adds the property key with b as an bool to the log.
func WithBool(key string, b bool) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Bool(key, b)
	}
}

// Bytes adds the property key with val as an []byte to the log.
func WithBytes(key string, val []byte) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Bytes(key, val)
	}
}

// Duration adds the property key with val as an time.Duration to the log.
func WithDuration(key string, d time.Duration) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		// Logging Std https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard#Custom-fields
		// time durations use ISO8601 Duration format.
		s := iso8601.Format(d)
		return lc.Str(key, s)
	}
}

// Time adds the property key with val as an uuid.UUID to the log.
func WithTime(key string, t time.Time) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		// uses zerolog.TimeFieldFormat which we set to time.RFC3339
		return lc.Time(key, t)
	}
}

// IPAddr adds the property key with val as an net.IP to the log.
func WithIPAddr(key string, ip net.IP) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.IPAddr(key, ip)
	}
}

// UUID adds the property key with val as an uuid.UUID to the log.
func WithUUID(key string, uuid uuid.UUID) LoggerOption {
	return func(lc zerolog.Context) zerolog.Context {
		return lc.Str(key, uuid.String())
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
