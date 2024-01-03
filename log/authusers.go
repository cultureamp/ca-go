package log

import (
	"context"
)

type contextAuthUserIDKey string

const authUserIDsKey = contextAuthUserIDKey("auth_user_ids")

// AuthUserIDs contains the account, user and realuser ids usually obtains from an authenticated JWT.
type AuthUserIDs struct {
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

// ContextWithAuthUserIDs returns a new context with the given AuthUserIDs
// embedded as a value.
func ContextWithAuthUserIDs(ctx context.Context, ids AuthUserIDs) context.Context {
	return context.WithValue(ctx, authUserIDsKey, ids)
}

// AuthUserIDsFromContext attempts to retrieve a AuthUserIDs struct from the given
// context, returning a AuthUserIDs struct along with a boolean signalling
// whether the retrieval was successful.
func AuthUserIDsFromContext(ctx context.Context) (AuthUserIDs, bool) {
	ids, ok := ctx.Value(authUserIDsKey).(AuthUserIDs)
	return ids, ok
}
