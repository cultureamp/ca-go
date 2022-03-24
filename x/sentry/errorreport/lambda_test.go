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
		mockSentryTransport := setupMockSentryTransport(t)

		t.Run(test.name, func(t *testing.T) {
			wrapped := errorreport.LambdaMiddleware(test.testEventHandler)

			_ = wrapped(context.Background(), "random body")

			assert.Len(t, mockSentryTransport.Events(), test.expectedSentryEvents)
		})
	}
}

func TestPanic(t *testing.T) {
	tr := true
	fls := false

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

	unstableHandler := func(ctx context.Context, payload interface{}) error {
		panic(fmt.Errorf("lol"))
	}

	for _, test := range tests {
		mockSentryTransport := setupMockSentryTransport(t)

		options := []errorreport.LambdaOption{}

		if test.repanic != nil {
			options = append(options, errorreport.WithRepanic(*test.repanic))
		}

		t.Run(test.name, func(t *testing.T) {
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
	}
}
