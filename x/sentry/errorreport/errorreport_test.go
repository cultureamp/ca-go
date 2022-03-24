package errorreport_test

import (
	"testing"

	"github.com/cultureamp/ca-go/x/sentry/errorreport"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
)

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
