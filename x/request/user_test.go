package request_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/stretchr/testify/assert"
)

func newAuthenticatedUser() request.AuthenticatedUser {
	return request.AuthenticatedUser{
		CustomerAccountID: "123",
		UserID:            "456",
		RealUserID:        "789",
	}
}

func TestContextWithAuthenticatedUser(t *testing.T) {
	user := newAuthenticatedUser()
	ctx := context.Background()

	ctx = request.ContextWithAuthenticatedUser(ctx, user)
	userFromContext, ok := request.AuthenticatedUserFromContext(ctx)

	assert.Equal(t, user, userFromContext)
	assert.True(t, ok)
}

func ExampleContextWithAuthenticatedUser() {
	authenticatedUser := request.AuthenticatedUser{
		CustomerAccountID: "123",
		UserID:            "456",
		RealUserID:        "789",
	}
	ctx := context.Background()

	ctx = request.ContextWithAuthenticatedUser(ctx, authenticatedUser)

	if authenticatedUserFromContext, ok := request.AuthenticatedUserFromContext(ctx); ok {
		fmt.Println(authenticatedUserFromContext.CustomerAccountID)
		fmt.Println(authenticatedUserFromContext.UserID)
		fmt.Println(authenticatedUserFromContext.RealUserID)

		// Output:
		// 123
		// 456
		// 789
	}
}

func TestAuthenticatedUserFromContextMissing(t *testing.T) {
	ctx := context.Background()

	_, ok := request.AuthenticatedUserFromContext(ctx)
	assert.False(t, ok)
}

func ExampleContextHasAuthenticatedUser() {
	authenticatedUser := request.AuthenticatedUser{
		CustomerAccountID: "123",
		UserID:            "456",
		RealUserID:        "789",
	}
	ctx := context.Background()

	ctx = request.ContextWithAuthenticatedUser(ctx, authenticatedUser)

	ok := request.ContextHasAuthenticatedUser(ctx)
	fmt.Println(ok)

	// Output: true
}

func TestContextHasAuthenticatedUserSuceeds(t *testing.T) {
	ctx := request.ContextWithAuthenticatedUser(context.Background(), newAuthenticatedUser())

	ok := request.ContextHasAuthenticatedUser(ctx)
	assert.True(t, ok)
}

func TestContextHasAuthenticatedUserFails(t *testing.T) {
	ctx := context.Background()

	ok := request.ContextHasAuthenticatedUser(ctx)
	assert.False(t, ok)
}
