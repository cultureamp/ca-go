package evaluationcontext_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/cultureamp/ca-go/x/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	t.Run("can create an anonymous account", func(t *testing.T) {
		account := evaluationcontext.NewAnonymousAccount("")
		ldContext := account.ToLDContext()
		assert.True(t, ldContext.Anonymous())
	})

	t.Run("can create an anonymous account with session/request key", func(t *testing.T) {
		account := evaluationcontext.NewAnonymousAccount("my-request-id")
		assert.Equal(t, "my-request-id", account.ToLDContext().Key())
	})

	t.Run("can create an anonymous account with subdomain", func(t *testing.T) {
		account := evaluationcontext.NewAnonymousAccountWithSubdomain("", "cultureamp")
		ldContext := account.ToLDContext()
		assert.True(t, ldContext.Anonymous())
		assert.Equal(t, "cultureamp", ldContext.Name().StringValue())
	})

	t.Run("can create an anonymous account with session/request key and subdomain", func(t *testing.T) {
		user := evaluationcontext.NewAnonymousAccountWithSubdomain("my-request-id", "cultureamp")

		ldContext := user.ToLDContext()
		assert.Equal(t, "my-request-id", ldContext.Key())
		assert.Equal(t, "cultureamp", ldContext.Name().StringValue())
	})

	t.Run("can create an identified account", func(t *testing.T) {
		account := evaluationcontext.NewAccount("not-a-uuid")
		assertAccountAttributes(t, account, "not-a-uuid", "")

		account = evaluationcontext.NewAccount(
			"not-a-uuid", evaluationcontext.WithSubdomain("not-a-subdomain"))
		assertAccountAttributes(t, account, "not-a-uuid", "not-a-subdomain")
	})

	t.Run("can create an account from context", func(t *testing.T) {
		user := request.AuthenticatedUser{
			CustomerAccountID: "123",
			RealUserID:        "456",
			UserID:            "789",
		}
		ctx := context.Background()

		ctx = request.ContextWithAuthenticatedUser(ctx, user)

		flagsAccount, err := evaluationcontext.AccountFromContext(ctx)
		require.NoError(t, err)
		assertAccountAttributes(t, flagsAccount, "123", "")
	})
}

func assertAccountAttributes(t *testing.T, account evaluationcontext.Account, accountID string, subdomain string) {
	t.Helper()

	ldContext := account.ToLDContext()
	ldAccount := ldContext.IndividualContextByKind("account")
	assert.Equal(t, accountID, ldAccount.Key())
	assert.Equal(t, subdomain, ldAccount.Name().StringValue())
}
