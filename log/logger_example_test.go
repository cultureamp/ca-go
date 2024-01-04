package log_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/log"
)

func TestBasicExamples(t *testing.T) {
	ctx := context.Background()

	log.Debug(ctx, "hander_added").
		Str("resource", "resource_id").
		Int("test-number", 1).
		Msg("detailed information explain")

	log.Info(ctx, "something_else_happened").
		Str("resource", "resource_id").
		Int("test-number", 2).
		Msg("detailed information explain")

	log.Warn(ctx, "something_did_not_work").
		Str("resource", "resource_id").
		Int("test-number", 3).
		Msg("detailed information explain")

	err := errors.New("exception")
	log.Error(ctx, "user_added", err).
		Str("resource", "resource_id").
		Int("test-number", 4).
		Dict("properties", log.Properties().
			Str("bar", "baz").
			Int("n", 1),
		).Msg("detailed information explain")

	// log.Fatal calls os.exit() so this is hard to test!

	defer recoverFromPanic(ctx)
	log.Panic(ctx, "panic_error", err).
		Int("test-number", 4).
		Dict("properties", log.Properties().
			Str("bar", "baz").
			Int("n", 1),
		).Msg("detailed information explain")

	// Output:
	// {"severity":"error","app":"test-app","app_version":"1.0.0","aws_region":"us-east-1","aws_account_id":"","farm":"test-farm","error":"exception","event":"event_details","test-str":"str","test-number":1,"time":"2024-01-03T15:57:02+11:00","details":"detailed message"}
}

func TestGlamplifyLogFieldExamples(t *testing.T) {
	ctx := context.Background()

	type fields map[string]interface{}

	now := time.Now()
	f := &fields{
		"key1":  "string value",
		"key2":  1,
		"now":   now.Format(time.RFC3339),
		"later": time.Now(),
	}

	log.Info(ctx, "log_fields").
		Interface("properties", f).
		Msg("detailed information explain")
}

func TestRequestExample(t *testing.T) {
	ctx := context.Background()

	// create a dummy request
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Add(log.TraceHeader, "trace_123_id")
	req.Header.Add(log.RequestHeader, "request_456_id")
	req.Header.Add(log.CorrelationHeader, "correlation_789_id")

	ctx = log.ContextWithRequest(ctx, req)

	log.Info(ctx, "info_event").
		Str("resource", "resource_id").
		Int("test-number", 1).
		Msg("logging should contain request headers")

	// Output:
	// {"severity":"info","app":"ca-go","app_version":"1.0.1","aws_region":"local","aws_account_id":"012345678901","farm":"local","product":"library","event":"info_event","trace_id":"trace_123_id","request_id":"request_456_id","correlation_id":"correlation_789_id","os":"darwin","test-str":"str","test-number":1,"time":"2024-01-03T19:46:45+11:00","details":"logging should contain request headers"}
}

func TestAuthPayloadExample(t *testing.T) {
	ctx := context.Background()

	// create a jwt payload
	auth := log.AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	ctx = log.ContextWithAuthPayload(ctx, auth)

	log.Info(ctx, "info_event").
		Str("resource", "resource_id").
		Int("test-number", 1).
		Msg("logging should contain auth payload")

	// Output:
	// {"severity":"info","app":"ca-go","app_version":"1.0.1","aws_region":"local","aws_account_id":"012345678901","farm":"local","product":"library","event":"info_event","account_id":"account_123_id","user_id":"user_789_id","real_user_id":"real_456_id","os":"darwin","test-str":"str","test-number":1,"time":"2024-01-03T19:49:33+11:00","details":"logging should contain request headers"}
}

func recoverFromPanic(ctx context.Context) {
	if saved := recover(); saved != nil {
		// convert to an error if it's not one already
		err, ok := saved.(error)
		if !ok {
			err = errors.New(fmt.Sprint(saved))
		}

		log.Error(ctx, "recovered_from_panic", err).Send()
	}
}
