package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// useful to create RS256 test tokens https://jwt.io/
// useful for PEM to JWKS https://jwkset.com/generate

func TestNewDecoder(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	assert.Nil(t, err)
	validJwks := string(b)

	testCases := []struct {
		desc           string
		jwks           string
		expectedErrMsg string
	}{
		{
			desc:           "Success 1: valid jwks and found kid",
			jwks:           validJwks,
			expectedErrMsg: "",
		},
		{
			desc:           "Error 1: missing jwk keys",
			jwks:           "",
			expectedErrMsg: "missing jwks",
		},
		{
			desc:           "Error 2: JWKS json bad",
			jwks:           "{\"bad\": \"jwks-json\" }",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
		{
			desc:           "Error 3: JWKS json invalid",
			jwks:           "invalid JSON",
			expectedErrMsg: "failed to unmarshal JWK set",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoder(func() string { return tC.jwks })
			if tC.expectedErrMsg != "" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedErrMsg)
				assert.Nil(t, decoder)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, decoder)
			}
		})
	}
}

func TestDecoderDecodeAllClaims(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	assert.Nil(t, err)
	testJWKS := string(b)

	jwks := func() string {
		return testJWKS
	}

	testCases := []struct {
		desc            string
		token           string
		expectedErrMsg  string
		accountId       string
		realUserId      string
		effectiveUserId string
		issuer          string
		subject         string
		audience        []string
		year            int
	}{
		{
			desc:            "Success 1: valid token RSA",
			token:           "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS01MTIiLCJ0eXAiOiJKV1QifQ.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0.hQLuOe8qZHUgstYe0A0n4-Pww7ovlReyKiDR1Y02ltUnUlgbm9qpp-Ef6YNFuIKdHmS-ynQbDx5pbI36szsggzi80apNpI48cwSXshx82TwuU-_Z4wNBXu7MdPvbA5FdjhxCvRqaqhglsGJ6NofC1bP9awVyyy4j9LGfkVuVEXJQrVpdvEs8Ks-LxlWz7_9Cr7BrZcLuBJnujhe4CbdSudkrfeFl19EY3i1wH9OatGjfjwOSJVqv-ZLnn3QkaZmrQ1xwXTm3MlMUH3KSQjBn8h6vbqosIB5iHDFtqR11mLCgYExGHBpzFjM1d5NEmcTNLV9MtZ_qDZwG0wkgv9O4rXVQ0JfdXypMwhchED2Z45_mc2OiLidtKtDmeoE5g0Daq8YpM0ZpVRbXUFeYIZ1doQKUNsbWNdITmrjVOC3Zn8BecYPu1pC4Hk1y-ViArDzxlCMHA7Bua64BfzVuaJ8pBTEmbqMiZ9VujWcimCOtJ5yfCks_RPAhFYOErcqy3B56fmyYdIN__mKl7VvRDtBSiiPGCq07BUjGywaMoZIULbyXYSV4zs3hX_R4_o4asGiVWCZgn7k4pZzCJo_y2e-Mf85nYoRlyr1MXx7IM4srFQCgO-KTjDWL_TXqpMJU5zDzKyelrMFkc6EaMQ2KP_yBhOrh4UW-Pm7ghusox_-bV1U",
			expectedErrMsg:  "",
			accountId:       "abc123",
			realUserId:      "xyz234",
			effectiveUserId: "xyz345",
			issuer:          "encoder-name",
			subject:         "test",
			audience:        []string{"decoder-name"},
			year:            2040,
		},
		{
			desc:            "Success 2: valid token ECDSA",
			token:           "eyJhbGciOiJFUzI1NiIsImtpZCI6ImVjZHNhLTI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0.tuevuDAPABxVW8_SqD6T8SeCrFucXxyBEIk3Yk8xKsfqYSNv1nW0HPNPE_a-E02fKYAsXJY8yn1R8u1idj9Z5Q",
			expectedErrMsg:  "",
			accountId:       "abc123",
			realUserId:      "xyz234",
			effectiveUserId: "xyz345",
			issuer:          "encoder-name",
			subject:         "test",
			audience:        []string{"decoder-name"},
			year:            2040,
		},
		{
			desc:            "Error 1: invalid token",
			token:           "eyJhbGciOiJSUzI1NiIsImtpZCI6ld2F50.eyJhY2g2MSwiaWF0IjoxNTc3ODAwODYxfQ.wV6z_kUjsKUebT7RUjELE",
			expectedErrMsg:  "token is malformed",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			issuer:          "encoder-name",
			subject:         "test",
			audience:        []string{"decoder-name"},
			year:            1,
		},
		{
			desc:            "Error 2: missing kid",
			token:           "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.S48Yfs0COCQR70jCrD2y6kD26nns6-c9CLvxKTahxzv0KrkdAXC7I62yIz2yD6j3v3rTKVQ8eGhSKOkN6EU_M8BZa5ltt7TmcIOnn5RWbwnfSfFLMR3njzlMiRT2MGAi2A2WMkx_LrTk9PZZRIlfceQxFVhThjc-Dp92C_zFJARZ8yss3upAW0m0pbeD5Y23GWs6bkkBbAAvh8Rw6rICjW6qROnqj6u8mcb4bS3kDlmmFkYnQdKMLu4bWa6twyLwUMg0N-Y5h2rp6GyAYuTrqyKif5IU1IEhNW63gj5h1xCLNyX4ZGJsNSZP_HOQGVVQMBDg2rsg7tBxow_E2wvYTgDYn8f1SjKE7vKdL2uYzA732hcd63fNwJpNcrwFs3lW8DjM_VYf-M4ePSr4GqHg6PawTCFgWCVNvi-lsmogfRUq_1t21GXlX7pQd029CFJ7mnnxUBau7KxTuX-Pxpny3jhYpJ87GlDA3WaB0r1tEg4Hl87VDawQ5Cb5ac6R6eXEO34i5oESVt6lFL-wpWUnU7KbiWVxKkSifN0M27IE8vAEsUBhgD8sKZpBsUvUDpRcb_atpFE_xU0K6DZXGUgFpZBx-CULmmfoDubTNfRNtqwmJiXI-M1YyiRbc_lOVQBAibuZ20ucixyhhqYSa-5fWa4m5NcjkRquTR2J-OaxmhA",
			expectedErrMsg:  "missing key_id (kid) header",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			issuer:          "encoder-name",
			subject:         "test",
			audience:        []string{"decoder-name"},
			year:            1,
		},
		{
			desc:            "Error 3: signed with EdDSA",
			token:           "eyJhbGciOiJFZERTQSIsImtpZCI6ImVkZHNhIiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.ZjBuTNqC74M525ROjM1kRANBCb7JUjl6ko8dSD52S-Q_f5p7EaUM5TxCPg4rICvVWxF26B99EkVzNNhivcw9Dw",
			expectedErrMsg:  "signing method EdDSA is invalid",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			issuer:          "encoder-name",
			subject:         "test",
			audience:        []string{"decoder-name"},
			year:            1,
		},
		{
			desc:            "Error 4: bad jwks key",
			token:           "eyJhbGciOiJFUzUxMiIsImtpZCI6ImJhZC1lY2RzYSIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.ALtrklCsBs5Z3QZp8p8qfClXKS8pj0asOzHspDRghmOoA5XWzYriMWFElb2UlhZeSazIsyR7P2k_oKt8S4qdY3jxAbpkBL6BQURudNeuw-bBipPPBehscuUlOhQl7ckl4RO7c5U60uHyaek0m4LMpEIuziWX9IHikDSVBzkNuji8zWQ1",
			expectedErrMsg:  "bad public key in jwks",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			issuer:          "encoder-name",
			subject:         "test",
			audience:        []string{"decoder-name"},
			year:            1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoder(jwks)
			assert.Nil(t, err)
			assert.NotNil(t, decoder)

			claim, err := decoder.Decode(tC.token)
			if tC.expectedErrMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, tC.accountId, claim.AccountId)
				assert.Equal(t, tC.realUserId, claim.RealUserId)
				assert.Equal(t, tC.effectiveUserId, claim.EffectiveUserId)
				assert.Equal(t, tC.issuer, claim.Issuer)
				assert.Equal(t, tC.subject, claim.Subject)
				assert.Equal(t, tC.audience, claim.Audience)
				assert.Equal(t, tC.year, claim.ExpiresAt.Year())
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedErrMsg)
			}
		})
	}
}

