package evaluationcontext

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/google/uuid"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

// User is a type of context, representing the identifiers and attributes of
// a human user to evaluate a flag against.
//
// Deprecated: User is now replaced by EvaluationContext
type User struct {
	EvaluationContext
}

// ToLDUser transforms the context implementation into an LDUser object that can
// be understood by LaunchDarkly when evaluating a flag.
//
// Deprecated: User is deprecated use EvaluationContext with ToLDContext instead
func (u User) ToLDUser() ldcontext.Context {
	return u.ToLDContext()
}

// UserOption are functions that can be supplied to configure a new user with
// additional attributes.
//
// Deprecated: User is now replaced by EvaluationContext, use ContextOption with NewEvaluationContext instead
type UserOption func(*User)

// WithUserAccountID configures the user with the given account ID.
// This is the ID of the currently logged in user's parent account/organization,
// sometimes known as the "account_aggregate_id".
//
// Deprecated: User is now replaced by EvaluationContext, use WithAccountID instead
func WithUserAccountID(id string) UserOption {
	return func(u *User) {
		u.accountID = id
	}
}

// WithRealUserID configures the user with the given real user ID.
// This is the ID of the user who is currently impersonating the current user.
//
// Deprecated: UserOption is replaced by ContextOption on EvaluationContext, use WithContextRealUserID
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
//
// Deprecated: User is now replaced by EvaluationContext, use NewEvaluationContext without any opts instead
func NewAnonymousUser(key string) User {
	if key == "" {
		key = uuid.NewString()
	}

	return User{
		EvaluationContext{
			ldContext: ldcontext.NewBuilder(key).
				Anonymous(true).
				Build(),
		},
	}
}

// NewAnonymousUserWithSubdomain returns a user object suitable for use in unauthenticated with known subdomain
// requests or requests with no access to user identifiers.
// Provide a unique session or request identifier as the key if possible. If the
// key is empty, it will default to an uuid so percentage rollouts will still apply.
// No userID will be given to an anonymous user.
//
// Deprecated: User is now replaced by EvaluationContext, use NewAnonymousContextWithSubdomain instead
func NewAnonymousUserWithSubdomain(key string, subdomain string) User {
	if key == "" {
		key = uuid.NewString()
	}
	evaluationContext := NewAnonymousContextWithSubdomain(key, subdomain)
	evaluationContext.ldContext = ldcontext.NewMultiBuilder().
		Add(evaluationContext.ldContext).
		Add(ldcontext.NewBuilder(key).
			Kind(contextKindUser).
			SetString(contextAttributeSubdomain, subdomain).
			Anonymous(true).
			Build()).
		Build()
	return User{evaluationContext}
}

// NewUser returns a new user object with the given user ID and options.
// userID is the ID of the currently authenticated user, and will generally
// be a "user_aggregate_id".
//
// Deprecated: User is now replaced by EvaluationContext, use NewEvaluationContext instead
func NewUser(userID string, opts ...UserOption) User {
	u := &User{}
	u.userID = userID
	for _, opt := range opts {
		opt(u)
	}

	userBuilder := ldcontext.NewBuilder(u.userID).Kind(contextKindUser)
	userBuilder.SetString(contextAttributeUserID, u.userID)
	if u.realUserID != "" {
		userBuilder.SetString(contextAttributeRealUserID, u.realUserID)
	}
	if u.accountID != "" {
		userBuilder.SetString(contextAttributeAccountID, u.accountID)
	}
	userContext := userBuilder.Build()

	contextBuilder := u.ContextMultiBuilder()
	contextBuilder.Add(userContext)

	u.ldContext = contextBuilder.Build()

	return *u
}

// UserFromContext extracts the effective user aggregate ID, real user aggregate
// ID, and account aggregate ID from the context. These values are used to
// create a new User object. An error is returned if user identifiers are not
// present in the context.
//
// Deprecated: User is now replaced by EvaluationContext, use EvaluationContextFromContext instead
func UserFromContext(ctx context.Context) (User, error) {
	authenticatedUser, ok := request.AuthenticatedUserFromContext(ctx)
	if !ok {
		return User{}, errors.New("no AuthenticatedUser in supplied context")
	}

	return NewUser(authenticatedUser.UserID, WithUserAccountID(authenticatedUser.CustomerAccountID), WithRealUserID(authenticatedUser.RealUserID)), nil
}
