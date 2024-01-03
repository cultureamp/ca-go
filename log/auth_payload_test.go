package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthPayload(t *testing.T) {
	ctx := context.Background()

	// create a jwt payload
	auth_expected := AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	ctx = ContextWithAuthPayload(ctx, auth_expected)
	payload, ok := AuthPayloadFromContext(ctx)
	assert.True(t, ok)
	assert.NotNil(t, payload)
	assert.Equal(t, "account_123_id", payload.CustomerAccountID)
	assert.Equal(t, "real_456_id", payload.RealUserID)
	assert.Equal(t, "user_789_id", payload.UserID)
}
