package log

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func ExampleLogger_Info_withGlobalExtensions() {
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
		WithRequestDiagnostics(req),
		WithRequestTracing(req),
		WithAuthorizationTracing(req),
		WithAuthenticatedUserTracing(auth),
		WithDatadogTracing(spanCtx),
		// WithSystemTracing(), // this logs "pid" which keeps changing
	)

	logger.Info("info_with_global_extensions").
		Properties(Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain global extensions fields")

	logger.Info("info_with_global_extensions").
		Properties(Add().
			Str("page", "home").
			Int("size", 42),
		).Details("logging should contain global extensions fields")

	// Output:
	// 2020-11-14T11:30:32Z INF event="logging should contain global extensions fields" app=logger-test app_version=1.0.0 authentication={"account_id":"account_123_id","realuser_id":"real_456_id","user_id":"user_789_id"} authorization={"authorization_token":"AWS**********ken","user_agent":"node","x_forwarded_for":"123.123.123","xca_service_authorization_token":"Bear**********oken"} aws_account_id=development aws_region=def dd.span_id=0 dd.trace_id=0 event=info_with_global_extensions farm=local product=cago properties={"resource":"resource_id","test-number":1} request={"host":"example.com","method":"GET","proto":"HTTP/1.1","scheme":"http"} tracing={"correlation_id":"correlation_789_id","request_id":"request_456_id","trace_id":"trace_123_id"}
	// 2020-11-14T11:30:32Z INF event="logging should contain global extensions fields" app=logger-test app_version=1.0.0 authentication={"account_id":"account_123_id","realuser_id":"real_456_id","user_id":"user_789_id"} authorization={"authorization_token":"AWS**********ken","user_agent":"node","x_forwarded_for":"123.123.123","xca_service_authorization_token":"Bear**********oken"} aws_account_id=development aws_region=def dd.span_id=0 dd.trace_id=0 event=info_with_global_extensions farm=local product=cago properties={"page":"home","size":42} request={"host":"example.com","method":"GET","proto":"HTTP/1.1","scheme":"http"} tracing={"correlation_id":"correlation_789_id","request_id":"request_456_id","trace_id":"trace_123_id"}
}

func ExampleLogger_Info_withNilGlobalExtensions() {
	config := getExampleLoggerConfig("INFO")
	logger := NewLogger(config,
		WithRequestDiagnostics(nil),
		WithRequestTracing(nil),
		WithAuthorizationTracing(nil),
		WithAuthenticatedUserTracing(nil),
		WithDatadogTracing(nil),
		// WithSystemTracing(), // this logs "pid" which keeps changing
	)

	logger.Info("info_with_global_extensions").
		Properties(Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain global extensions fields")

	logger.Info("info_with_global_extensions").
		Properties(Add().
			Str("page", "home").
			Int("size", 42),
		).Details("logging should contain global extensions fields")

	// Output:
	// 2020-11-14T11:30:32Z INF event="logging should contain global extensions fields" app=logger-test app_version=1.0.0 aws_account_id=development aws_region=def event=info_with_global_extensions farm=local product=cago properties={"resource":"resource_id","test-number":1}
	// 2020-11-14T11:30:32Z INF event="logging should contain global extensions fields" app=logger-test app_version=1.0.0 aws_account_id=development aws_region=def event=info_with_global_extensions farm=local product=cago properties={"page":"home","size":42}
}

func ExampleLogger_Warn_withAllGlobalProperties() {
	then := time.Date(2023, 11, 14, 11, 30, 32, 0, time.UTC)
	u := uuid.MustParse("e5fa7acf-1846-41b4-a2ee-80ecd86fb060")
	duration := time.Second * 42
	f := func(e *zerolog.Event) { e.Str("func", "val") }
	b := []byte("some bytes")

	var ui uint
	var i64 int64
	var ui64 uint64
	var f32 float32
	var f64 float64

	ui = 234
	i64 = 123
	ui64 = 123
	f32 = 32.32
	f64 = 64.64

	props := Add().
		Str("str", "value").
		Int("int", 1).
		UInt("uint", ui).
		Int64("int64", i64).
		UInt64("uint64", ui64).
		Float32("float32", f32).
		Float64("float64", f64).
		Bool("bool", true).
		Bytes("bytes", b).
		Duration("dur", duration).
		Time("time", then).
		IPAddr("ipaddr", net.IPv4bcast).
		UUID("uuid", u).
		Func(f)

	config := getExampleLoggerConfig("DEBUG")
	logger := NewLogger(config,
		WithProperties(props),
	)

	logger.Warn("debug_with_all_global_field_types").
		Detailsf("logging should contain all types: %s", "ok")

	// Output:
	// 2020-11-14T11:30:32Z WRN event="logging should contain all types: ok" app=logger-test app_version=1.0.0 aws_account_id=development aws_region=def default_properties={"bool":true,"bytes":"some bytes","dur":"PT42S","float32":32.32,"float64":64.64,"func":"val","int":1,"int64":123,"ipaddr":"255.255.255.255","str":"value","time":"2023-11-14T11:30:32Z","uint":234,"uint64":123,"uuid":"e5fa7acf-1846-41b4-a2ee-80ecd86fb060"} event=debug_with_all_global_field_types farm=local product=cago
}

func ExampleLogger_Warn_withDupicatePropertiesDocs() {
	global_props := Add().
		Str("global_str", "value").
		Int("global_int", 1)

	local_props := Add().
		Str("local_str", "value").
		Int("local_int", 1)

	config := getExampleLoggerConfig("DEBUG")
	logger := NewLogger(config, WithProperties(global_props))

	logger.Warn("debug_with_all_global_field_types").
		Properties(local_props).
		Detailsf("logging should contain all types: %s", "ok")

	// Output:
	// 2020-11-14T11:30:32Z WRN event="logging should contain all types: ok" app=logger-test app_version=1.0.0 aws_account_id=development aws_region=def default_properties={"global_int":1,"global_str":"value"} event=debug_with_all_global_field_types farm=local product=cago properties={"local_int":1,"local_str":"value"}
}
