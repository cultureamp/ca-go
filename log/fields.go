package log

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	iso8601 "github.com/sosodev/duration"
)

// Field contains an element of the log, usually a key-value pair.
type Field struct {
	impl *zerolog.Event
}

func newLoggerField(impl *zerolog.Event) *Field {
	return &Field{impl: impl}
}

// Add creates a new custom log properties list.
func Add() *Field {
	subDoc := zerolog.Dict()
	return newLoggerField(subDoc)
}

// Str adds the property key with val as a string to the log.
// Note: Empty string values will not be logged.
func (lf *Field) Str(key string, val string) *Field {
	if val == "" {
		return lf
	}

	lf.impl = lf.impl.Str(key, val)
	return lf
}

// Int adds the property key with val as an int to the log.
func (lf *Field) Int(key string, val int) *Field {
	lf.impl = lf.impl.Int(key, val)
	return lf
}

// UInt adds the property key with val as an int to the log.
func (lf *Field) UInt(key string, val uint) *Field {
	lf.impl = lf.impl.Uint(key, val)
	return lf
}

// Int64 adds the property key with val as an int64 to the log.
func (lf *Field) Int64(key string, val int64) *Field {
	lf.impl = lf.impl.Int64(key, val)
	return lf
}

// UInt64 adds the property key with val as an uint64 to the log.
func (lf *Field) UInt64(key string, val uint64) *Field {
	lf.impl = lf.impl.Uint64(key, val)
	return lf
}

// Float32 adds the property key with val as an float32 to the log.
func (lf *Field) Float32(key string, val float32) *Field {
	lf.impl = lf.impl.Float32(key, val)
	return lf
}

// Float64 adds the property key with val as an float64 to the log.
func (lf *Field) Float64(key string, val float64) *Field {
	lf.impl = lf.impl.Float64(key, val)
	return lf
}

// Bool adds the property key with b as an bool to the log.
func (lf *Field) Bool(key string, b bool) *Field {
	lf.impl = lf.impl.Bool(key, b)
	return lf
}

// Bytes adds the property key with val as an []byte to the log.
func (lf *Field) Bytes(key string, val []byte) *Field {
	lf.impl = lf.impl.Bytes(key, val)
	return lf
}

// Duration adds the property key with val as an time.Duration to the log.
func (lf *Field) Duration(key string, d time.Duration) *Field {
	// Logging Std https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard#Custom-fields
	// time durations use ISO8601 Duration format.
	s := iso8601.Format(d)
	lf.impl = lf.impl.Str(key, s)
	return lf
}

// Time adds the property key with val as an uuid.UUID to the log.
func (lf *Field) Time(key string, t time.Time) *Field {
	// uses zerolog.TimeFieldFormat which we set to time.RFC3339
	lf.impl = lf.impl.Time(key, t)
	return lf
}

// IPAddr adds the property key with val as an net.IP to the log.
func (lf *Field) IPAddr(key string, ip net.IP) *Field {
	lf.impl = lf.impl.IPAddr(key, ip)
	return lf
}

// UUID adds the property key with val as an uuid.UUID to the log.
func (lf *Field) UUID(key string, uuid uuid.UUID) *Field {
	lf.impl = lf.impl.Str(key, uuid.String())
	return lf
}

// Func allows an anonymous func to run for the accumulated event.
func (lf *Field) Func(f func(e *zerolog.Event)) *Field {
	lf.impl = lf.impl.Func(f)
	return lf
}

func requestTracingFields(req *http.Request) *Field {
	traceID := req.Header.Get(TraceIDHeader)
	requestID := req.Header.Get(RequestIDHeader)
	correlationID := req.Header.Get(CorrelationIDHeader)

	return Add().
		Str("trace_id", traceID).
		Str("request_id", requestID).
		Str("correlation_id", correlationID)
}

func requestDiagnosticsFields(req *http.Request) *Field {
	url := req.URL

	return Add().
		Str("method", req.Method).
		Str("proto", req.Proto).
		Str("host", req.Host).
		Str("scheme", url.Scheme).
		Str("path", url.Path).
		Str("query", url.RawQuery).
		Str("fragment", url.Fragment)
}

func authenticatedUserTracingFields(auth *AuthPayload) *Field {
	return Add().
		Str("account_id", auth.CustomerAccountID).
		Str("realuser_id", auth.RealUserID).
		Str("user_id", auth.UserID)
}

func authorizationTracingFields(req *http.Request) *Field {
	authToken := req.Header.Get(AuthorizationHeader)
	xcaAuthToken := req.Header.Get(XCAServiceGatewayAuthorizationHeader)
	userAgent := req.Header.Get(UserAgentHeader)
	forwardFor := req.Header.Get(XForwardedForHeader)

	return Add().
		Str("authorization_token", redactString(authToken)).
		Str("xca_service_authorization_token", redactString(xcaAuthToken)).
		Str("user_agent", userAgent).
		Str("x_forwarded_for", forwardFor)
}

func systemTracingFields() *Field {
	host, _ := os.Hostname()
	_, path, line, ok := runtime.Caller(1)
	file := "unknown"
	if ok {
		file = filepath.Base(path)
	}
	buildInfo, _ := debug.ReadBuildInfo()

	return Add().
		Str("os", runtime.GOOS).
		Int("num_cpu", runtime.NumCPU()).
		Str("host", host).
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Str("loc", file+":"+strconv.Itoa(line))
}
