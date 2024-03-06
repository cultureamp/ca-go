package log

import (
	"context"
	"testing"
)

func TestExtensionWithGlamplifyRequestFieldsFromCtx(t *testing.T) {
	config := NewLoggerConfig()
	config.Quiet = false
	logger := NewLogger(config)

	// Log with no context
	logger.Info("info_with_glampify_request_field_tracing_no_ctx").
		WithGlamplifyRequestFieldsFromCtx(nil).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain glamplify request fields tracing")

	rsFields := RequestScopedFields{
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
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain glamplify request fields tracing")

	ctx = AddRequestFields(ctx, rsFields)

	// Log with context and request fields
	logger.Info("info_with_glampify_request_field_tracing").
		WithGlamplifyRequestFieldsFromCtx(ctx).
		Properties(SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain glamplify request fields tracing")

	// Local Console Output:
	// 2024-03-07T08:51:08+11:00 INF event="logging should contain glamplify request fields tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_glampify_request_field_tracing_no_ctx farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2024-03-07T08:51:08+11:00 INF event="logging should contain glamplify request fields tracing" app= app_version=1.0.0 aws_account_id=development aws_region= event=info_with_glampify_request_field_tracing_no_request_fields farm=local product= properties={"resource":"resource_id","test-number":1}
	// 2024-03-07T08:51:08+11:00 INF event="logging should contain glamplify request fields tracing" app= app_version=1.0.0 authentication={"account_id":"account-123-id","user_id":"user-123-id"} aws_account_id=development aws_region= event=info_with_glampify_request_field_tracing farm=local product= properties={"resource":"resource_id","test-number":1} tracing={"correlation_id":"correlation-123-id","request_id":"request-123-id","trace_id":"trace-123-id"}
}
