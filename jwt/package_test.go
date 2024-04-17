package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPackageEncodeDecode(t *testing.T) {
	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
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

func TestPackageEncodeDecodeNotBeforeExpiryChecks(t *testing.T) {
	now := time.Now()

	expiry := now.Add(10 * time.Hour)
	notBefore := now.Add(-10 * time.Hour)
	okClaims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       expiry,
		NotBefore:       notBefore,
	}

	expiry = now.Add(-1 * time.Hour)
	notBefore = now.Add(-10 * time.Hour)
	expiryClaims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       expiry,
		NotBefore:       notBefore,
	}

	expiry = now.Add(10 * time.Hour)
	notBefore = now.Add(1 * time.Hour)
	notBeforeClaims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       expiry,
		NotBefore:       notBefore,
	}

	testCases := []struct {
		desc                  string
		claims                *StandardClaims
		expectedDecoderErrMsg string
	}{
		{
			desc:                  "Success 1: valid expiry and not before times",
			claims:                okClaims,
			expectedDecoderErrMsg: "",
		},
		{
			desc:                  "Error 1: expired token",
			claims:                expiryClaims,
			expectedDecoderErrMsg: "token is expired",
		},
		{
			desc:                  "Error 2: not available token",
			claims:                notBeforeClaims,
			expectedDecoderErrMsg: "token is not valid yet",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Encode this claim
			token, err := Encode(tC.claims)
			assert.Nil(t, err)

			// Decode it back again
			sc, err := Decode(token)
			if tC.expectedDecoderErrMsg == "" {
				assert.Nil(t, err)
				// check it matches
				assert.Equal(t, "abc123", sc.AccountId)
				assert.Equal(t, "xyz234", sc.RealUserId)
				assert.Equal(t, "xyz345", sc.EffectiveUserId)
				assert.Equal(t, "encoder-name", sc.Issuer)
				assert.Equal(t, "test", sc.Subject)
				assert.Equal(t, []string{"decoder-name"}, sc.Audience)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedDecoderErrMsg)
			}
		})
	}
}

func TestPackageEncodeDecodeWithDifferentEnvVars(t *testing.T) {
	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
	}

	testCases := []struct {
		desc                  string
		encodedKeyId          string
		privateKey            string
		expectedEncoderErrMsg string
		expectedDecoderErrMsg string
	}{
		{
			desc:                  "Success 1: missing env vars defaults to test values",
			encodedKeyId:          "",
			privateKey:            "",
			expectedEncoderErrMsg: "",
			expectedDecoderErrMsg: "",
		},
		{
			desc:                  "Error 1: missing decoder kid",
			encodedKeyId:          "missing-kid",
			privateKey:            testECDSA521PrivateKey,
			expectedEncoderErrMsg: "",
			expectedDecoderErrMsg: "no matching key_id (kid) header",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Setenv("AUTH_PRIVATE_KEY_ID", tC.encodedKeyId)
			privateKeys := ""
			if tC.privateKey != "" {
				b, err := os.ReadFile(filepath.Clean(tC.privateKey))
				assert.Nil(t, err)
				privateKeys = string(b)
			}
			t.Setenv("AUTH_PRIVATE_KEY", privateKeys)

			// Encode this claim
			DefaultJwtEncoder = nil // force re-create
			token, err := Encode(claims)
			if tC.expectedEncoderErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedEncoderErrMsg)
			}

			t.Setenv("AUTH_PUBLIC_JWK_KEYS", "")

			// Decode it back again
			DefaultJwtDecoder = nil // force re-create
			sc, err := Decode(token)
			if tC.expectedDecoderErrMsg == "" {
				assert.Nil(t, err)
				// check it matches
				assert.Equal(t, "abc123", sc.AccountId)
				assert.Equal(t, "xyz234", sc.RealUserId)
				assert.Equal(t, "xyz345", sc.EffectiveUserId)
				assert.Equal(t, "encoder-name", sc.Issuer)
				assert.Equal(t, "test", sc.Subject)
				assert.Equal(t, []string{"decoder-name"}, sc.Audience)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedDecoderErrMsg)
			}
		})
	}

	DefaultJwtEncoder = nil // force re-create for next test
	DefaultJwtDecoder = nil // force re-create for next test
}
