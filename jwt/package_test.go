package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
	}

	// Encode this claim
	token, err := Encode(claims)
	assert.Nil(t, err)

	// Decode it back again
	sc, err := Decode(token)
	assert.Nil(t, err)

	// check it matches
	assert.Equal(t, "abc123", sc.AccountId)
	assert.Equal(t, "xyz234", sc.RealUserId)
	assert.Equal(t, "xyz345", sc.EffectiveUserId)
}
