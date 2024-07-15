package request_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/request"
	"github.com/stretchr/testify/assert"
)

func newRequestIDs() request.HTTPFieldIDs {
	return request.HTTPFieldIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
}

func TestContextWithHTTPFieldIDs(t *testing.T) {
	ids := newRequestIDs()
	ctx := context.Background()

	ctx = request.ContextWithHTTPFieldIDs(ctx, ids)
	idsFromContext, ok := request.HTTPFieldIDsFromContext(ctx)

	assert.Equal(t, ids, idsFromContext)
	assert.True(t, ok)
}

func ExampleContextWithHTTPFieldIDs() {
	requestIDs := request.HTTPFieldIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
	ctx := context.Background()

	ctx = request.ContextWithHTTPFieldIDs(ctx, requestIDs)

	if requestIDsFromContext, ok := request.HTTPFieldIDsFromContext(ctx); ok {
		fmt.Println(requestIDsFromContext.RequestID)
		fmt.Println(requestIDsFromContext.CorrelationID)

		// Output:
		// 123
		// 456
	}
}

func TestHTTPFieldIDsFromContextMissing(t *testing.T) {
	ctx := context.Background()

	_, ok := request.HTTPFieldIDsFromContext(ctx)
	assert.False(t, ok)
}

func ExampleContextHasHTTPFieldIDs() {
	requestIDs := request.HTTPFieldIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
	ctx := context.Background()

	ctx = request.ContextWithHTTPFieldIDs(ctx, requestIDs)

	ok := request.ContextHasHTTPFieldIDs(ctx)
	fmt.Println(ok)
	// Output: true
}

func TestContextHasHTTPFieldIDsSucceeds(t *testing.T) {
	ctx := request.ContextWithHTTPFieldIDs(context.Background(), newRequestIDs())

	ok := request.ContextHasHTTPFieldIDs(ctx)
	assert.True(t, ok)
}

func TestContextHasHTTPFieldIDsFails(t *testing.T) {
	ctx := context.Background()

	ok := request.ContextHasHTTPFieldIDs(ctx)
	assert.False(t, ok)
}
