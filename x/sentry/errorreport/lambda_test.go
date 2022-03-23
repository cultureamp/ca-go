package errorreport

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// From https://github.com/getsentry/sentry-go/blob/bd116d6ce79b604297c6497aa07d7ac01768adbb/mocks_test.go#L24-L44
type transportMock struct {
	mu        sync.Mutex
	events    []*sentry.Event
	lastEvent *sentry.Event
}

func (t *transportMock) Configure(_ sentry.ClientOptions) {}

func (t *transportMock) SendEvent(event *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
	t.lastEvent = event
}

func (t *transportMock) Flush(_ time.Duration) bool {
	return true
}

func (t *transportMock) Events() []*sentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.events
}

func setupSentry(t *testing.T) *transportMock {
	t.Helper()

	mockSentryTransport := &transportMock{}
	err := Init(
		WithEnvironment("test"),
		WithDSN("https://public@sentry.example.com/1"),
		WithRelease("my-app", "1.0.0"),
		WithTransport(mockSentryTransport),
	)
	require.NoError(t, err)

	return mockSentryTransport
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name                 string
		testEventHandler     func(ctx context.Context, payload interface{}) error
		expectedSentryEvents int
	}{
		{
			name: "test not error",
			testEventHandler: func(ctx context.Context, payload interface{}) error {
				return nil
			},
			expectedSentryEvents: 0,
		},
		{
			name: "test error",
			testEventHandler: func(ctx context.Context, payload interface{}) error {
				return fmt.Errorf("test error")
			},
			expectedSentryEvents: 1,
		},
	}

	for _, test := range tests {
		mockSentryTransport := setupSentry(t)
		handler := New(PanicOptions{})

		t.Run(test.name, func(t *testing.T) {
			wrapped := handler.Handle(test.testEventHandler)
			_ = wrapped(context.Background(), "random body")
			if !assert.Len(t, mockSentryTransport.Events(), test.expectedSentryEvents) {
				t.Errorf("got %v, want %v", len(mockSentryTransport.Events()), test.expectedSentryEvents)
			}
		})
	}
}

func TestPanic(t *testing.T) {
	handler := New(PanicOptions{
		Repanic: true,
	})
	wrapped := handler.Handle(func(ctx context.Context, payload interface{}) error {
		panic(fmt.Errorf("lol"))
	})
	if !assert.PanicsWithError(t, "lol", func() {
		_ = wrapped(context.Background(), "random body")
	}) {
		t.Errorf("should panic with errString %s but didn't", "lol")
	}
}
