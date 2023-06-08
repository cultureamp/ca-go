package evaluationcontext

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/google/uuid"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

const (
	accountAttributeAccountID = "accountID"
)

// Account is a type of context, representing the identifiers and attributes of
// a singular account to evaluate a flag against.
type Account struct {
	key       string
	accountID string
	subdomain string

	ldContext ldcontext.Context
}

// ToLDContext transforms the context implementation into an LDContext object that can
// be understood by LaunchDarkly when evaluating a flag.
func (a Account) ToLDContext() ldcontext.Context {
	return a.ldContext
}

// ToLDUser transforms the context implementation into an LDUser object that can
// be understood by LaunchDarkly when evaluating a flag.
//
// Deprecated: use ToLDContext() instead
func (a Account) ToLDUser() ldcontext.Context {
	return a.ToLDContext()
}

// AccountOption are functions that can be supplied to configure a new Account with
// additional attributes.
type AccountOption func(*Account)

// WithSubdomain configures the account with the given subdomain.
// This is the subdomain of the currently logged in account/organization
func WithSubdomain(subdomain string) AccountOption {
	return func(a *Account) {
		a.subdomain = subdomain
	}
}

// NewAnonymousAccount returns a context object suitable for use in unauthenticated
// requests or requests with no access to context identifiers.
// Provide a unique session or request identifier as the key if possible. If the
// key is empty, it will default to an uuid so percentage rollouts will still apply.
func NewAnonymousAccount(key string) Account {
	if key == "" {
		key = uuid.NewString()
	}

	return Account{
		key: key,
		ldContext: ldcontext.NewBuilder(key).
			SetString(accountAttributeAccountID, key).
			Kind("account").
			Anonymous(true).
			Build(),
	}
}

// NewAnonymousAccountWithSubdomain returns an account object suitable for use in unauthenticated with known subdomain
// requests or requests with no access to context identifiers.
// Provide a unique session or request identifier as the key if possible. If the
// key is empty, it will default to an uuid so percentage rollouts will still apply.
func NewAnonymousAccountWithSubdomain(key string, subdomain string) Account {
	if key == "" {
		key = uuid.NewString()
	}

	return Account{
		key:       key,
		ldContext: ldcontext.NewBuilder(key).Anonymous(true).Kind("account").Name(subdomain).Build(),
	}
}

// NewAccount returns a new account object with the given account ID and options.
// accountID is the ID of the current account, and will generally
// be an "account_aggregate_id".
func NewAccount(accountID string, opts ...AccountOption) Account {
	a := &Account{
		key:       accountID,
		accountID: accountID,
	}

	for _, opt := range opts {
		opt(a)
	}

	accountBuilder := ldcontext.NewBuilder(a.accountID).Kind("account").SetString(accountAttributeAccountID, accountID)
	if a.subdomain != "" {
		accountBuilder.Name(a.subdomain)
	}
	a.ldContext = accountBuilder.Build()
	return *a
}

// AccountFromContext extracts the effective account aggregate ID from the context. These values are used to
// create a new Account object. An error is returned if account identifiers are not
// present in the context.
func AccountFromContext(ctx context.Context) (Account, error) {
	authenticatedUser, ok := request.AuthenticatedUserFromContext(ctx)
	if !ok {
		return Account{}, errors.New("no AuthenticatedUser in supplied context")
	}

	return NewAccount(authenticatedUser.CustomerAccountID), nil
}
