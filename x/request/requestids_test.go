package request_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/stretchr/testify/assert"
)

func newRequestIDs() request.RequestIDs {
	return request.RequestIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
}

func TestContextWithRequestIDs(t *testing.T) {
	ids := newRequestIDs()
	ctx := context.Background()

	ctx = request.ContextWithRequestIDs(ctx, ids)
	idsFromContext, ok := request.RequestIDsFromContext(ctx)

	assert.Equal(t, ids, idsFromContext)
	assert.True(t, ok)
}

func ExampleContextWithRequestIDs() {
	requestIDs := request.RequestIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
	ctx := context.Background()

	ctx = request.ContextWithRequestIDs(ctx, requestIDs)

	if requestIDsFromContext, ok := request.RequestIDsFromContext(ctx); ok {
		fmt.Println(requestIDsFromContext.RequestID)
		fmt.Println(requestIDsFromContext.CorrelationID)

		// Output:
		// 123
		// 456
	}
}

func TestRequestIDsFromContextMissing(t *testing.T) {
	ctx := context.Background()

	_, ok := request.RequestIDsFromContext(ctx)
	assert.False(t, ok)
}

func ExampleContextHasRequestIDs() {
	requestIDs := request.RequestIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
	ctx := context.Background()

	ctx = request.ContextWithRequestIDs(ctx, requestIDs)

	ok := request.ContextHasRequestIDs(ctx)
	fmt.Println(ok)
	// Output: true
}

func TestContextHasRequestIDsSucceeds(t *testing.T) {
	ctx := request.ContextWithRequestIDs(context.Background(), newRequestIDs())

	ok := request.ContextHasRequestIDs(ctx)
	assert.True(t, ok)
}

func TestContextHasRequestIDsFails(t *testing.T) {
	ctx := context.Background()

	ok := request.ContextHasRequestIDs(ctx)
	assert.False(t, ok)
}
