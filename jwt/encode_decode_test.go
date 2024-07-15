package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testAuthJwks           string = "./testKeys/development.jwks"
	testRSA256PrivateKey   string = "./testKeys/jwt-rsa256-test.key"
	testRSA384PrivateKey   string = "./testKeys/jwt-rsa384-test.key"
	testRSA512PrivateKey   string = "./testKeys/jwt-rsa512-test.key"
	testECDSA521PrivateKey string = "./testKeys/jwt-ecdsa521-test.key"
	testECDSA384PrivateKey string = "./testKeys/jwt-ecdsa384-test.key"
	testECDSA256PrivateKey string = "./testKeys/jwt-ecdsa256-test.key"
)

// useful to create RS256 test tokens https://jwt.io/
// useful for PEM to JWKS https://jwkset.com/generate

func TestEncodeDecode(t *testing.T) {
	claims := &StandardClaims{
		AccountID:       "abc123",
		RealUserID:      "xyz234",
		EffectiveUserID: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
	}

	testCases := []struct {
		desc    string
		privkey string
		kid     string
	}{
		{
			desc:    "Success 1: RSA 256 Key",
			privkey: testRSA256PrivateKey,
			kid:     "rsa-256",
		},
		{
			desc:    "Success 2: RSA 384 Key",
			privkey: testRSA384PrivateKey,
			kid:     "rsa-384",
		},
		{
			desc:    "Success 3: RSA 512 Key",
			privkey: testRSA512PrivateKey,
			kid:     "rsa-512",
		},
		{
			desc:    "Success 4: ECDSA 256 Key",
			privkey: testECDSA256PrivateKey,
			kid:     "ecdsa-256",
		},
		{
			desc:    "Success 5: ECDSA 384 Key",
			privkey: testECDSA384PrivateKey,
			kid:     "ecdsa-384",
		},
		{
			desc:    "Success 6: ECDSA 521 Key",
			privkey: testECDSA521PrivateKey,
			kid:     "ecdsa-test",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// 1. Load and create encoder
			b, err := os.ReadFile(filepath.Clean(tC.privkey))
			assert.Nil(t, err)
			privKey := string(b)

			encoder, err := NewJwtEncoder(func() (string, string) { return privKey, tC.kid })
			assert.Nil(t, err)
			assert.NotNil(t, encoder)

			// 2. Load and create decoder
			b, err = os.ReadFile(filepath.Clean(testAuthJwks))
			assert.Nil(t, err)
			jwks := string(b)

			decoder, err := NewJwtDecoder(func() string { return jwks })
			assert.Nil(t, err)
			assert.NotNil(t, decoder)

			// 3. Encode the claims and then Decode the token
			token, err := encoder.Encode(claims)
			assert.Nil(t, err)
			// fmt.Printf("Token: '%s'", token)

			actual, err := decoder.Decode(token)
			assert.Nil(t, err)

			// 4. Assert its the same
			assert.Equal(t, claims.AccountID, actual.AccountID)
			assert.Equal(t, claims.RealUserID, actual.RealUserID)
			assert.Equal(t, claims.EffectiveUserID, actual.EffectiveUserID)
			assert.Equal(t, claims.Issuer, actual.Issuer)
			assert.Equal(t, claims.Subject, actual.Subject)
			assert.Equal(t, claims.Audience, actual.Audience)
			assert.Equal(t, claims.ExpiresAt.Year(), actual.ExpiresAt.Year())
		})
	}
}
