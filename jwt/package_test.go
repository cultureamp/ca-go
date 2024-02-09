package jwt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestEncodeDecodeWhenMissingEnvVars(t *testing.T) {
	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
	}

	testCases := []struct {
		desc                  string
		privateKeyId          string
		privateKey            string
		expectedEncoderErrMsg string
		expectedDecoderErrMsg string
	}{
		{
			desc:                  "Success 1: missing env vars defaults to test values",
			privateKeyId:          "",
			privateKey:            "",
			expectedEncoderErrMsg: "",
			expectedDecoderErrMsg: "",
		},
		{
			desc:                  "Error 1: missing decoder kid",
			privateKeyId:          "missing-kid",
			privateKey:            testOtherAuthPrivateKey,
			expectedEncoderErrMsg: "",
			expectedDecoderErrMsg: "token signature is invalid",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Setenv("AUTH_PRIVATE_KEY_ID", tC.privateKeyId)
			privateKeys := ""
			if tC.privateKey != "" {
				b, err := os.ReadFile(filepath.Clean(tC.privateKey))
				require.NoError(t, err)
				privateKeys = string(b)
			}
			t.Setenv("AUTH_PRIVATE_KEY", privateKeys)

			// Encode this claim
			encoder := getEncoderInstance()
			token, err := encoder.Encode(claims)
			if tC.expectedEncoderErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedEncoderErrMsg)
			}

			t.Setenv("AUTH_PUBLIC_DEFAULT_KEY_ID", "")
			t.Setenv("AUTH_PUBLIC_JWK_KEYS", "")

			// Decode it back again
			decoder := getDecoderInstance()
			sc, err := decoder.Decode(token)
			if tC.expectedDecoderErrMsg == "" {
				assert.Nil(t, err)
				// check it matches
				assert.Equal(t, "abc123", sc.AccountId)
				assert.Equal(t, "xyz234", sc.RealUserId)
				assert.Equal(t, "xyz345", sc.EffectiveUserId)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedDecoderErrMsg)
			}
		})
	}
}