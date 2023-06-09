package evaluationcontext_test

import (
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/stretchr/testify/assert"
)

func TestNewSurvey(t *testing.T) {
	t.Run("can create a survey", func(t *testing.T) {
		survey := evaluationcontext.NewSurvey("not-a-uuid")
		assertSurveyAttributes(t, survey, "not-a-uuid", "")

		survey = evaluationcontext.NewSurvey(
			"not-a-uuid",
		)
		assertSurveyAttributes(t, survey, "not-a-uuid", "")
	})

	t.Run("can create a survey with account ID", func(t *testing.T) {
		survey := evaluationcontext.NewSurvey("not-a-uuid")
		assertSurveyAttributes(t, survey, "not-a-uuid", "")

		survey = evaluationcontext.NewSurvey(
			"not-a-uuid",
			evaluationcontext.WithSurveyAccountID("not-a-uuid"))
		assertSurveyAttributes(t, survey, "not-a-uuid", "not-a-uuid")
	})
}

func assertSurveyAttributes(t *testing.T, survey evaluationcontext.Survey, surveyID string, accountID string) {
	t.Helper()

	ldContext := survey.ToLDContext()
	ldSurvey := ldContext.IndividualContextByKind("survey")
	ldAccount := ldContext.IndividualContextByKind("account")
	ldUser := ldContext.IndividualContextByKind("user")
	assert.Equal(t, surveyID, ldSurvey.Key())
	assert.Equal(t, surveyID, ldUser.GetValue("surveyID").StringValue())
	assert.Equal(t, accountID, ldUser.GetValue("accountID").StringValue())
	assert.Equal(t, accountID, ldAccount.Key())
}
