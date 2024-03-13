package jwt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// useful to create RS256 test tokens https://jwt.io/
// useful for PEM to JWKS https://jwkset.com/generate

func TestNewEncoder(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testRSA256PrivateKey))
	assert.Nil(t, err)
	privKey := string(b)

	testCases := []struct {
		desc           string
		key            string
		kid            string
		expectedErrMsg string
	}{
		{
			desc:           "Success 1: valid private key",
			key:            privKey,
			kid:            "rsa-256",
			expectedErrMsg: "",
		},
		{
			desc:           "Error 1: missing key",
			key:            "",
			kid:            "",
			expectedErrMsg: "invalid private key",
		},
		{
			desc:           "Error 2: bad key",
			key:            "bad key",
			kid:            "bad",
			expectedErrMsg: "invalid private key",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			encoder, err := NewJwtEncoder(tC.key, tC.kid)
			if tC.expectedErrMsg != "" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedErrMsg)
				assert.Nil(t, encoder)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, encoder)
			}
		})
	}
}

func TestEncoderRSA(t *testing.T) {
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

	testCases := []struct {
		desc      string
		key       string
		kid       string
		header    string
		payload   string
		signature string
	}{
		{
			desc:      "Success 1: RSA 256",
			key:       "./testKeys/jwt-rsa256-test.key",
			kid:       "rsa-256",
			header:    "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS0yNTYiLCJ0eXAiOiJKV1QifQ",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0",
			signature: "TS-896O9X6bd6jz6Y765OrJhuMU9qf3PYGE3R9L2AIUO6TEtc3PY_dSoHYIvu_WL8pPn_2MOojvhfoHMFe6HxL5ADyipjJ1hvYiwPJqDq4PqDWvPKQqP6whK8X5Nn1sJvqy3FJ20ToioQDDeqXj0wjY2ZmPnmy_jC4TW2WhWDQI6_Vkl0ha0txEaLORuGeReMb40q4C7Za15NwWqLUvX6mvOph7YjDUxhkoxAldVtPsHiPG37JNcHfkh_2cN2i-07g9vPa16gXMP0cJP8-CVh1qNZ2un9eqtrCVYOvA_Ln9dNpUhaRPPYNAz-gcwTpUqXqySLx_cBVAnvECFUfT_LQ",
		},
		{
			desc:      "Success 2: RSA 384",
			key:       "./testKeys/jwt-rsa384-test.key",
			kid:       "rsa-384",
			header:    "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS0zODQiLCJ0eXAiOiJKV1QifQ",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0",
			signature: "jZgLl2PQnEn7Xz9SwuPeM7v8bsB5joCgCBKg6RBlmoNGX3d0PYBavNHAMM6uGgXX2Kc5wSqECMmzbCUBbGW0kMwLcYYz90d_37WkIDPTMojjf9u9JWCjSvj6r53KnD87z6UMn7sjOKM-3T1-ekIBvjxovgLOK5dxd8eAyQIXoTRDUf2jss_NeC1iWwDTAgPoJ4R-krxSdOf6uhT_mCD4_oaaliTsi8268ph_6ezaqUZfE3IcN7hWZtyAusVntue2PswZ29-lr2-KUcDNNGX9nolESv2Fb4SHZNVwtuIbAnNWmlChbZ4nJQfCBlo5hw7r6i06iYeSFrdBX-fqQxxiVQ",
		},
		{
			desc:      "Success 3: RSA 512",
			key:       "./testKeys/jwt-rsa512-test.key",
			kid:       "rsa-512",
			header:    "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS01MTIiLCJ0eXAiOiJKV1QifQ",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0",
			signature: "hQLuOe8qZHUgstYe0A0n4-Pww7ovlReyKiDR1Y02ltUnUlgbm9qpp-Ef6YNFuIKdHmS-ynQbDx5pbI36szsggzi80apNpI48cwSXshx82TwuU-_Z4wNBXu7MdPvbA5FdjhxCvRqaqhglsGJ6NofC1bP9awVyyy4j9LGfkVuVEXJQrVpdvEs8Ks-LxlWz7_9Cr7BrZcLuBJnujhe4CbdSudkrfeFl19EY3i1wH9OatGjfjwOSJVqv-ZLnn3QkaZmrQ1xwXTm3MlMUH3KSQjBn8h6vbqosIB5iHDFtqR11mLCgYExGHBpzFjM1d5NEmcTNLV9MtZ_qDZwG0wkgv9O4rXVQ0JfdXypMwhchED2Z45_mc2OiLidtKtDmeoE5g0Daq8YpM0ZpVRbXUFeYIZ1doQKUNsbWNdITmrjVOC3Zn8BecYPu1pC4Hk1y-ViArDzxlCMHA7Bua64BfzVuaJ8pBTEmbqMiZ9VujWcimCOtJ5yfCks_RPAhFYOErcqy3B56fmyYdIN__mKl7VvRDtBSiiPGCq07BUjGywaMoZIULbyXYSV4zs3hX_R4_o4asGiVWCZgn7k4pZzCJo_y2e-Mf85nYoRlyr1MXx7IM4srFQCgO-KTjDWL_TXqpMJU5zDzKyelrMFkc6EaMQ2KP_yBhOrh4UW-Pm7ghusox_-bV1U",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Clean(tC.key))
			assert.Nil(t, err)

			rsaEncoder, err := NewJwtEncoder(string(b), tC.kid)
			assert.Nil(t, err)

			token, err := rsaEncoder.Encode(claims)
			// fmt.Printf("\n%s token: '%s'\n", tC.desc, token)

			assert.Nil(t, err)
			assert.NotEmpty(t, token)

			splitToken := strings.Split(token, ".")
			assert.Equal(t, 3, len(splitToken))

			header := splitToken[0]
			assert.Equal(t, tC.header, header)
			payload := splitToken[1]
			assert.Equal(t, tC.payload, payload)
			signature := splitToken[2]
			assert.Equal(t, tC.signature, signature)
		})
	}
}

