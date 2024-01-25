package log

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtensionWithRequestTracing(t *testing.T) {
	config := NewLoggerConfig()
	config.Quiet = false
	logger := NewLogger(config)

	// First test nil Request
	logger.Info("info_with_nil_request_tracing").
		WithRequestTracing(nil).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain request tracing")

	// Next with Request but no headers
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	logger.Info("info_with_missing_headers_request_tracing").
		WithRequestTracing(req).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log empty request tracing")

	// Finally with headers set
	req.Header.Add(TraceIDHeader, "trace_123_id")
	req.Header.Add(RequestIDHeader, "request_456_id")
	req.Header.Add(CorrelationIDHeader, "correlation_789_id")

	logger.Info("info_with_request_tracing").
		WithRequestTracing(req).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain request tracing")

	// Output:
	// "severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_with_nil_request_tracing","properties":{"resource":"resource_id","test-number":1},"time":"2024-01-17T15:38:13+11:00","details":"logging should not contain request tracing"}
	// {"severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_with_missing_headers_request_tracing","tracing":{"trace_id":"","request_id":"","correlation_id":""},"properties":{"resource":"resource_id","test-number":1},"time":"2024-01-17T15:38:13+11:00","details":"logging should not contain request tracing"}
	// {"severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_with_request_tracing","tracing":{"trace_id":"trace_123_id","request_id":"request_456_id","correlation_id":"correlation_789_id"},"properties":{"resource":"resource_id","test-number":1},"time":"2024-01-17T15:38:13+11:00","details":"logging should contain request tracing"}
}

func TestExtensionWithAuthenticationUserTracing(t *testing.T) {
	config := NewLoggerConfig()
	config.Quiet = false
	logger := NewLogger(config)

	// First test nil Auth Payload
	logger.Info("info_with_nil_auth_tracing").
		WithAuthenticatedUserTracing(nil).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain auth tracing")

	// Next with empty Auth Payload
	auth := &AuthPayload{}

	logger.Info("info_with_missing_auth_tracing").
		WithAuthenticatedUserTracing(auth).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log empty auth tracing")

	// Finally with Auth Payload set
	auth = &AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	logger.Info("info_with_auth_tracing").
		WithAuthenticatedUserTracing(auth).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain auth tracing")

	// Output:
	//{"severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_with_nil_auth_tracing","properties":{"resource":"resource_id","test-number":1},"time":"2024-01-17T15:38:13+11:00","details":"logging should not contain auth tracing"}
	// {"severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_with_missing_auth_tracing","authentication":{"account_id":"","realuser_id":"","user_id":""},"properties":{"resource":"resource_id","test-number":1},"time":"2024-01-17T15:38:13+11:00","details":"logging should not contain auth tracing"}
	// {"severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_with_auth_tracing","authentication":{"account_id":"account_123_id","realuser_id":"real_456_id","user_id":"user_789_id"},"properties":{"resource":"resource_id","test-number":1},"time":"2024-01-17T15:38:13+11:00","details":"logging should contain auth tracing"}
}

func TestExtensionWithSystemTracing(t *testing.T) {
	config := NewLoggerConfig()
	config.Quiet = false
	logger := NewLogger(config)

	logger.Info("info_with_nil_auth_tracing").
		WithSystemTracing().
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain system tracing")

	// Output:
	// {"severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_with_nil_auth_tracing","system":{"os":"darwin","num_cpu":8,"host":"mridgway-6RR4DK","loc":"/opt/homebrew/opt/go/libexec/src/runtime/asm_arm64.s:1197"},"properties":{"resource":"resource_id","test-number":1},"time":"2024-01-17T15:38:13+11:00","details":"logging should contain system tracing"}
}
