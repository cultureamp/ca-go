package sentry_test

import (
	"sync"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/sentry"
	getsentry "github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
)

func setupMockSentryTransport(t *testing.T, opts ...sentry.Option) *transportMock {
	t.Helper()

	mockSentryTransport := &transportMock{}

	defaultOpts := []sentry.Option{
		sentry.WithEnvironment("test"),
		sentry.WithDSN("https://public@sentry.example.com/1"),
		sentry.WithRelease("my-app", "1.0.0"),
	}

	allOpts := make([]sentry.Option, 0, len(defaultOpts)+len(opts)+1)

	// merge default options with user-supplied options, ensuring that the transport is the last option
	allOpts = append(allOpts, defaultOpts...)
	allOpts = append(allOpts, opts...)
	allOpts = append(allOpts, sentry.WithTransport(mockSentryTransport))

	err := sentry.Init(
		allOpts...,
	)
	require.NoError(t, err)

	return mockSentryTransport
}

// From https://github.com/getsentry/sentry-go/blob/bd116d6ce79b604297c6497aa07d7ac01768adbb/mocks_test.go#L24-L44
type transportMock struct {
	mu        sync.Mutex
	events    []*getsentry.Event
	lastEvent *getsentry.Event
}

func (t *transportMock) Configure(options getsentry.ClientOptions) {}
func (t *transportMock) SendEvent(event *getsentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
	t.lastEvent = event
}

func (t *transportMock) Flush(timeout time.Duration) bool {
	return true
}

func (t *transportMock) Events() []*getsentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.events
}