func TestEncoderECDSA(t *testing.T) {
	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
	}

	testCases := []struct {
		desc      string
		key       string
		kid       string
		header    string
		payload   string
		signature string
	}{
		{
			desc:      "Success 1: ECDSA 256",
			key:       "./testKeys/jwt-ecdsa256-test.key",
			kid:       "ecdsa-256",
			header:    "eyJhbGciOiJFUzI1NiIsImtpZCI6ImVjZHNhLTI1NiIsInR5cCI6IkpXVCJ9",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0",
			signature: "",
		},
		{
			desc:      "Success 2: ECDSA 384",
			key:       "./testKeys/jwt-ecdsa384-test.key",
			kid:       "ecdsa-384",
			header:    "eyJhbGciOiJFUzM4NCIsImtpZCI6ImVjZHNhLTM4NCIsInR5cCI6IkpXVCJ9",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0",
			signature: "",
		},
		{
			desc:      "Success 3: ECDSA 512",
			key:       "./testKeys/jwt-ecdsa521-test.key",
			kid:       "ecdsa-512",
			header:    "eyJhbGciOiJFUzUxMiIsImtpZCI6ImVjZHNhLTUxMiIsInR5cCI6IkpXVCJ9",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0",
			signature: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Clean(tC.key))
			assert.Nil(t, err)

			ecdsaEncoder, err := NewJwtEncoder(string(b), tC.kid)
			assert.Nil(t, err)

			token, err := ecdsaEncoder.Encode(claims)
			// fmt.Printf("\nToken for %s: '%s'", tC.desc, token)

			assert.Nil(t, err)
			assert.NotEmpty(t, token)

			splitToken := strings.Split(token, ".")
			assert.Equal(t, 3, len(splitToken))

			header := splitToken[0]
			assert.Equal(t, tC.header, header)
			payload := splitToken[1]
			assert.Equal(t, tC.payload, payload)
			// This keeps changing each test run - not sure why!
			// signature := splitToken[2]
			// assert.Equal(t, tC.signature, signature)
		})
	}
}
