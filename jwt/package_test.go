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

	defer func() {
		// Force re-create of defaults when test finished
		DefaultJwtEncoder = nil
		DefaultJwtDecoder = nil
	}()

	// Encode this claim
	setupEncodeDecodeForTests()
	token, err := Encode(claims)
	assert.Nil(t, err)

	// Decode it back again
	sc, err := Decode(token)
	assert.Nil(t, err)

	// check it matches
	assert.Equal(t, "abc123", sc.AccountId)
	assert.Equal(t, "xyz234", sc.RealUserId)
	assert.Equal(t, "xyz345", sc.EffectiveUserId)

	// Decode it back again, checking aud, iss and sub all match
	sc, err = Decode(token, MustMatchAudience("decoder-name"), MustMatchIssuer("encoder-name"), MustMatchSubject("test"))
	assert.Nil(t, err)

	// check it matches
	assert.Equal(t, "abc123", sc.AccountId)
	assert.Equal(t, "xyz234", sc.RealUserId)
	assert.Equal(t, "xyz345", sc.EffectiveUserId)

	// Decode it back again, checking aud should fail
	sc, err = Decode(token, MustMatchAudience("incorrect-aud"))
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "token has invalid audience")

	// Decode it back again, checking iss should fail
	sc, err = Decode(token, MustMatchIssuer("incorrect-iss"))
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "token has invalid issuer")

	// Decode it back again, checking sub should fail
	sc, err = Decode(token, MustMatchSubject("incorrect-sub"))
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "token has invalid subject")
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

	defer func() {
		// Force re-create of defaults when test finished
		DefaultJwtEncoder = nil
		DefaultJwtDecoder = nil
	}()

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
			setupEncodeDecodeForTests()
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

	defer func() {
		// Force re-create of defaults when test finished
		DefaultJwtEncoder = nil
		DefaultJwtDecoder = nil
	}()

	testCases := []struct {
		desc                  string
		encodedKeyId          string
		privateKey            string
		expectedEncoderErrMsg string
		expectedDecoderErrMsg string
	}{
		{
			desc:                  "Error 1: missing env vars",
			encodedKeyId:          "",
			privateKey:            "",
			expectedEncoderErrMsg: "error loading jwk encoder, maybe missing env vars",
			expectedDecoderErrMsg: "token is malformed: token contains an invalid number of segments",
		},
		{
			desc:                  "Error 2: missing decoder kid",
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

func setupEncodeDecodeForTests() {
	DefaultJwtEncoder = nil
	wb, _ := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
	privKey := string(wb)
	os.Setenv("AUTH_PRIVATE_KEY", privKey)
	os.Setenv("AUTH_PRIVATE_KEY_ID", webGatewayKid)

	DefaultJwtDecoder = nil
	db, _ := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
	jwkKeys := string(db)
	os.Setenv("AUTH_PUBLIC_JWK_KEYS", jwkKeys)
}
