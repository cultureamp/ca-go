package errorreport_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cultureamp/ca-go/x/sentry/errorreport"
	"github.com/getsentry/sentry-go"
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
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithDSN("https://public@sentry.example.com/1"),
			errorreport.WithRelease("my-app", "1.0.0"),
		)
		require.NoError(t, err)
	})

	t.Run("errors when environment is missing", func(t *testing.T) {
		err := errorreport.Init(
			errorreport.WithDSN("https://public@sentry.example.com/1"),
			errorreport.WithRelease("my-app", "1.0.0"),
		)
		require.EqualError(t, err, "mandatory fields missing: environment")
	})

	t.Run("errors when DSN is missing", func(t *testing.T) {
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithRelease("my-app", "1.0.0"),
		)
		require.EqualError(t, err, "mandatory fields missing: DSN")
	})

	t.Run("errors when release is missing", func(t *testing.T) {
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithDSN("https://public@sentry.example.com/1"),
		)
		require.EqualError(t, err, "mandatory fields missing: release")
	})

	t.Run("allows build details, transport, debug mode, and before filter to be supplied", func(t *testing.T) {
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
		err := errorreport.Init(
			errorreport.WithEnvironment("test"),
			errorreport.WithDSN("https://public@sentry.example.com/1"),
			errorreport.WithRelease("my-app", "1.0.0"),
			errorreport.WithServerlessTransport(),
		)
		require.NoError(t, err)
	})
}
