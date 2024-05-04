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
	"github.com/stretchr/testify/mock"
)

func TestMockedPackageLogger(t *testing.T) {
	// Revert the DefaultLogger once the test if finished
	stdLogger := log.DefaultLogger
	defer func() { log.DefaultLogger = stdLogger }()

	mockLogger := new(mockLogger)
	log.DefaultLogger = mockLogger

	nilProperty := &log.Property{}
	mockLogger.On("Debug", "should_call_mock").Return(nilProperty)

	log.Debug("should_call_mock").
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("detailed information explain")

	// Output:
	//
}

func TestCommonExamples(t *testing.T) {
	defer func() {
		// Force re-create of defaults when test finished
		log.DefaultLogger = nil
	}()

	t.Setenv("APP", "logger-test")
	t.Setenv("AWS_REGION", "dev")
	t.Setenv("PRODUCT", "cago")

	log.DefaultOptions(
		log.WithBool("global_bool", true),
		log.WithInt("global_int", 42),
		log.WithDuration("global_dur", 42*time.Second),
	)

	log.Debug("hander_added").
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("detailed information explain")

	props := log.Add().
		Str("custom", "field").
		Int("test-number", 2).
		Str("bar", "baz").
		Int("n", 1)

	log.Info("something_else_happened").
		Properties(props).
		Details("detailed information explain")

	log.Warn("something_did_not_work").
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 3),
		).Details("detailed information explain")

	err := errors.New("exception")
	log.Error("user_added", err).
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 4).
			Str("bar", "baz").
			Int("n", 1),
		).Details("detailed information explain")

	// log.Fatal calls os.exit() so this is hard to test!

	defer recoverFromPanic()
	log.Panic("panic_error", err).
		Properties(log.Add().
			Str("custom", "field").
			Int("test-number", 4).
			Str("bar", "baz").
			Int("n", 1),
		).Details("detailed information explain")

	// Output:
	// {"severity":"error","app":"test-app","app_version":"1.0.0","aws_region":"us-east-1","aws_account_id":"","farm":"test-farm","error":"exception","event":"event_details","test-str":"str","test-number":1,"time":"2024-01-03T15:57:02+11:00","details":"detailed message"}
}

func TestRequestExample(t *testing.T) {
	defer func() {
		// Force re-create of defaults when test finished
		log.DefaultLogger = nil
	}()

	t.Setenv("APP", "logger-test")
	t.Setenv("AWS_REGION", "dev")
	t.Setenv("PRODUCT", "cago")

	log.DefaultLogger = nil // force re-create

	// create a dummy request
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Add(log.TraceIDHeader, "trace_123_id")
	req.Header.Add(log.RequestIDHeader, "request_456_id")
	req.Header.Add(log.CorrelationIDHeader, "correlation_789_id")

	log.Info("info_event").
		WithRequestTracing(req).
		WithSystemTracing().
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain request headers")

	// Output:
	// {"severity":"info","app":"ca-go","app_version":"1.0.1","aws_region":"local","aws_account_id":"012345678901","farm":"local","product":"library","event":"info_event","trace_id":"trace_123_id","request_id":"request_456_id","correlation_id":"correlation_789_id","os":"darwin","test-str":"str","test-number":1,"time":"2024-01-03T19:46:45+11:00","details":"logging should contain request headers"}
}

func TestAuthPayloadExample(t *testing.T) {
	defer func() {
		// Force re-create of defaults when test finished
		log.DefaultLogger = nil
	}()

	t.Setenv("APP", "logger-test")
	t.Setenv("AWS_REGION", "dev")
	t.Setenv("PRODUCT", "cago")

	log.DefaultLogger = nil // force re-create

	// create a jwt payload
	auth := &log.AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	log.Info("info_event").
		WithAuthenticatedUserTracing(auth).
		WithSystemTracing().
		Properties(log.Add().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain auth payload")

	// Output:
	// {"severity":"info","app":"ca-go","app_version":"1.0.1","aws_region":"local","aws_account_id":"012345678901","farm":"local","product":"library","event":"info_event","account_id":"account_123_id","user_id":"user_789_id","real_user_id":"real_456_id","os":"darwin","test-str":"str","test-number":1,"time":"2024-01-03T19:49:33+11:00","details":"logging should contain request headers"}
}

func recoverFromPanic() {
	if saved := recover(); saved != nil {
		// convert to an error if it's not one already
		err, ok := saved.(error)
		if !ok {
			err = errors.New(fmt.Sprint(saved))
		}

		log.Error("recovered_from_panic", err).Send()
	}
}

type mockLogger struct {
	mock.Mock
}

func (ml *mockLogger) Debug(event string) *log.Property {
	args := ml.Called(event)
	return args.Get(0).(*log.Property)
}

func (ml *mockLogger) Info(event string) *log.Property {
	args := ml.Called(event)
	return args.Get(0).(*log.Property)
}

func (ml *mockLogger) Warn(event string) *log.Property {
	args := ml.Called(event)
	return args.Get(0).(*log.Property)
}

func (ml *mockLogger) Error(event string, err error) *log.Property {
	args := ml.Called(event)
	return args.Get(0).(*log.Property)
}

func (ml *mockLogger) Fatal(event string, err error) *log.Property {
	args := ml.Called(event)
	return args.Get(0).(*log.Property)
}

func (ml *mockLogger) Panic(event string, err error) *log.Property {
	args := ml.Called(event)
	return args.Get(0).(*log.Property)
}

func (ml *mockLogger) Child(options ...log.LoggerOption) log.Logger {
	args := ml.Called(options)
	return args.Get(0).(log.Logger)
}

func (ml *mockLogger) WithContext(ctx context.Context) context.Context {
	args := ml.Called(ctx)
	return args.Get(0).(context.Context)
}
