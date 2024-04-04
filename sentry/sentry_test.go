package sentry_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/x/sentry/errorreport"
	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleDecorate() {
	defer errorreport.Decorate(map[string]string{
		"key":    "123",
		"animal": "flamingo",
	})()

	// Since this API is designed around "defer", don't use it in a loop.
	// Instead, create a function and call that function in a loop.

	// Output:
}

func TestDecorate(t *testing.T) {
	ctx := context.Background()
	mockSentryTransport := setupMockSentryTransport(t)

	// record event with tag value in context
	popFn := errorreport.Decorate(map[string]string{
		"animal": "flamingo",
	})
	errorreport.ReportError(ctx, errors.New("with a flamingo"))

	// pop tagging context
	popFn()

	// report event without a tag value in context
	errorreport.ReportError(ctx, errors.New("i have no flamingo"))

	require.Len(t, mockSentryTransport.events, 2)

	eventWithTag := mockSentryTransport.events[0]
	assert.Equal(t, "flamingo", eventWithTag.Tags["animal"])

	eventNoTag := mockSentryTransport.events[1]
	assert.Equal(t, "", eventNoTag.Tags["animal"])
}

func TestConfigure(t *testing.T) {
	t.Run("no errors when all mandatory options supplied", func(t *testing.T) {
		testingScope(t)
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithDSN("https://public@sentry.example.com/1"),
			errorreport.WithRelease("my-app", "1.0.0"),
		)
		require.NoError(t, err)
	})

	t.Run("errors when environment is missing", func(t *testing.T) {
		testingScope(t)
		err := errorreport.Init(
			errorreport.WithDSN("https://public@sentry.example.com/1"),
			errorreport.WithRelease("my-app", "1.0.0"),
		)
		require.EqualError(t, err, "mandatory fields missing: environment")
	})

	t.Run("errors when DSN is missing", func(t *testing.T) {
		testingScope(t)
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithRelease("my-app", "1.0.0"),
		)
		require.EqualError(t, err, "mandatory fields missing: DSN")
	})

	t.Run("No error when DSN is missing, but environment is local", func(t *testing.T) {
		err := errorreport.Init(
			errorreport.WithEnvironment("local"),
			errorreport.WithRelease("my-app", "1.0.0"),
		)
		require.NoError(t, err)
	})

	t.Run("errors when release is missing", func(t *testing.T) {
		testingScope(t)
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithDSN("https://public@sentry.example.com/1"),
		)
		require.EqualError(t, err, "mandatory fields missing: release")
	})

	t.Run("allows build details, transport, debug mode, and before filter to be supplied", func(t *testing.T) {
		testingScope(t)
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithDSN("https://public@sentry.example.com/1"),
			errorreport.WithRelease("my-app", "1.0.0"),
			errorreport.WithBuildDetails("dolly", "100", "main", "ffff"),
			errorreport.WithTransport(&transportMock{}),
			errorreport.WithDebug(),
			errorreport.WithBeforeFilter(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				return event
			}),
		)
		require.NoError(t, err)
	})

	t.Run("allows a default serverless transport to be set", func(t *testing.T) {
		testingScope(t)
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithDSN("https://public@sentry.example.com/1"),
			errorreport.WithRelease("my-app", "1.0.0"),
			errorreport.WithServerlessTransport(),
		)
		require.NoError(t, err)
	})
}

func TestConfigureWithBeforeFilter(t *testing.T) {
	testingScope(t)
	ctx := context.Background()
	mockSentryTransport := setupMockSentryTransport(t,
		errorreport.WithBeforeFilter(errorreport.RootCauseAsTitle),
	)

	errorreport.ReportError(ctx, errors.New("a flamingo"))

	require.Len(t, mockSentryTransport.events, 1)

	event := mockSentryTransport.events[0]
	assert.Equal(t, "a flamingo", event.Exception[0].Type)
}

func TestConfigureWithTag(t *testing.T) {
	testingScope(t)
	ctx := context.Background()
	mockSentryTransport := setupMockSentryTransport(t,
		errorreport.WithTag("common_name", "james' flamingo"),
	)

	errorreport.ReportError(ctx, errors.New("with a flamingo"))

	require.Len(t, mockSentryTransport.events, 1)

	event := mockSentryTransport.events[0]
	assert.Equal(t, "james' flamingo", event.Tags["common_name"])
}

func TestConfigureWithTags(t *testing.T) {
	cases := []struct {
		name string
		tags map[string]string
	}{
		{
			name: "nil tags",
			tags: nil,
		},
		{
			name: "empty tags",
			tags: map[string]string{},
		},
		{
			name: "one tag",
			tags: map[string]string{
				"common_name": "james' flamingo",
			},
		},
		{
			name: "multiple tags",
			tags: map[string]string{
				"genus":   "phoenicoparrus",
				"species": "jamesi",
			},
		},
	}
	ctx := context.Background()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testingScope(t)
			mockSentryTransport := setupMockSentryTransport(t,

				errorreport.WithTags(tc.tags),
			)

			errorreport.ReportError(ctx, errors.New("with a flamingo"))

			require.Len(t, mockSentryTransport.events, 1)

			expected := tc.tags
			if expected == nil {
				expected = map[string]string{}
			}

			tags := mockSentryTransport.events[0].Tags

			assert.Equal(t, expected, tags)
		})
	}
}

func TestGracefullyShutdown(t *testing.T) {
	type myError struct{}

	cases := []struct {
		name      string
		panicFunc func(error)
		err       error
	}{
		{name: "shut down due to panic", panicFunc: func(err error) { panic(err) }, err: errors.New("panic error")},
		{name: "shut down due to log panic", panicFunc: func(err error) { log.Panic(err) }, err: errors.New("log panic error")},
		{name: "shut down due to unknown panic", panicFunc: func(err error) { panic(myError{}) }, err: errors.New("unknown panic: sentry_test.myError")},
		{name: "shut down due to panic with string message", panicFunc: func(err error) { panic("error with string message") }, err: errors.New("error with string message")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testingScope(t)
			mockSentryTransport := setupMockSentryTransport(t,
				errorreport.WithBeforeFilter(errorreport.RootCauseAsTitle),
			)

			defer func() {
				if err := recover(); err != nil {
					errorreport.GracefullyShutdown(err, time.Second*1)
				}
				require.Len(t, mockSentryTransport.events, 1)
				event := mockSentryTransport.events[0]
				assert.Equal(t, tc.err.Error(), event.Exception[0].Type)
			}()
			tc.panicFunc(tc.err)
		})
	}
}

func testingScope(t *testing.T) {
	t.Helper()
	sentry.PushScope()
	t.Cleanup(sentry.PopScope)
}
