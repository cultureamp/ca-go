package log_test

import (
	"context"
	"time"

	"github.com/cultureamp/ca-go/log"
)

func ExampleLogInfoWithGlamplifyRequestFieldsFromCtx() {
	logger := getGlamplifyExtensionExampleLogger("INFO")

	// Log with no context
	logger.Info("info_with_glampify_request_field_tracing_no_ctx").
		WithGlamplifyRequestFieldsFromCtx(nil).
		Properties(log.SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain glamplify request fields tracing")

	rsFields := log.RequestScopedFields{
		TraceID:             "trace-123-id",
		RequestID:           "request-123-id",
		CorrelationID:       "correlation-123-id",
		UserAggregateID:     "user-123-id",
		CustomerAggregateID: "account-123-id",
	}

	ctx := context.Background()

	// Log with context but no request fields
	logger.Info("info_with_glampify_request_field_tracing_no_request_fields").
		WithGlamplifyRequestFieldsFromCtx(ctx).
		Properties(log.SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain glamplify request fields tracing")

	ctx = log.AddRequestFields(ctx, rsFields)

	// Log with context and request fields
	logger.Info("info_with_glampify_request_field_tracing").
		WithGlamplifyRequestFieldsFromCtx(ctx).
		Properties(log.SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain glamplify request fields tracing")

	// Output:
	// 2020-02-02T13:02:02+11:00 INF event="logging should contain glamplify request fields tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_glampify_request_field_tracing_no_ctx farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-02-02T13:02:02+11:00 INF event="logging should contain glamplify request fields tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_glampify_request_field_tracing_no_request_fields farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2020-02-02T13:02:02+11:00 INF event="logging should contain glamplify request fields tracing" app= app_version=1.0.0 authentication={"account_id":"account-123-id","user_id":"user-123-id"} aws_account_id=development aws_region= event=info_with_glampify_request_field_tracing farm=local product= properties={"resource":"resource_id","test-number":1} tracing={"correlation_id":"correlation-123-id","request_id":"request-123-id","trace_id":"trace-123-id"}
}

func getGlamplifyExtensionExampleLogger(sev string) log.Logger {
	config := log.NewLoggerConfig()
	config.LogLevel = sev
	config.Quiet = false
	config.ConsoleWriter = true
	config.ConsoleColour = false
	config.TimeNow = func() time.Time { return time.Unix(1580608922, 0) } // 1/1/2020
	return log.NewLogger(config)
}
