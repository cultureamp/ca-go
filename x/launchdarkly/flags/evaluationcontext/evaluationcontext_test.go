package evaluationcontext_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/cultureamp/ca-go/x/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEvaluationContext(t *testing.T) {
	t.Run("can create an anonymous context", func(t *testing.T) {
		evalcontext := evaluationcontext.NewEvaluationContext()
		ldContext := evalcontext.ToLDContext()
		assert.True(t, ldContext.IndividualContextByKind("user").Anonymous())
	})

	t.Run("can create a context with only userID", func(t *testing.T) {
		evalcontext := evaluationcontext.NewEvaluationContext(evaluationcontext.WithUserID("not-a-uuid"))
		assertContextAttributes(t, evalcontext, "not-a-uuid", "", "", "", 1)
	})

	t.Run("can create a context only realUserID", func(t *testing.T) {
		evalcontext := evaluationcontext.NewEvaluationContext(
			evaluationcontext.WithContextRealUserID("not-a-real-user-uuid"))
		assertContextAttributes(t, evalcontext, "not-a-real-user-uuid", "not-a-real-user-uuid", "", "", 1)
	})

	t.Run("can create a context with all attributes", func(t *testing.T) {
		evalcontext := evaluationcontext.NewEvaluationContext(
			evaluationcontext.WithUserID("not-a-user-uuid"),
			evaluationcontext.WithAccountID("not-a-account-uuid"),
			evaluationcontext.WithContextRealUserID("not-a-real-user-uuid"),
			evaluationcontext.WithSurveyID("not-a-survey-uuid"))
		assertContextAttributes(t, evalcontext, "not-a-user-uuid", "not-a-real-user-uuid", "not-a-account-uuid", "not-a-survey-uuid", 3)
	})
}
func TestEvaluationContextFromContext(t *testing.T) {
	t.Run("can create an evaluation context from context", func(t *testing.T) {
		user := request.AuthenticatedUser{
			CustomerAccountID: "123",
			RealUserID:        "456",
			UserID:            "789",
		}
		ctx := context.Background()

		ctx = request.ContextWithAuthenticatedUser(ctx, user)

		flagsEvalContext, err := evaluationcontext.EvaluationContextFromContext(ctx)
		require.NoError(t, err)
		assertContextAttributes(t, flagsEvalContext, "789", "456", "123", "", 2)
	})
}

func TestNewAnonymousContextWithSubdomain(t *testing.T) {
	t.Run("can create an anonymous context with subdomain", func(t *testing.T) {
		evalcontext := evaluationcontext.NewAnonymousContextWithSubdomain("", "cultureamp")
		ldContext := evalcontext.ToLDContext()
		// loop through multicontexts to check they are anon

		assert.True(t, ldContext.Anonymous())
		value := ldContext.GetValue("subdomain")
		assert.Equal(t, "cultureamp", value.StringValue())
	})

	t.Run("can create an anonymous context with session/request key and subdomain", func(t *testing.T) {
		evalcontext := evaluationcontext.NewAnonymousContextWithSubdomain("my-request-id", "cultureamp")

		ldContext := evalcontext.ToLDContext()
		assert.True(t, ldContext.Anonymous())
		assert.Equal(t, "my-request-id", ldContext.Key())
		value := ldContext.GetValue("subdomain")
		assert.Equal(t, "cultureamp", value.StringValue())
	})
}

func assertContextAttributes(t *testing.T, context evaluationcontext.EvaluationContext, userID, realUserID, accountID string, surveyID string, expectedContextNo int) {
	t.Helper()

	ldContext := context.ToLDContext()
	assert.Equal(t, expectedContextNo, ldContext.IndividualContextCount())
	ldUser := ldContext.IndividualContextByKind("user")
	ldAccount := ldContext.IndividualContextByKind("account")
	ldSurvey := ldContext.IndividualContextByKind("survey")
	assert.Equal(t, userID, ldUser.Key())
	assert.Equal(t, realUserID, ldUser.GetValue("realUserID").StringValue())
	assert.Equal(t, surveyID, ldSurvey.Key())
	assert.Equal(t, accountID, ldAccount.Key())
}
