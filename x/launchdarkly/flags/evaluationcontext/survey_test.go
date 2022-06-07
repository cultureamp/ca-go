package evaluationcontext_test

import (
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/stretchr/testify/assert"
)

func TestNewSurvey(t *testing.T) {
	t.Run("can create a survey", func(t *testing.T) {
		survey := evaluationcontext.NewSurvey("not-a-uuid")
		assertSurveyAttributes(t, survey, "not-a-uuid")

		survey = evaluationcontext.NewSurvey(
			"not-a-uuid",
		)
		assertSurveyAttributes(t, survey, "not-a-uuid")
	})
}

func assertSurveyAttributes(t *testing.T, survey evaluationcontext.Survey, surveyID string) {
	t.Helper()

	ldUser := survey.ToLDUser()

	assert.Equal(t, surveyID, ldUser.GetKey())
	assert.Equal(t, surveyID, ldUser.GetAttribute("surveyID").StringValue())
}
