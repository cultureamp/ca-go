package request_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/request"
	"github.com/stretchr/testify/assert"
)

func newUniqueIDs() request.UniqueIDs {
	return request.UniqueIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
}

func TestContextWithUniqueIDs(t *testing.T) {
	ids := newUniqueIDs()
	ctx := context.Background()

	ctx = request.ContextWithUniqueIDs(ctx, ids)
	idsFromContext, ok := request.UniqueIDsFromContext(ctx)

	assert.Equal(t, ids, idsFromContext)
	assert.True(t, ok)
}

func ExampleContextWithUniqueIDs() {
	requestIDs := request.UniqueIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
	ctx := context.Background()

	ctx = request.ContextWithUniqueIDs(ctx, requestIDs)

	if requestIDsFromContext, ok := request.UniqueIDsFromContext(ctx); ok {
		fmt.Println(requestIDsFromContext.RequestID)
		fmt.Println(requestIDsFromContext.CorrelationID)

		// Output:
		// 123
		// 456
	}
}

func TestUniqueIDsFromContextMissing(t *testing.T) {
	ctx := context.Background()

	_, ok := request.UniqueIDsFromContext(ctx)
	assert.False(t, ok)
}

func ExampleContextHasUniqueIDs() {
	requestIDs := request.UniqueIDs{
		RequestID:     "123",
		CorrelationID: "456",
	}
	ctx := context.Background()

	ctx = request.ContextWithUniqueIDs(ctx, requestIDs)

	ok := request.ContextHasUniqueIDs(ctx)
	fmt.Println(ok)
	// Output: true
}

func TestContextHasUniqueIDsSucceeds(t *testing.T) {
	ctx := request.ContextWithUniqueIDs(context.Background(), newUniqueIDs())

	ok := request.ContextHasUniqueIDs(ctx)
	assert.True(t, ok)
}

func TestContextHasUniqueIDsFails(t *testing.T) {
	ctx := context.Background()

	ok := request.ContextHasUniqueIDs(ctx)
	assert.False(t, ok)
}
