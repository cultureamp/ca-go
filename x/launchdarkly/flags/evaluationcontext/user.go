package evaluationcontext

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/google/uuid"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

const (
	userAttributeAccountID  = "accountID"
	userAttributeRealUserID = "realUserID"
)

// User is a type of context, representing the identifiers and attributes of
// a human user to evaluate a flag against.
type User struct {
	key        string
	realUserID string
	accountID  string

	ldUser lduser.User
}

func (u User) ToLDUser() lduser.User {
	return u.ldUser
}

// UserOption are functions that can be supplied to configure a new user with
// additional attributes.
type UserOption func(*User)

// WithAccountID configures the user with the given account ID.
// This is the ID of the currently logged in user's parent account/organization,
// sometimes known as the "account_aggregate_id".
func WithAccountID(id string) UserOption {
	return func(u *User) {
		u.accountID = id
	}
}

// WithRealUserID configures the user with the given real user ID.
// This is the ID of the user who is currently impersonating the current user.
func WithRealUserID(id string) UserOption {
	return func(u *User) {
		u.realUserID = id
	}
}

// NewAnonymousUser returns a user object suitable for use in unauthenticated
// requests or requests with no access to user identifiers.
// Provide a unique session or request identifier as the key if possible. If the
// key is empty, it will default to 'ANONYMOUS_USER' and percentage rollouts
// will not be supported.
func NewAnonymousUser(key string) User {
	if key == "" {
		key = uuid.NewString()
	}

	u := User{
		key: key,
	}

	userBuilder := lduser.NewUserBuilder(u.key)
	userBuilder.Anonymous(true)
	u.ldUser = userBuilder.Build()

	return u
}

// NewUser returns a new user object with the given user ID and options.
// userID is the ID of the currently authenticated user, and will generally
// be a "user_aggregate_id".
func NewUser(userID string, opts ...UserOption) User {
	u := &User{
		key: userID,
	}

	for _, opt := range opts {
		opt(u)
	}

	userBuilder := lduser.NewUserBuilder(u.key)
	userBuilder.Custom(
		userAttributeAccountID,
		ldvalue.String(u.accountID))
	userBuilder.Custom(
		userAttributeRealUserID,
		ldvalue.String(u.realUserID))
	u.ldUser = userBuilder.Build()

	return *u
}

// UserFromContext extracts the effective user aggregate ID, real user aggregate
// ID, and account aggregate ID from the context. These values are used to
// create a new User object. An error is returned if user identifiers are not
// present in the context.
func UserFromContext(ctx context.Context) (User, error) {
	authenticatedUser, ok := request.AuthenticatedUserFromContext(ctx)
	if !ok {
		return User{}, errors.New("no AuthenticatedUser in supplied context")
	}

	return NewUser(
		authenticatedUser.UserID,
		WithAccountID(authenticatedUser.CustomerAccountID),
		WithRealUserID(authenticatedUser.RealUserID)), nil
}
