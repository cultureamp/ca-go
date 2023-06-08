package evaluationcontext

import (
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

const (
	surveyAttributeAccountID = "accountID"
	surveyAttributeSurveyID  = "surveyID"
)

// Survey is a type of context, representing the identifiers and attributes of
// a survey to evaluate a flag against.
type Survey struct {
	key       string
	surveyID  string
	accountID string

	ldContext ldcontext.Context
}

// ToLDContext transforms the context implementation into an LDcontext object that can
// be understood by LaunchDarkly when evaluating a flag.
func (u Survey) ToLDContext() ldcontext.Context {
	return u.ldContext
}

// ToLDUser transforms context into LD context object that can be understood by
// LaunchDarkly when evaluating a flag
//
// Deprecated: use ToLDContext() instead
func (u Survey) ToLDUser() ldcontext.Context {
	return u.ToLDContext()
}

// SurveyOption are functions that can be supplied to configure a new survey with
// additional attributes.
type SurveyOption func(*Survey)

// WithSurveyAccountID configures the survey with the given account ID.
// This is the ID of the currently logged in user's parent account/organization,
// sometimes known as the "account_aggregate_id".
func WithSurveyAccountID(id string) SurveyOption {
	return func(u *Survey) {
		u.accountID = id
	}
}

// NewSurvey returns a new Survey object with the given survey ID, there are no options.
// surveyID is the ID of the currently authenticated survey, and will generally
// be a "survey_aggregate_id".
func NewSurvey(surveyID string, opts ...SurveyOption) Survey {
	u := &Survey{
		key:      surveyID,
		surveyID: surveyID,
	}

	for _, opt := range opts {
		opt(u)
	}

	userBuilder := lduser.NewUserBuilder(u.key)
	userBuilder.Custom(
		surveyAttributeAccountID,
		ldvalue.String(u.accountID))
	userBuilder.Custom(
		surveyAttributeSurveyID,
		ldvalue.String(u.surveyID))
	u.ldUser = userBuilder.Build()

	return *u
}
