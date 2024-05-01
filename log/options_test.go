package log

import (
	"context"
	"net/http"
	"net/http/httptest"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func ExampleLogger_Info_withGlobalOptions() {
	ctx := context.Background()

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Add(TraceIDHeader, "trace_123_id")
	req.Header.Add(RequestIDHeader, "request_456_id")
	req.Header.Add(CorrelationIDHeader, "correlation_789_id")
	req.Header.Add(AuthorizationHeader, "AWS 123 token")
	req.Header.Add(XCAServiceGatewayAuthorizationHeader, "Bearer 456 token")
	req.Header.Add(XForwardedForHeader, "123.123.123")
	req.Header.Add(UserAgentHeader, "node")

	auth := &AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	span, spanCtx := tracer.StartSpanFromContext(ctx, "test_span")
	defer span.Finish()

	config := getExampleLoggerConfig("INFO")
	logger := NewLogger(config,
		WithBool("global_bool", true),
		WithRequestDiagnostics(req),
		WithAuthorizationTracing(req),
		WithAuthenticatedUserTracing(auth),
		WithDatadogTracing(spanCtx),
		// WithSystemTracing(), // this logs "pid" which keeps changing
	)

	logger.Info("info_with_global_settings").
		Properties(Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain global fields")

	logger.Info("info_with_global_settings").
		Properties(Add().
			Str("page", "home").
			Int("size", 42),
		).Details("logging should contain global fields")

	// Output:
	// 2020-11-14T11:30:32Z INF event="logging should contain global fields" app=logger-test app_version=1.0.0 authentication={"account_id":"account_123_id","realuser_id":"real_456_id","user_id":"user_789_id"} authorization={"authorization_token":"AWS**********ken","user_agent":"node","x_forwarded_for":"123.123.123","xca_service_authorization_token":"Bear**********oken"} aws_account_id=development aws_region=def dd.span_id=0 dd.trace_id=0 event=info_with_global_settings farm=local global_bool=true product=cago properties={"resource":"resource_id","test-number":1} request={"host":"example.com","method":"GET","proto":"HTTP/1.1","scheme":"http"}
	// 2020-11-14T11:30:32Z INF event="logging should contain global fields" app=logger-test app_version=1.0.0 authentication={"account_id":"account_123_id","realuser_id":"real_456_id","user_id":"user_789_id"} authorization={"authorization_token":"AWS**********ken","user_agent":"node","x_forwarded_for":"123.123.123","xca_service_authorization_token":"Bear**********oken"} aws_account_id=development aws_region=def dd.span_id=0 dd.trace_id=0 event=info_with_global_settings farm=local global_bool=true product=cago properties={"page":"home","size":42} request={"host":"example.com","method":"GET","proto":"HTTP/1.1","scheme":"http"}
}
