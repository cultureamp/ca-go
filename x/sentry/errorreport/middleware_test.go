//nolint:bodyclose
package errorreport_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/cultureamp/ca-go/x/sentry/errorreport"
)

// setupContextForSentry returns a new context populated with sample
// request and user IDs, along with an assertion function. The function
// can be called by tests to check that the request and user IDs have
// successfully been sent to Sentry in the error report.
func setupContextForSentry() (context.Context, func(t *testing.T, transport *transportMock)) {
	ctx := context.Background()
	ctx = request.ContextWithAuthenticatedUser(ctx, request.AuthenticatedUser{
		CustomerAccountID: "123",
		UserID:            "456",
		RealUserID:        "789",
	})
	ctx = request.ContextWithRequestIDs(ctx, request.RequestIDs{
		RequestID:     "abc",
		CorrelationID: "def",
	})

	return ctx, func(t *testing.T, mockSentryTransport *transportMock) {
		t.Helper()

		assert.Len(t, mockSentryTransport.Events(), 1)
		sentryEvent := mockSentryTransport.Events()[0]

		assert.Equal(t, "123", sentryEvent.Tags["customer"])
		assert.Equal(t, "456", sentryEvent.User.ID)
		assert.Equal(t, "789", sentryEvent.Tags["user.real"])
		assert.Equal(t, "abc", sentryEvent.Tags["RequestID"])

		tracingContext, ok := sentryEvent.Contexts["Culture Amp - Tracing"]
		require.True(t, ok)
		assert.Equal(t, "def", tracingContext["CorrelationID"])
	}
}

func TestHTTPMiddleware(t *testing.T) {
	ctx, sentryContextAssertions := setupContextForSentry()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://www.example.com/happy_path",
		nil)
	require.NoError(t, err)

	t.Run("successful request", func(t *testing.T) {
		mockSentryTransport := setupMockSentryTransport(t)
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

		assert.Empty(t, mockSentryTransport.Events())
	})

	t.Run("unsuccessful request", func(t *testing.T) {
		mockSentryTransport := setupMockSentryTransport(t)
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

		assert.Equal(t, http.StatusTeapot, w.Result().StatusCode)

		// ...and reports the error to Sentry.
		sentryContextAssertions(t, mockSentryTransport)
	})

	t.Run("unsuccessful request with default panic handler", func(t *testing.T) {
		mockSentryTransport := setupMockSentryTransport(t)
		w := httptest.NewRecorder()

		mw := errorreport.NewHTTPMiddleware(nil)

		innerHandlerCalled := false
		innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			innerHandlerCalled = true
			panic("boom")
		})

		sut := mw(innerHandler)
		sut.ServeHTTP(w, req)

		// Executes the request handler...
		assert.True(t, innerHandlerCalled)

		// ...reports the error to Sentry...
		sentryContextAssertions(t, mockSentryTransport)

		// ...and produces a JSON:API style error response.
		assert.Equal(t, "{\"errors\":[{\"status\":\"500\",\"title\":\"Internal Server Error\"}]}", w.Body.String())
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})
}

func TestGoaEndpointMiddleware(t *testing.T) {
	ctx, sentryContextAssertion := setupContextForSentry()

	t.Run("successful request", func(t *testing.T) {
		mockSentryTransport := setupMockSentryTransport(t)

		endpointCalled := false
		endpoint := func(ctx context.Context, req interface{}) (interface{}, error) {
			endpointCalled = true

			return "foobar", nil
		}

		mw := errorreport.NewGoaEndpointMiddleware()

		sut := mw(endpoint)
		res, err := sut(ctx, nil)
		require.NoError(t, err)
		assert.Equal(t, "foobar", res)

		assert.True(t, endpointCalled)
		assert.Empty(t, mockSentryTransport.Events())
	})

	t.Run("unsuccessful request", func(t *testing.T) {
		mockSentryTransport := setupMockSentryTransport(t)

		endpointCalled := false
		endpoint := func(ctx context.Context, req interface{}) (interface{}, error) {
			endpointCalled = true

			return nil, errors.New("boom")
		}

		mw := errorreport.NewGoaEndpointMiddleware()

		sut := mw(endpoint)
		_, err := sut(ctx, nil)
		require.Error(t, err)

		assert.True(t, endpointCalled)
		sentryContextAssertion(t, mockSentryTransport)
	})
}
