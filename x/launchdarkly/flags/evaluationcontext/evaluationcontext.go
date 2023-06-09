package evaluationcontext

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/google/uuid"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

// Context represents a set of attributes which a flag is evaluated against. The
// only contexts supported now are EvaluationContext, User and Survey
type Context interface {
	// ToLDContext transforms the context implementation into an LDContext object that can
	// be understood by LaunchDarkly when evaluating a flag.
	ToLDContext() ldcontext.Context
}

const (
	contextAttributeUserID     = "userID"
	contextAttributeAccountID  = "accountID"
	contextAttributeRealUserID = "realUserID"
	contextAttributeSubdomain  = "subdomain"
	contextKindAccount         = "account"
	contextKindSurvey          = "survey"
	contextKindUser            = "user"
)

// EvaluationContext is the context that is evaluating a flag, it contains all the attributes required for targeting
type EvaluationContext struct {
	userID     string
	realUserID string
	accountID  string
	surveyID   string

	ldContext ldcontext.Context
}

// ToLDContext transforms the context implementation into an LDContext object that can
// be understood by LaunchDarkly when evaluating a flag.
func (e EvaluationContext) ToLDContext() ldcontext.Context {
	return e.ldContext
}

func (c EvaluationContext) ContextMultiBuilder() *ldcontext.MultiBuilder {
	contextBuilder := ldcontext.NewMultiBuilder()
	if c.realUserID != "" && c.userID == "" {
		contextBuilder.Add(ldcontext.NewBuilder(c.realUserID).Kind(contextKindUser).SetString(contextAttributeRealUserID, c.realUserID).Build())
	}
	if c.accountID != "" {
		accountContext := ldcontext.NewBuilder(c.accountID).Kind(contextKindAccount).Build()
		contextBuilder.Add(accountContext)
	}
	if c.surveyID != "" {
		surveyContext := ldcontext.NewBuilder(c.surveyID).Kind(contextKindSurvey).Build()
		contextBuilder.Add(surveyContext)
	}

	return contextBuilder
}

// ContextOption are functions that can be supplied to configure a new evaluation context with
// additional attributes.
type ContextOption func(*EvaluationContext)

// WithUserID configures the context with the given userID.
// userID is the ID of the currently authenticated user, and will generally
// be a "user_aggregate_id".
func WithUserID(id string) ContextOption {
	return func(e *EvaluationContext) {
		e.userID = id
	}
}

// WithAccountID configures the user with the given account ID.
// This is the ID of the currently logged in user's parent account/organization,
// sometimes known as the "account_aggregate_id".
func WithAccountID(id string) ContextOption {
	return func(e *EvaluationContext) {
		e.accountID = id
	}
}

// WithRealUserID configures the user with the given real user ID.
// This is the ID of the user who is currently impersonating the current user.
func WithContextRealUserID(id string) ContextOption {
	return func(e *EvaluationContext) {
		e.realUserID = id
	}
}

// WithSurveyID configures the context with the given survey ID.
// This is the ID of the related survey to target against.
func WithSurveyID(id string) ContextOption {
	return func(e *EvaluationContext) {
		e.surveyID = id
	}
}

// NewAnonymousContextWithSubdomain returns an evaluation context object suitable for use in unauthenticated
// environments with known subdomain requests or requests with no access to user identifiers.
// Provide a unique session or request identifier as the key if possible. If the
// key is empty, it will default to an uuid so percentage rollouts will still apply.
func NewAnonymousContextWithSubdomain(key string, subdomain string) EvaluationContext {
	if key == "" {
		key = uuid.NewString()
	}

	return EvaluationContext{
		accountID: key,
		// not using Name attribute for subdomain to avoid clashing with segment sync named account contexts
		ldContext: ldcontext.NewBuilder(key).Anonymous(true).Kind("account").SetString(contextAttributeSubdomain, subdomain).Build(),
	}
}

// NewEvaluationContext returns a new context object with the given options.
// As many options as are available should be provided for increased targeting ability
// If no options are provided, an anonymous context with a randomly generated key will be returned. This is to be used
// for unauthenticated users where no information is available.
func NewEvaluationContext(opts ...ContextOption) EvaluationContext {
	c := &EvaluationContext{}

	// if no options provided then context is anonymous
	if len(opts) == 0 {
		key := uuid.NewString()
		c.userID = key
		c.ldContext = ldcontext.NewBuilder(key).Anonymous(true).Build()
		return *c
	}

	// apply the options
	for _, opt := range opts {
		opt(c)
	}

	// Separating the user context out of ContextMultiBuilder to avoid duplicated user contexts for legacy
	// TODO: move this logic back into the function when User/Survey is removed
	contextBuilder := c.ContextMultiBuilder()
	if c.userID != "" {
		userContext := ldcontext.NewBuilder(c.userID).Kind(contextKindUser)
		if c.realUserID != "" {
			userContext.SetString(contextAttributeRealUserID, c.realUserID)
		}
		contextBuilder.Add(userContext.Build())
	}
	c.ldContext = contextBuilder.Build()

	return *c
}

// EvaluationContextFromContext extracts the effective user aggregate ID, real user aggregate
// ID, and account aggregate ID from the context. These values are used to
// create a new EvaluationContext object. An error is returned if user identifiers are not
// present in the context.
func EvaluationContextFromContext(ctx context.Context) (EvaluationContext, error) {
	authenticatedUser, ok := request.AuthenticatedUserFromContext(ctx)
	if !ok {
		return EvaluationContext{}, errors.New("no AuthenticatedUser in supplied context")
	}

	return NewEvaluationContext(WithUserID(authenticatedUser.UserID), WithAccountID(authenticatedUser.CustomerAccountID), WithContextRealUserID(authenticatedUser.RealUserID)), nil
}
