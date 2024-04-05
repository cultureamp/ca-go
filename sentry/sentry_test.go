package sentry_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/sentry"
	getsentry "github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleDecorate() {
	defer sentry.Decorate(map[string]string{
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
	popFn := sentry.Decorate(map[string]string{
		"animal": "flamingo",
	})
	sentry.ReportError(ctx, errors.New("with a flamingo"))

	// pop tagging context
	popFn()

	// report event without a tag value in context
	sentry.ReportError(ctx, errors.New("i have no flamingo"))

	require.Len(t, mockSentryTransport.events, 2)

	eventWithTag := mockSentryTransport.events[0]
	assert.Equal(t, "flamingo", eventWithTag.Tags["animal"])

	eventNoTag := mockSentryTransport.events[1]
	assert.Equal(t, "", eventNoTag.Tags["animal"])
}

func TestConfigure(t *testing.T) {
	t.Run("no errors when all mandatory options supplied", func(t *testing.T) {
		testingScope(t)
		err := sentry.Init(
			sentry.WithEnvironment("test"),
			sentry.WithDSN("https://public@sentry.example.com/1"),
			sentry.WithRelease("my-app", "1.0.0"),
		)
		require.NoError(t, err)
	})

	t.Run("errors when environment is missing", func(t *testing.T) {
		testingScope(t)
		err := sentry.Init(
			sentry.WithDSN("https://public@sentry.example.com/1"),
			sentry.WithRelease("my-app", "1.0.0"),
		)
		require.EqualError(t, err, "mandatory fields missing: environment")
	})

	t.Run("errors when DSN is missing", func(t *testing.T) {
		testingScope(t)
		err := sentry.Init(
			sentry.WithEnvironment("test"),
			sentry.WithRelease("my-app", "1.0.0"),
		)
		require.EqualError(t, err, "mandatory fields missing: DSN")
	})

	t.Run("No error when DSN is missing, but environment is local", func(t *testing.T) {
		err := sentry.Init(
			sentry.WithEnvironment("local"),
			sentry.WithRelease("my-app", "1.0.0"),
		)
		require.NoError(t, err)
	})

	t.Run("errors when release is missing", func(t *testing.T) {
		testingScope(t)
		err := sentry.Init(
			sentry.WithEnvironment("test"),
			sentry.WithDSN("https://public@sentry.example.com/1"),
		)
		require.EqualError(t, err, "mandatory fields missing: release")
	})

	t.Run("allows build details, transport, debug mode, and before filter to be supplied", func(t *testing.T) {
		testingScope(t)
		err := sentry.Init(
			sentry.WithEnvironment("test"),
			sentry.WithDSN("https://public@sentry.example.com/1"),
			sentry.WithRelease("my-app", "1.0.0"),
			sentry.WithBuildDetails("dolly", "100", "main", "ffff"),
			sentry.WithTransport(&transportMock{}),
			sentry.WithDebug(),
			sentry.WithBeforeFilter(func(event *getsentry.Event, hint *getsentry.EventHint) *getsentry.Event {
				return event
			}),
		)
		require.NoError(t, err)
	})

	t.Run("allows a default serverless transport to be set", func(t *testing.T) {
		testingScope(t)
		err := sentry.Init(
			sentry.WithEnvironment("test"),
			sentry.WithDSN("https://public@sentry.example.com/1"),
			sentry.WithRelease("my-app", "1.0.0"),
			sentry.WithServerlessTransport(),
		)
		require.NoError(t, err)
	})
}

func TestConfigureWithBeforeFilter(t *testing.T) {
	testingScope(t)
	ctx := context.Background()
	mockSentryTransport := setupMockSentryTransport(t,
		sentry.WithBeforeFilter(sentry.RootCauseAsTitle),
	)

	sentry.ReportError(ctx, errors.New("a flamingo"))

	require.Len(t, mockSentryTransport.events, 1)

	event := mockSentryTransport.events[0]
	assert.Equal(t, "a flamingo", event.Exception[0].Type)
}

func TestConfigureWithTag(t *testing.T) {
	testingScope(t)
	ctx := context.Background()
	mockSentryTransport := setupMockSentryTransport(t,
		sentry.WithTag("common_name", "james' flamingo"),
	)

	sentry.ReportError(ctx, errors.New("with a flamingo"))

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

				sentry.WithTags(tc.tags),
			)

			sentry.ReportError(ctx, errors.New("with a flamingo"))

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
				sentry.WithBeforeFilter(sentry.RootCauseAsTitle),
			)

			defer func() {
				if err := recover(); err != nil {
					sentry.GracefullyShutdown(err, time.Second*1)
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
	getsentry.PushScope()
	t.Cleanup(getsentry.PopScope)
}
