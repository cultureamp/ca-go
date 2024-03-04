package jwt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// useful to create RS256 test tokens https://jwt.io/
// useful for PEM to JWKS https://jwkset.com/generate

func TestNewDecoder(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	require.NoError(t, err)
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
			decoder, err := NewJwtDecoder(tC.jwks)
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
	require.NoError(t, err)

	testCases := []struct {
		desc            string
		token           string
		expectedErrMsg  string
		accountId       string
		realUserId      string
		effectiveUserId string
		year            int
	}{
		{
			desc:            "Success 1: valid token RSA",
			token:           "eyJhbGciOiJSUzI1NiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.IYPu_PGUO7lpy_wTSObA4S-n9HQUwPf2kTG2AzvSFUwfz994SHZOazYL7CyiRqqhIndIt5R4CQ3cXY7_Lok_wgBQ-U4FAciJw0Fx9tawJIEqwVeL10P4w0h5OIU21E7jeNmlcLOO57QN-ip7hc_--zyAFVKV5qjlbemuHWWpeUGu62SsdHr4J33O6hR8ubTyfXVF7wxKhNM4hCdM7PNanP9OOyAgEWxhwutURiA1nJsATwDf6QKNceGpqkb5A31PvFdfPHoktY4u6e4feBt2KjYJ1xy9opDlllFOEIwTw4nuksQk4q3437bGtfoQkC_CTGO83YTX5GHs70rxu_AubBxCazqSxqMwagiekkpgKZd6d0g7u5F5K8QImRJsore3oHNDAuVg7pbZmH9sApFN_bJhonOkECoPeeF5oYLSLHOXjN7CakvAsmCW01_ENPVXXO2E1yObzwmsY28_Ox5r_jC6XugGdXVfco6l1Oqbxb0ogG6BbOngYEZwVMbEO5qsBnUtBfr0nNUjFKIYCYXdpoeT_bxlt8GI4H2cMAb6FGa_XIEd60fJGazgAk9axA61xHEnqxgUyZv5PEL908zPBRvcNGpQuMsDpGOXTOQ_fgJO1IRBx4VwWcobzKbOyRNarTNwQZH0OY13HMMnFoiPjk8U0fWkJdj1ujobTQYYtz0",
			expectedErrMsg:  "",
			accountId:       "abc123",
			realUserId:      "xyz234",
			effectiveUserId: "xyz345",
			year:            2040,
		},
		{
			desc:            "Success 2: valid token ECDSA",
			token:           "eyJhbGciOiJFUzUxMiIsImtpZCI6ImVjZHNhLXRlc3QiLCJ0eXAiOiJKV1QifQ.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.AQvawAQSDqhtgsF6zwdULe1b9csrSOzp-2zgjjBLpweex3v-KYMSP6rc65aeGiSTqVNhifrmLoeF1lcb-OFh9hASATehAAEYEZEVnFycDqcxRjdYTNwY048RzhiY2zkK61uyyLu8HOtEvXj827NHjvPBbNjl9uStwQZlDRouwqyS_Elg",
			expectedErrMsg:  "",
			accountId:       "abc123",
			realUserId:      "xyz234",
			effectiveUserId: "xyz345",
			year:            2040,
		},
		{
			desc:            "Error 1: invalid token",
			token:           "eyJhbGciOiJSUzI1NiIsImtpZCI6ld2F50.eyJhY2g2MSwiaWF0IjoxNTc3ODAwODYxfQ.wV6z_kUjsKUebT7RUjELE",
			expectedErrMsg:  "token is malformed",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
		{
			desc:            "Error 2: missing kid",
			token:           "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.S48Yfs0COCQR70jCrD2y6kD26nns6-c9CLvxKTahxzv0KrkdAXC7I62yIz2yD6j3v3rTKVQ8eGhSKOkN6EU_M8BZa5ltt7TmcIOnn5RWbwnfSfFLMR3njzlMiRT2MGAi2A2WMkx_LrTk9PZZRIlfceQxFVhThjc-Dp92C_zFJARZ8yss3upAW0m0pbeD5Y23GWs6bkkBbAAvh8Rw6rICjW6qROnqj6u8mcb4bS3kDlmmFkYnQdKMLu4bWa6twyLwUMg0N-Y5h2rp6GyAYuTrqyKif5IU1IEhNW63gj5h1xCLNyX4ZGJsNSZP_HOQGVVQMBDg2rsg7tBxow_E2wvYTgDYn8f1SjKE7vKdL2uYzA732hcd63fNwJpNcrwFs3lW8DjM_VYf-M4ePSr4GqHg6PawTCFgWCVNvi-lsmogfRUq_1t21GXlX7pQd029CFJ7mnnxUBau7KxTuX-Pxpny3jhYpJ87GlDA3WaB0r1tEg4Hl87VDawQ5Cb5ac6R6eXEO34i5oESVt6lFL-wpWUnU7KbiWVxKkSifN0M27IE8vAEsUBhgD8sKZpBsUvUDpRcb_atpFE_xU0K6DZXGUgFpZBx-CULmmfoDubTNfRNtqwmJiXI-M1YyiRbc_lOVQBAibuZ20ucixyhhqYSa-5fWa4m5NcjkRquTR2J-OaxmhA",
			expectedErrMsg:  "missing key_id (kid) header",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
		{
			desc:            "Error 3: signed with EdDSA",
			token:           "eyJhbGciOiJFZERTQSIsImtpZCI6ImVkZHNhIiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.ZjBuTNqC74M525ROjM1kRANBCb7JUjl6ko8dSD52S-Q_f5p7EaUM5TxCPg4rICvVWxF26B99EkVzNNhivcw9Dw",
			expectedErrMsg:  "unexpected signing method",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
		{
			desc:            "Error 4: bad jwks key",
			token:           "eyJhbGciOiJFUzUxMiIsImtpZCI6ImJhZC1lY2RzYSIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.ALtrklCsBs5Z3QZp8p8qfClXKS8pj0asOzHspDRghmOoA5XWzYriMWFElb2UlhZeSazIsyR7P2k_oKt8S4qdY3jxAbpkBL6BQURudNeuw-bBipPPBehscuUlOhQl7ckl4RO7c5U60uHyaek0m4LMpEIuziWX9IHikDSVBzkNuji8zWQ1",
			expectedErrMsg:  "bad public key in jwks",
			accountId:       "",
			realUserId:      "",
			effectiveUserId: "",
			year:            1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			decoder, err := NewJwtDecoder(string(b))
			assert.Nil(t, err)
			assert.NotNil(t, decoder)

			claim, err := decoder.Decode(tC.token)

			if tC.expectedErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedErrMsg)
			}

			assert.Equal(t, tC.accountId, claim.AccountId)
			assert.Equal(t, tC.realUserId, claim.RealUserId)
			assert.Equal(t, tC.effectiveUserId, claim.EffectiveUserId)
			assert.Equal(t, tC.year, claim.ExpiresAt.Year())
		})
	}
}
