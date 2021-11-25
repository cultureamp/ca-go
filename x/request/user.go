package request

import "context"

const authenticatedUserKey = contextValueKey("authenticatedUser")

// AuthenticatedUser holds the identifiers of a user making an authenticated
// request.
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

// ContextWithAuthenticatedUser returns a new context with the given user
// embedded as a value.
func ContextWithAuthenticatedUser(parent context.Context, user AuthenticatedUser) context.Context {
	ctx := context.WithValue(parent, authenticatedUserKey, user)
	return ctx
}

// AuthenticatedUserFromContext attempts to retrieve an AuthenticatedUser
// from the given context, returning an AuthenticatedUser along with a boolean
// signalling whether the retrieval was successful.
func AuthenticatedUserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	value := ctx.Value(authenticatedUserKey)

	user, ok := value.(AuthenticatedUser)
	return user, ok
}

// ContextHasAuthenticatedUser returns whether the given context contains
// an AuthenticatedUser value.
func ContextHasAuthenticatedUser(ctx context.Context) bool {
	_, ok := AuthenticatedUserFromContext(ctx)
	return ok
}