func TestDecoderRotateKeys(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	assert.Nil(t, err)
	testJWKS := string(b)

	attempt := 0
	jwks := func() string {
		attempt++
		if attempt == 1 {
			// first time return just 1 key
			return `{"keys" : [{
      "alg":"RS256",
      "kty":"RSA",
      "e":"AQAB",
      "kid":"web-gateway",
      "n":"zkzpPa8QB5JwYWJI1W3WmxnMwusvFZb-0EVY4Sko3C1zwBcY8P6NucHo1epXTO-rFQy8JPiSMyTBINkmDP0d1jfvJF_RDL8Gzi1_aM2mScsPxmXA7ftqHdvcaqP0aobuYNJSEk_3erM6iddBJwsKY5BNkzS-R9szsfCgnDdfN-9JvChpfrTvoOwI-vtsqpkgIgGB4uCeQ0CPvqZzsRMJyWouEt0Jj7huKXBOvDBuoZdInuh-2kzNpm9KEkdbB0wzhC57MnyA3ap0I-ES374utQGM1EbZfW68T0QU3t--Q7L7yQ4D8WjRLZw_WTS8amcLRYf0urb3yTmvQFA4ryhc25dBUF68xPrC2kETljf6SLtig2bWvr-TGqGiyLnqiPloSxeBtpZhWSBgH8KJ7iHjwCyT2dSMEhf-ouivT2rEn5wEP3joDPywBqywKs-hbJrOB_x9cg4dGqERuljvW02tMGHu1JTK8tb23wWl8_5RSPHGetM526G3MW8r8hJ4mPHASPzQ2jWM_XhHtvLOg4_0V3CczMe93e6ilWkxala1hnZA180lOFoOOscdQmcH7LbOjkH6Iwb_9lc0Ez6n2tcfuY9p1aujcsJ5uQNBJtoX4kOSTM7LfUJa88ZbUkOeJ9AHhCe9xqaAS-W0LJYR00-JZcsaZz31F2DSFMmOWLUCVZ8",
      "use":"sig"
    }]}	`
		}

		// for subsequent calls return all the keys to simulate a new key being added
		return testJWKS
	}

	decoder, err := NewJwtDecoder(jwks,
		WithDecoderJwksExpiry(time.Millisecond*500),
		WithDecoderRotateWindow(time.Millisecond*100),
	)
	assert.Nil(t, err)
	assert.NotNil(t, decoder)

	// try and decode is straight away - should fail as the kid isn't in the jwks
	token := "eyJhbGciOiJFUzI1NiIsImtpZCI6ImVjZHNhLTI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0.tuevuDAPABxVW8_SqD6T8SeCrFucXxyBEIk3Yk8xKsfqYSNv1nW0HPNPE_a-E02fKYAsXJY8yn1R8u1idj9Z5Q"
	claim, err := decoder.Decode(token)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "no matching key_id (kid) header")

	// sleep so now the JWKS is within the rotation window
	time.Sleep(110 * time.Millisecond)
	claim, err = decoder.Decode(token)
	assert.Equal(t, 2, attempt)
	assert.Nil(t, err)

	require.NotNil(t, claim)
	assert.Equal(t, "abc123", claim.AccountId)
	assert.Equal(t, "xyz234", claim.RealUserId)
	assert.Equal(t, "xyz345", claim.EffectiveUserId)
	assert.Equal(t, "encoder-name", claim.Issuer)
	assert.Equal(t, "test", claim.Subject)
	assert.Equal(t, []string{"decoder-name"}, claim.Audience)
	assert.Equal(t, 2040, claim.ExpiresAt.Year())
}
