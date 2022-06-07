package evaluationcontext

import (
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

const (
	surveyAttributeSurveyID = "surveyID"
)

// Survey is a type of context, representing the identifiers and attributes of
// a survey to evaluate a flag against.
type Survey struct {
	key      string
	surveyID string

	ldUser lduser.User
}

func (u Survey) ToLDUser() lduser.User {
	return u.ldUser
}

// NewSurvey returns a new Survey object with the given survey ID, there are no options.
// surveyID is the ID of the currently authenticated survey, and will generally
// be a "survey_aggregate_id".
func NewSurvey(surveyID string) Survey {
	u := &Survey{
		key:      surveyID,
		surveyID: surveyID,
	}

	userBuilder := lduser.NewUserBuilder(u.key)
	userBuilder.Custom(
		surveyAttributeSurveyID,
		ldvalue.String(u.surveyID))
	u.ldUser = userBuilder.Build()

	return *u
}
