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
	key       string
	surveyID  string
	accountID string

	ldUser lduser.User
}

func (u Survey) ToLDUser() lduser.User {
	return u.ldUser
}

func (u *Survey) SetAccountID(accountID string) {
	u.accountID = accountID
}

func (u *Survey) SetRealUserID(realUserID string) {
	// I guess no need to implement this. Maybe logging error here log.fatal("xxxx")?
	//u.realUserID = realUserID
}

// NewSurvey returns a new Survey object with the given survey ID, there are no options.
// surveyID is the ID of the currently authenticated survey, and will generally
// be a "survey_aggregate_id".
func NewSurvey(surveyID string, opts ...Option) Survey {
	u := &Survey{
		key:      surveyID,
		surveyID: surveyID,
	}

	for _, opt := range opts {
		opt(u)
	}

	userBuilder := lduser.NewUserBuilder(u.key)
	userBuilder.Custom(
		surveyAttributeSurveyID,
		ldvalue.String(u.surveyID))
	userBuilder.Custom(
		userAttributeAccountID,
		ldvalue.String(u.accountID))
	u.ldUser = userBuilder.Build()

	return *u
}
