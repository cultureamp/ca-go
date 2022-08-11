package evaluationcontext_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/cultureamp/ca-go/x/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

func TestNewUser(t *testing.T) {
	t.Run("can create an anonymous user", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUser("")
		ldUser := user.ToLDUser()
		assert.True(t, ldUser.GetAnonymous())
	})

	t.Run("can create an anonymous user with session/request key", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUser("my-request-id")
		assert.Equal(t, "my-request-id", user.ToLDUser().GetKey())
	})

	t.Run("can create an anonymous user with subdomain", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUserWithSubdomain("", "cultureamp")
		ldUser := user.ToLDUser()
		// assert.True(t, ldUser.GetAnonymous())

		value, result := ldUser.GetCustom("subdomain")
		assert.Equal(t, ldvalue.String("cultureamp"), value)
		assert.True(t, result)
	})

	t.Run("can create an anonymous user with session/request key and subdomain", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousUserWithSubdomain("my-request-id", "cultureamp")
		assert.Equal(t, "my-request-id", user.ToLDUser().GetKey())

		ldUser := user.ToLDUser()
		value, result := ldUser.GetCustom("subdomain")
		assert.Equal(t, ldvalue.String("cultureamp"), value)
		assert.True(t, result)
	})

	t.Run("can create an identified user", func(t *testing.T) {
		user := evaluationcontext.NewUser("not-a-uuid")
		assertUserAttributes(t, user, "not-a-uuid", "", "")

		user = evaluationcontext.NewUser(
			"not-a-uuid",
			evaluationcontext.WithAccountID("not-a-uuid"),
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

	ldUser := user.ToLDUser()

	assert.Equal(t, userID, ldUser.GetKey())
	assert.Equal(t, userID, ldUser.GetAttribute("userID").StringValue())
	assert.Equal(t, realUserID, ldUser.GetAttribute("realUserID").StringValue())
	assert.Equal(t, accountID, ldUser.GetAttribute("accountID").StringValue())
}
