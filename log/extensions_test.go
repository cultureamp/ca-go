package log

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// Local Console Output:
	// 2024-02-14T12:38:29+11:00 INF event="logging should not contain request tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_request_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2024-02-14T12:38:29+11:00 INF event="logging should log empty request tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_missing_headers_request_tracing farm=local product= properties={"resource":"resource_id","test-number":1} tracing={"correlation_id":"","request_id":"","trace_id":""}
	// 2024-02-14T12:38:29+11:00 INF event="logging should contain request tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_request_tracing farm=local product= properties={"resource":"resource_id","test-number":1} tracing={"correlation_id":"correlation_789_id","request_id":"request_456_id","trace_id":"trace_123_id"}
}

func TestExtensionWithAuthenticationUserTracing(t *testing.T) {
	config := NewLoggerConfig()
	config.Quiet = false
	logger := NewLogger(config)

	// First test nil Auth Payload
	logger.Info("info_with_nil_authN_tracing").
		WithAuthenticatedUserTracing(nil).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain authN tracing")

	// Next with empty Auth Payload
	auth := &AuthPayload{}

	logger.Info("info_with_missing_authN_tracing").
		WithAuthenticatedUserTracing(auth).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log empty authN tracing")

	// Finally with Auth Payload set
	auth = &AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	logger.Info("info_with_authN_tracing").
		WithAuthenticatedUserTracing(auth).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain authN tracing")

	// Local Console Output:
	// 2024-02-14T12:41:18+11:00 INF event="logging should not contain authN tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_auth_n_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2024-02-14T12:41:18+11:00 INF event="logging should log empty authN tracing" app= app_version=1.0.0 authentication={"account_id":"","realuser_id":"","user_id":""} aws_account_id=development aws_region= event=info_with_missing_auth_n_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2024-02-14T12:41:18+11:00 INF event="logging should contain authN tracing" app= app_version=1.0.0 authentication={"account_id":"account_123_id","realuser_id":"real_456_id","user_id":"user_789_id"} aws_account_id=development aws_region= event=info_with_auth_n_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
}

func TestExtensionWithAuthorizationTracing(t *testing.T) {
	config := NewLoggerConfig()
	config.Quiet = false
	logger := NewLogger(config)

	// First test nil Request
	logger.Info("info_with_nil_authZ_tracing").
		WithAuthorizationTracing(nil).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should not contain authZ tracing")

	// Next with Request but no headers
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	logger.Info("info_with_missing_headers_authZ_tracing").
		WithAuthorizationTracing(req).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should log empty authZ tracing")

	// Finally with headers set
	req.Header.Add(AuthorizationHeader, "AWS 123 token")
	req.Header.Add(XCAServiceGatewayAuthorizationHeader, "Bearer 456 token")
	req.Header.Add(XForwardedForHeader, "123.123.123")
	req.Header.Add(UserAgentHeader, "node")

	logger.Info("info_with_authZ_tracing").
		WithAuthorizationTracing(req).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain authZ tracing")

	// Local Console Output:
	// 2024-02-14T12:42:29+11:00 INF event="logging should not contain authZ tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_auth_z_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2024-02-14T12:42:29+11:00 INF event="logging should log empty authZ tracing" app= app_version=1.0.0 authorization={"authorization_token":"","user_agent":"","x_forwarded_for":"","xca_service_authorization_token":""} aws_account_id=development aws_region= event=info_with_missing_headers_auth_z_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2024-02-14T12:42:29+11:00 INF event="logging should contain authZ tracing" app= app_version=1.0.0 authorization={"authorization_token":"AWS**********ken","user_agent":"node","x_forwarded_for":"123.123.123","xca_service_authorization_token":"Bear**********oken"} aws_account_id=development aws_region= event=info_with_auth_z_tracing farm=local product= properties={"resource":"resource_id","test-number":1}
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

	// Local Console Output:
	// 2024-02-14T12:42:29+11:00 INF event="logging should contain system tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_nil_auth_tracing farm=local product= properties={"resource":"resource_id","test-number":1} system={"go_version":"go1.21.7","host":"mridgway-6RR4DK","loc":"extensions_test.go:143","num_cpu":8,"os":"darwin","pid":3721}
}

func TestStringRedaction(t *testing.T) {
	testCases := []struct {
		desc     string
		str      string
		redacted string
	}{
		{
			desc:     "Empty string returns empty string",
			str:      "",
			redacted: "",
		},
		{
			desc:     "String less than 10 chars shows 10 stars",
			str:      "1234",
			redacted: "**********",
		},
		{
			desc:     "String equals 10 chars shows 10 stars",
			str:      "1234567890",
			redacted: "**********",
		},
		{
			desc:     "String equals 11 chars shows first char and last chars and 10 stars",
			str:      "12345678901",
			redacted: "12**********01",
		},
		{
			desc:     "String equals 12 chars shows first and last chars with 10 stars in the middle",
			str:      "123456789012",
			redacted: "123**********012",
		},
		{
			desc:     "String equals 20 chars shows first and last chars with 10 stars in the middle",
			str:      "12345678901234567890",
			redacted: "12345**********67890",
		},
		{
			desc:     "String equals 30 chars shows first and last chars with 10 stars in the middle",
			str:      "123456789012345678901234567890",
			redacted: "1234567**********4567890",
		},
		{
			desc:     "Real world test",
			str:      "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJkNDc1ZGQ1Yi1mMTZjLTRiZmItODk4Yy1kMzQzNWEyMTUyMzkiLCJlZmZlY3RpdmVVc2VySWQiOiI1YjMxNjY0YS03NjEwLTRmYjAtYmM4OS1mOWY4ZTIwYmY4Y2UiLCJyZWFsVXNlcklkIjoiNWIzMTY2NGEtNzYxMC00ZmIwLWJjODktZjlmOGUrUC021FhB_zuETHmhQUXOfIyTkpvhcJfrrqwdcc-KmJGznckACLj65VmnayoltCce_3JGJ361GuutgrDaqp1aW4D05mvO8CCIRwGq8hTcRoi7IdXYSnA6UlXtLNYvttz92jaAAoNDCZmbbP-umHac4x5AT1xY-kVyh7VAadZG_Qe7dZWU9WCHtCV3mqTMwX9B9zrqY2NrpblevbbYpoiJiXOU7kex4BEivF1K6VWI-mpcmKtEOZLx2E",
			redacted: "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJkNDc1ZGQ1Yi1mMTZjLTRiZmItODk4Yy1kMzQzNWEyMTUyMzkiLCJlZmZlY3R********************vttz92jaAAoNDCZmbbP-umHac4x5AT1xY-kVyh7VAadZG_Qe7dZWU9WCHtCV3mqTMwX9B9zrqY2NrpblevbbYpoiJiXOU7kex4BEivF1K6VWI-mpcmKtEOZLx2E",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r := redactString(tC.str)
			assert.Equal(t, tC.redacted, r, tC.desc)
		})
	}
}
