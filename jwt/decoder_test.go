package jwt

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	day  = 24 * time.Hour
	year = 365 * day // approx

	decoderAuthKey       string = "./testKeys/jwt.rs256.key.development.pub"
	decoderSecondAuthKey string = "./testKeys/jwt.rs256.key.development.extra_1.pub"
	decoderThirdAuthKey  string = "./testKeys/jwt.rs256.key.development.extra_2.pub"
	decoderAuthJwks      string = "./testKeys/development.jwks"
)

func TestNewDecoderSuccess(t *testing.T) {
	pubKeyBytes, err := os.ReadFile(decoderAuthKey)
	require.NoError(t, err)

	pubSecondKeyBytes, err := os.ReadFile(decoderSecondAuthKey)
	require.NoError(t, err)

	pubJwkKeyBytes, err := os.ReadFile(decoderAuthJwks)
	require.NoError(t, err)

	decoder, err := NewJwtDecoder(string(pubKeyBytes), string(pubSecondKeyBytes), string(pubJwkKeyBytes))
	assert.Nil(t, err)
	assert.NotNil(t, decoder)
}

func TestNewDecoderErrors(t *testing.T) {
	invalidPublicKey := "invalid-public-key"

	pubKeyBytes, err := os.ReadFile(decoderAuthKey)
	require.NoError(t, err)

	pubSecondKeyBytes, err := os.ReadFile(decoderSecondAuthKey)
	require.NoError(t, err)

	testCases := []struct {
		desc           string
		pubKey         string
		secondPubKey   string
		jwks           string
		expectedErrMsg string
	}{
		{
			desc:           "Error 1: missing key",
			pubKey:         "",
			secondPubKey:   "",
			jwks:           "",
			expectedErrMsg: "invalid key",
		},
		{
			desc:           "Error 2: bad key",
			pubKey:         invalidPublicKey,
			secondPubKey:   invalidPublicKey,
			jwks:           invalidPublicKey,
			expectedErrMsg: "invalid key",
		},
		{
			desc:           "Error 3: keys ok, JWKS json bad",
			pubKey:         string(pubKeyBytes),
			secondPubKey:   string(pubSecondKeyBytes),
			jwks:           "{\"bad\": \"jwks-json\" }",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
		{
			desc:           "Error 3: keys ok, JWKS json invalid",
			pubKey:         string(pubKeyBytes),
			secondPubKey:   string(pubSecondKeyBytes),
			jwks:           "invalid JSON",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoder(tC.pubKey, tC.secondPubKey, tC.jwks)
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tC.expectedErrMsg)
			assert.Nil(t, decoder)
		})
	}
}

func TestDecoderDecodeAllClaims(t *testing.T) {
	pubKeyBytes, err := os.ReadFile(decoderAuthKey)
	require.NoError(t, err)

	pubSecondKeyBytes, err := os.ReadFile(decoderSecondAuthKey)
	require.NoError(t, err)

	pubJwksBytes, err := os.ReadFile(decoderAuthJwks)
	require.NoError(t, err)

	decoder, err := NewJwtDecoder(string(pubKeyBytes), string(pubSecondKeyBytes), string(pubJwksBytes))
	assert.Nil(t, err)
	assert.NotNil(t, decoder)

	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiZXhwIjoxOTAzOTMwNzA0LCJpYXQiOjE1ODg1NzA3MDR9.XGm34FDIgtBFvx5yC2HTUu-cf3DaQI4TmIBVLx0H7y89oNVNWJaKA3dLvWS0oOZoYIuGhj6GzPREBEmou2f9JsUerqnc-_Tf8oekFZWU7kEfzu9ECBiSWPk7ljPJeZLbau62sSqD7rYb-m3v1mohqz4tKJ_7leWu9L1uHHliC7YGlSRl1ptVDllJjKXKjOg9ifeGSXDEMeU35KgCFwIwKdu8WmCTd8ztLSKEnLT1OSaRZ7MSpmHQ4wUZtS6qvhLBiquvHub9KdQmc4mYWLmfKdDiR5DH-aswJFGLVu3yisFRY8uSfeTPQRhQXd_UfdgifCTXdWTnCvNZT-BxULYG-5mlvAFu-JInTga_9-r-wHRzFD1SrcKjuECF7vUG8czxGNE4sPjFrGVyBxE6fzzcFsdrhdqS-LB_shVoG940fD-ecAhXQZ9VKgr-rmCvmxuv5vYI2HoMfg9j_-zeXkucKxvPYvDQZYMdeW4wFsUORliGplThoHEeRQxTX8d_gvZFCy_gGg0H57FmJwCRymWk9v29s6uyHUMor_r-e7e6ZlShFBrCPAghXL04S9IFJUxUv30wNie8aaSyvPuiTqCgGiEwF_20ZaHCgYX0zupdGm4pHTyJrx2wv31yZ4VZYt8tKjEW6-BlB0nxzLGk5OUN83vq-RzH-92WmY5kMndF6Jo"
	claim, err := decoder.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "abc123", claim.AccountId)
	assert.Equal(t, "xyz234", claim.RealUserId)
	assert.Equal(t, "xyz345", claim.EffectiveUserId)
	assert.Equal(t, 2030, claim.Expiry.Year())
}
