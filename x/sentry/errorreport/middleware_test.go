package errorreport_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cultureamp/ca-go/x/sentry/errorreport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSentry(t *testing.T) *transportMock {
	t.Helper()

	mockSentryTransport := &transportMock{}
	err := errorreport.Configure(
		errorreport.WithEnvironment("test"),
		errorreport.WithDSN("https://public@sentry.example.com/1"),
		errorreport.WithRelease("1.0.0"),
		errorreport.WithTransport(mockSentryTransport),
	)
	require.NoError(t, err)
	require.NoError(t, errorreport.Connect())

	return mockSentryTransport
}

func TestHTTPMiddleware(t *testing.T) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"http://www.example.com/happy_path",
		nil)
	require.NoError(t, err)

	t.Run("successful request", func(t *testing.T) {
		mockSentryTransport := setupSentry(t)
		w := httptest.NewRecorder()

		panicHandlerCalled := false
		mw := errorreport.NewHTTPMiddleware(func(c context.Context, w http.ResponseWriter, err error) {
			panicHandlerCalled = true
		})

		innerHandlerCalled := false
		innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			innerHandlerCalled = true
		})

		sut := mw(innerHandler)
		sut.ServeHTTP(w, req)

		assert.False(t, panicHandlerCalled)
		assert.True(t, innerHandlerCalled)

		assert.Len(t, mockSentryTransport.Events(), 0)
	})

	t.Run("unsuccessful request", func(t *testing.T) {
		mockSentryTransport := setupSentry(t)
		w := httptest.NewRecorder()

		panicHandlerCalled := false
		mw := errorreport.NewHTTPMiddleware(func(c context.Context, w http.ResponseWriter, err error) {
			panicHandlerCalled = true
		})

		innerHandlerCalled := false
		innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			innerHandlerCalled = true
			w.WriteHeader(http.StatusTeapot)
			panic("boom")
		})

		sut := mw(innerHandler)
		sut.ServeHTTP(w, req)

		// Executes both handlers...
		assert.True(t, panicHandlerCalled)
		assert.True(t, innerHandlerCalled)

		// ...recovers the panic...
		// nolint:bodyclose
		assert.Equal(t, http.StatusTeapot, w.Result().StatusCode)

		// ...and reports the error to Sentry.
		assert.Len(t, mockSentryTransport.Events(), 1)
	})
}
