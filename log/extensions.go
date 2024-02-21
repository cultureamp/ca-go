package log

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
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

	return lf.doc("tracing", SubDoc().
		Str("trace_id", traceID).
		Str("request_id", requestID).
		Str("correlation_id", correlationID),
	)
}

// WithAuthenticatedUserTracing adds an "authentication" subdocument to the log that
// includes important account, user and realuser fields.
func (lf *Property) WithAuthenticatedUserTracing(auth *AuthPayload) *Property {
	if auth == nil {
		return lf
	}

	return lf.doc("authentication", SubDoc().
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

	return lf.doc("authorization", SubDoc().
		Str("authorization_token", redactString(auth_token)).
		Str("xca_service_authorization_token", redactString(xca_auth_token)).
		Str("user_agent", user_agent).
		Str("x_forwarded_for", forward_for),
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

	return lf.doc("system", SubDoc().
		Str("os", runtime.GOOS).
		Int("num_cpu", runtime.NumCPU()).
		Str("host", host).
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Str("loc", file+":"+strconv.Itoa(line)),
	)
}

func redactString(s string) string {
	const minStars = 10
	const maxStars = 20

	l := len(s)
	if l == 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(l + maxStars)

	aQuarter := l / 4
	stars := minStars
	if stars < aQuarter {
		stars = maxStars
	}

	// write first quater of the chars if greater than minimum number of stars
	if l > minStars {
		b.WriteString(s[:aQuarter])
	}

	// no matter how long the string, show at least 10 "*" in the middle
	for i := 0; i < stars; i++ {
		b.WriteString("*")
	}

	// write remaining 1 or 2 chars
	if l > minStars {
		i := l - aQuarter
		b.WriteString(s[i:])
	}

	return b.String()
}
