package errorreport_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/x/sentry/errorreport"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name                 string
		err                  error
		expectedSentryEvents int
	}{
		{
			name:                 "test not error",
			err:                  nil,
			expectedSentryEvents: 0,
		},
		{
			name:                 "test error",
			err:                  fmt.Errorf("test error"),
			expectedSentryEvents: 1,
		},
	}

	for _, test := range tests {
		t.Run("LambdaMiddleware: "+test.name, func(t *testing.T) {
			mockSentryTransport := setupMockSentryTransport(t)
			eventHandler := func(ctx context.Context, payload any) error {
				return test.err
			}
			wrapped := errorreport.LambdaMiddleware(eventHandler)

			err := wrapped(context.Background(), "random body")

			assert.Equal(t, test.err, err)
			assert.Len(t, mockSentryTransport.Events(), test.expectedSentryEvents)
		})

		t.Run("LambdaWithOutputMiddleware: "+test.name, func(t *testing.T) {
			mockSentryTransport := setupMockSentryTransport(t)
			eventHandler := func(ctx context.Context, payload any) (any, error) {
				return nil, test.err
			}
			wrapped := errorreport.LambdaWithOutputMiddleware(eventHandler)

			_, err := wrapped(context.Background(), "random body")

			assert.Equal(t, test.err, err)
			assert.Len(t, mockSentryTransport.Events(), test.expectedSentryEvents)
		})
	}
}

func TestPanic(t *testing.T) {
	tr := true
	fls := false

	unstableHandler := func(ctx context.Context, payload any) error {
		panic(fmt.Errorf("lol"))
	}

	unstableOutputHandler := func(ctx context.Context, payload any) (any, error) {
		panic(fmt.Errorf("lol"))
	}

	tests := []struct {
		name    string
		repanic *bool
	}{
		{
			name:    "sends event and rethrows panic by default",
			repanic: nil,
		},
		{
			name:    "sends event and rethrows panic when configured",
			repanic: &tr,
		},
		{
			name:    "sends event and swallows panic when configured",
			repanic: &fls,
		},
	}

	for _, test := range tests {
		options := []errorreport.LambdaOption{}

		if test.repanic != nil {
			options = append(options, errorreport.WithRepanic(*test.repanic))
		}

		t.Run("LambdaMiddleware: "+test.name, func(t *testing.T) {
			mockSentryTransport := setupMockSentryTransport(t)
			wrapped := errorreport.LambdaMiddleware(unstableHandler, options...)

			testFunc := func() {
				_ = wrapped(context.Background(), "random body")
			}

			if test.repanic == nil || *test.repanic {
				assert.PanicsWithError(t, "lol", testFunc)
			} else {
				assert.NotPanics(t, testFunc)
			}

			assert.Len(t, mockSentryTransport.Events(), 1)
		})

		t.Run("LambdaWithOutputMiddleware: "+test.name, func(t *testing.T) {
			mockSentryTransport := setupMockSentryTransport(t)
			wrapped := errorreport.LambdaWithOutputMiddleware(unstableOutputHandler, options...)

			testFunc := func() {
				_, _ = wrapped(context.Background(), "random body")
			}

			if test.repanic == nil || *test.repanic {
				assert.PanicsWithError(t, "lol", testFunc)
			} else {
				assert.NotPanics(t, testFunc)
			}

			assert.Len(t, mockSentryTransport.Events(), 1)
		})
	}
}
