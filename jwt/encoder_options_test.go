package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncoderOptions(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testRSA256PrivateKey))
	assert.Nil(t, err)
	privKey := string(b)

	i := 0
	encoderKeyFunc := func() (string, string) {
		i++
		return privKey, "rsa-256"
	}

	encoder, err := NewJwtEncoder(encoderKeyFunc, WithEncoderCacheExpiry(100*time.Millisecond, 50*time.Millisecond))
	assert.Nil(t, err)
	assert.NotNil(t, encoder)

	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	for j := 0; j < 10; j++ {
		token, err := encoder.Encode(claims)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		time.Sleep(50 * time.Millisecond)
	}

	assert.Greater(t, i, 1)
}
