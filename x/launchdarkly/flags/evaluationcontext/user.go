package evaluationcontext

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/google/uuid"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

const (
	userAttributeUserID     = "userID"
	userAttributeAccountID  = "accountID"
	userAttributeRealUserID = "realUserID"
	userAttributeSubdomain  = "subdomain"
)

// User is a type of context, representing the identifiers and attributes of
// a human user to evaluate a flag against.
type User struct {
	key        string
	userID     string
	realUserID string
	accountID  string

	ldContext ldcontext.Context
}

// ToLDContext transforms the context implementation into an LDContext object that can
// be understood by LaunchDarkly when evaluating a flag.
func (u User) ToLDContext() ldcontext.Context {
	return u.ldContext
}

// ToLDUser transforms the context implementation into an LDUser object that can
// be understood by LaunchDarkly when evaluating a flag.
//
// Deprecated: use ToLDContext() instead
func (u User) ToLDUser() ldcontext.Context {
	return u.ToLDContext()
}

// UserOption are functions that can be supplied to configure a new user with
// additional attributes.
type UserOption func(*User)

// WithUserAccountID configures the user with the given account ID.
// This is the ID of the currently logged in user's parent account/organization,
// sometimes known as the "account_aggregate_id".
func WithUserAccountID(id string) UserOption {
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
// key is empty, it will default to an uuid so percentage rollouts will still apply.
// No userID will be given to an anonymous user.
func NewAnonymousUser(key string) User {
	if key == "" {
		key = uuid.NewString()
	}

	return User{
		key: key,
		ldContext: ldcontext.NewBuilder(key).
			Anonymous(true).
			Build(),
	}
}

// NewAnonymousUserWithSubdomain returns a user object suitable for use in unauthenticated with known subdomain
// requests or requests with no access to user identifiers.
// Provide a unique session or request identifier as the key if possible. If the
// key is empty, it will default to an uuid so percentage rollouts will still apply.
// No userID will be given to an anonymous user.
func NewAnonymousUserWithSubdomain(key string, subdomain string) User {
	if key == "" {
		key = uuid.NewString()
	}

	return User{
		key: key,
		ldContext: ldcontext.NewMultiBuilder().
			// use existing user shape for current implementation
			Add(ldcontext.NewBuilder(key).SetString(userAttributeSubdomain, subdomain).Anonymous(true).Build()).
			Add(ldcontext.NewBuilder(key).Anonymous(true).Kind("account").Name(subdomain).Build()).
			Build(),
	}
}

// NewUser returns a new user object with the given user ID and options.
// userID is the ID of the currently authenticated user, and will generally
// be a "user_aggregate_id".
func NewUser(userID string, opts ...UserOption) User {
	u := &User{
		key:    userID,
		userID: userID,
	}

	for _, opt := range opts {
		opt(u)
	}

	userBuilder := ldcontext.NewBuilder(u.key)
	userBuilder.SetString(
		userAttributeRealUserID,
		u.realUserID)
	userBuilder.SetString(
		userAttributeUserID,
		u.userID)
	userBuilder.SetString(
		userAttributeAccountID,
		u.accountID)
	userContext := userBuilder.Build()
	u.ldContext = userContext

	if u.accountID != "" {
		accountContext := ldcontext.NewBuilder(u.accountID).Kind("account").Build()
		u.ldContext = ldcontext.NewMulti(userContext, accountContext)
	}

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
		WithUserAccountID(authenticatedUser.CustomerAccountID),
		WithRealUserID(authenticatedUser.RealUserID)), nil
}
