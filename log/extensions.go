package log

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	// TraceIDHeader = "X-amzn-Trace-ID".
	TraceIDHeader = "X-amzn-Trace-ID"
	// RequestIDHeader = "X-Request-ID".
	RequestIDHeader = "X-Request-ID"
	// CorrelationIDHeader = "X-Correlation-ID".
	CorrelationIDHeader = "X-Correlation-ID"
	// ErrorUUID = "00000000-0000-0000-0000-000000000000".
	ErrorUUID = "00000000-0000-0000-0000-000000000000"
)

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

// WithRequestTracing.
func (lf *LoggerField) WithRequestTracing(req *http.Request) *LoggerField {
	if req == nil {
		return lf
	}

	traceID := req.Header.Get(TraceIDHeader)
	requestID := req.Header.Get(RequestIDHeader)
	correlationID := req.Header.Get(CorrelationIDHeader)

	return lf.doc("tracing", SubDoc().
		Str("trace_id", traceID).
		Str("request_id", requestID).
		Str("correlation_id", correlationID),
	)
}

func (lf *LoggerField) WithAuthenticatedUserTracing(auth *AuthPayload) *LoggerField {
	if auth == nil {
		return lf
	}

	return lf.doc("authentication", SubDoc().
		Str("account_id", auth.CustomerAccountID).
		Str("realuser_id", auth.RealUserID).
		Str("user_id", auth.UserID),
	)
}

func (lf *LoggerField) WithSystemTracing() *LoggerField {
	host, _ := os.Hostname()
	_, path, line, ok := runtime.Caller(1)
	file := "unknown"
	if ok {
		file = filepath.Base(path)
	}

	return lf.doc("system", SubDoc().
		Str("os", runtime.GOOS).
		Int("num_cpu", runtime.NumCPU()).
		Str("host", host).
		Str("loc", file+":"+strconv.Itoa(line)),
	)
}
