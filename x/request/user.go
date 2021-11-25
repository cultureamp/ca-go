package request

import "context"

const authenticatedUserKey = contextValueKey("authenticatedUser")

type AuthenticatedUser struct {
	// CustomerAccountID is the ID of the currently logged in user's parent
	// account/organization, sometimes known as the "account_aggregate_id".
	CustomerAccountID string
	// UserID is the ID of the currently authenticated user, and will
	// generally be a "user_aggregate_id".
	UserID string
	// RealUserID, when supplied, is the ID of the user who is currently
	// impersonating the current "UserID". This value is optional.
	RealUserID string
}

func ContextWithAuthenticatedUser(parent context.Context, user AuthenticatedUser) context.Context {
	ctx := context.WithValue(parent, authenticatedUserKey, user)
	return ctx
}

func AuthenticatedUserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	value := ctx.Value(authenticatedUserKey)

	user, ok := value.(AuthenticatedUser)
	return user, ok
}

func ContextHasAuthenticatedUser(ctx context.Context) bool {
	_, ok := AuthenticatedUserFromContext(ctx)
	return ok
}
