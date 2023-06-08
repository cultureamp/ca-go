package evaluationcontext_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/cultureamp/ca-go/x/request"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("can create an anonymous user", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUser("")
		ldContext := user.ToLDContext()
		assert.True(t, ldContext.Anonymous())
	})

	t.Run("can create an anonymous user with session/request key", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUser("my-request-id")
		assert.Equal(t, "my-request-id", user.ToLDContext().Key())
	})

	t.Run("can create an anonymous user with subdomain", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUserWithSubdomain("", "cultureamp")
		ldContext := user.ToLDContext()
		// loop through multicontexts to check they are anon
		if ldContext.IndividualContextCount() > 1 {
			preallocContexts := make([]ldcontext.Context, 0, 2)
			for _, individualContext := range ldContext.GetAllIndividualContexts(preallocContexts) {
				assert.True(t, individualContext.Anonymous())
			}
		} else {
			ldContext.Anonymous()
		}
		userSubdomain := ldContext.IndividualContextByKind("user").GetValue("subdomain")
		assert.Equal(t, "cultureamp", userSubdomain.StringValue())
		value := ldContext.IndividualContextByKind("account").Name()
		assert.Equal(t, "cultureamp", value.StringValue())
	})

	t.Run("can create an anonymous user with session/request key and subdomain", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUserWithSubdomain("my-request-id", "cultureamp")

		ldContext := user.ToLDContext()
		assert.Equal(t, "my-request-id", ldContext.IndividualContextByKind("account").Key())
		value := ldContext.IndividualContextByKind("account").Name()
		assert.Equal(t, "cultureamp", value.StringValue())
	})

	t.Run("can create an identified user", func(t *testing.T) {
		user := evaluationcontext.NewUser("not-a-uuid")
		assertUserAttributes(t, user, "not-a-uuid", "", "")

		user = evaluationcontext.NewUser(
			"not-a-uuid",
			evaluationcontext.WithUserAccountID("not-a-uuid"),
			evaluationcontext.WithRealUserID("not-a-uuid"))
		assertUserAttributes(t, user, "not-a-uuid", "not-a-uuid", "not-a-uuid")
	})

	t.Run("can create a user from context", func(t *testing.T) {
		user := request.AuthenticatedUser{
			CustomerAccountID: "123",
			RealUserID:        "456",
			UserID:            "789",
		}
		ctx := context.Background()

		ctx = request.ContextWithAuthenticatedUser(ctx, user)

		flagsUser, err := evaluationcontext.UserFromContext(ctx)
		require.NoError(t, err)
		assertUserAttributes(t, flagsUser, "789", "456", "123")
	})
}

func assertUserAttributes(t *testing.T, user evaluationcontext.User, userID, realUserID, accountID string) {
	t.Helper()

	ldContext := user.ToLDContext()
	ldUser := ldContext.IndividualContextByKind("user")
	ldAccount := ldContext.IndividualContextByKind("account")
	assert.Equal(t, userID, ldUser.Key())
	assert.Equal(t, userID, ldUser.GetValue("userID").StringValue())
	assert.Equal(t, realUserID, ldUser.GetValue("realUserID").StringValue())
	assert.Equal(t, accountID, ldUser.GetValue("accountID").StringValue())
	assert.Equal(t, accountID, ldAccount.Key())
}
