package log

import (
	"context"
)

type contextAuthPayloadKey string

const authPayloadKey = contextAuthPayloadKey("auth_payload")

// AuthPayload contains the account, user and realuser ids usually obtains from an authenticated JWT.
type AuthPayload struct {
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

// ContextWithAuthPayload returns a new context with the given AuthUserIDs
// embedded as a value.
func ContextWithAuthPayload(ctx context.Context, ids AuthPayload) context.Context {
	return context.WithValue(ctx, authPayloadKey, ids)
}

// AuthPayloadFromContext attempts to retrieve a AuthUserIDs struct from the given
// context, returning a AuthUserIDs struct along with a boolean signalling
// whether the retrieval was successful.
func AuthPayloadFromContext(ctx context.Context) (AuthPayload, bool) {
	ids, ok := ctx.Value(authPayloadKey).(AuthPayload)
	return ids, ok
}
