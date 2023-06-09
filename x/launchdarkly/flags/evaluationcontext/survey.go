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
//
// Deprecated: use EvaluationContext instead
type Survey struct {
	EvaluationContext
}

// ToLDUser transforms context into LD context object that can be understood by
// LaunchDarkly when evaluating a flag
//
// Deprecated: Survey is deprecated, use EvaluationContext with ToLDContext instead
func (u Survey) ToLDUser() ldcontext.Context {
	return u.ToLDContext()
}

// SurveyOption are functions that can be supplied to configure a new survey with
// additional attributes.
//
// Deprecated: use ContextOption with EvaluationContext instead
type SurveyOption func(*Survey)

// WithSurveyAccountID configures the survey with the given account ID.
// This is the ID of the currently logged in user's parent account/organization,
// sometimes known as the "account_aggregate_id".
//
// Deprecated: use ContextOptions with EvaluationContext instead
func WithSurveyAccountID(id string) SurveyOption {
	return func(u *Survey) {
		u.accountID = id
	}
}

// NewSurvey returns a new Survey object with the given survey ID, there are no options.
// surveyID is the ID of the currently authenticated survey, and will generally
// be a "survey_aggregate_id".
//
// Deprecated: use NewEvaluationContext with surveyID instead
func NewSurvey(surveyID string, opts ...SurveyOption) Survey {
	s := &Survey{}
	s.surveyID = surveyID
	for _, opt := range opts {
		opt(s)
	}

	// for backwards compatibility
	userContext := ldcontext.NewBuilder(s.surveyID).
		Kind(contextKindUser).
		SetString(surveyAttributeSurveyID, s.surveyID).
		SetString(surveyAttributeAccountID, s.accountID).
		Build()
	contextBuilder := s.ContextMultiBuilder()
	contextBuilder.Add(userContext)
	s.ldContext = contextBuilder.Build()
	return *s
}
