package log_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/log"
)

func ExampleLogger_Info_withRequestTracing() {
	config := getExampleLoggerConfig("INFO")
	logger := log.NewLogger(config)

	// First test nil Request
	logger.Info("info_with_nil_request_tracing").
		WithRequestTracing(nil).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain request tracing")

	// Next with Request but no headers
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	logger.Info("info_with_missing_headers_request_tracing").
		WithRequestTracing(req).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log empty request tracing")

	// Finally with headers set
	req.Header.Add(log.TraceIDHeader, "trace_123_id")
	req.Header.Add(log.RequestIDHeader, "request_456_id")
	req.Header.Add(log.CorrelationIDHeader, "correlation_789_id")

	logger.Info("info_with_request_tracing").
		WithRequestTracing(req).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain request tracing")

	// Output:
	// 2020-11-14T11:30:32Z INF event="logging should not contain request tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_request_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-11-14T11:30:32Z INF event="logging should log empty request tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_missing_headers_request_tracing farm=local product= properties={"resource":"resource_id","test-number":1} tracing={"correlation_id":"","request_id":"","trace_id":""}
	// 2020-11-14T11:30:32Z INF event="logging should contain request tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_request_tracing farm=local product= properties={"resource":"resource_id","test-number":1} tracing={"correlation_id":"correlation_789_id","request_id":"request_456_id","trace_id":"trace_123_id"}
}

func ExampleLogger_Info_withRequestDiagnostics() {
	config := getExampleLoggerConfig("INFO")
	logger := log.NewLogger(config)

	// First test nil Request
	logger.Info("info_with_nil_request_diagnostics").
		WithRequestDiagnostics(nil).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain request diagnostics")

	// Next with Request but empty url
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	logger.Info("info_with_minimal_url_request_diagnostics").
		WithRequestDiagnostics(req).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log minimal request diagnostics")

	// Next with Request but with url
	req = httptest.NewRequest(http.MethodGet, "https://example.com/foo?arg=value#fragment-1", nil)

	logger.Info("info_with_request_diagnostics").
		WithRequestDiagnostics(req).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain request diagnostics")

	// Output:
	// 2020-11-14T11:30:32Z INF event="logging should not contain request diagnostics" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_request_diagnostics farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-11-14T11:30:32Z INF event="logging should log minimal request diagnostics" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_minimal_url_request_diagnostics farm=local product= properties={"resource":"resource_id","test-number":1} request={"fragment":"","host":"example.com","method":"GET","path":"","proto":"HTTP/1.1","query":"","scheme":"http"}
	// 2020-11-14T11:30:32Z INF event="logging should contain request diagnostics" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_request_diagnostics farm=local product= properties={"resource":"resource_id","test-number":1} request={"fragment":"","host":"example.com","method":"GET","path":"/foo","proto":"HTTP/1.1","query":"arg=value#fragment-1","scheme":"https"}
}

func ExampleLogger_Info_withAuthenticationUserTracing() {
	config := getExampleLoggerConfig("INFO")
	logger := log.NewLogger(config)

	// First test nil Auth Payload
	logger.Info("info_with_nil_authN_tracing").
		WithAuthenticatedUserTracing(nil).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain authN tracing")

	// Next with empty Auth Payload
	auth := &log.AuthPayload{}

	logger.Info("info_with_missing_authN_tracing").
		WithAuthenticatedUserTracing(auth).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log empty authN tracing")

	// Finally with Auth Payload set
	auth = &log.AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	logger.Info("info_with_authN_tracing").
		WithAuthenticatedUserTracing(auth).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain authN tracing")

	// Output:
	// 2020-11-14T11:30:32Z INF event="logging should not contain authN tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_auth_n_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-11-14T11:30:32Z INF event="logging should log empty authN tracing" app= app_version=1.0.0 authentication={"account_id":"","realuser_id":"","user_id":""} aws_account_id=development aws_region= event=info_with_missing_auth_n_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-11-14T11:30:32Z INF event="logging should contain authN tracing" app= app_version=1.0.0 authentication={"account_id":"account_123_id","realuser_id":"real_456_id","user_id":"user_789_id"} aws_account_id=development aws_region= event=info_with_auth_n_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
}

func ExampleLogger_Info_withAuthorizationTracing() {
	config := getExampleLoggerConfig("INFO")
	logger := log.NewLogger(config)

	// First test nil Request
	logger.Info("info_with_nil_authZ_tracing").
		WithAuthorizationTracing(nil).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain authZ tracing")

	// Next with Request but no headers
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	logger.Info("info_with_missing_headers_authZ_tracing").
		WithAuthorizationTracing(req).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log empty authZ tracing")

	// Finally with headers set
	req.Header.Add(log.AuthorizationHeader, "AWS 123 token")
	req.Header.Add(log.XCAServiceGatewayAuthorizationHeader, "Bearer 456 token")
	req.Header.Add(log.XForwardedForHeader, "123.123.123")
	req.Header.Add(log.UserAgentHeader, "node")

	logger.Info("info_with_authZ_tracing").
		WithAuthorizationTracing(req).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain authZ tracing")

	// Output:
	// 2020-11-14T11:30:32Z INF event="logging should not contain authZ tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_auth_z_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-11-14T11:30:32Z INF event="logging should log empty authZ tracing" app= app_version=1.0.0 authorization={"authorization_token":"","user_agent":"","x_forwarded_for":"","xca_service_authorization_token":""} aws_account_id=development aws_region= event=info_with_missing_headers_auth_z_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-11-14T11:30:32Z INF event="logging should contain authZ tracing" app= app_version=1.0.0 authorization={"authorization_token":"AWS**********ken","user_agent":"node","x_forwarded_for":"123.123.123","xca_service_authorization_token":"Bear**********oken"} aws_account_id=development aws_region= event=info_with_auth_z_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
}

func TestExtensionWithSystemTracing(t *testing.T) {
	config := getExampleLoggerConfig("INFO")
	config.Quiet = true
	logger := log.NewLogger(config)

	logger.Info("info_with_nil_auth_tracing").
		WithSystemTracing().
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain system tracing")

	// System tracing add a "pid" "num_cpus" etc. which changes from run to run / machine to machine
	// So we don't use a testable Example()
}

func getExampleLoggerConfig(sev string) *log.Config {
	config := log.NewLoggerConfig()
	config.LogLevel = sev
	config.Quiet = false
	config.ConsoleWriter = true
	config.ConsoleColour = false
	config.TimeNow = func() time.Time {
		return time.Date(2020, 11, 14, 11, 30, 32, 0, time.UTC)
	}
	return config
}
