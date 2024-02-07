package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	day  = 24 * time.Hour
	year = 365 * day // approx

	decoderAuthJwks string = "./testKeys/development.jwks"
)

func TestNewDecoder(t *testing.T) {
	pubJwkKeyBytes, err := os.ReadFile(filepath.Clean(decoderAuthJwks))
	require.NoError(t, err)
	validJwks := string(pubJwkKeyBytes)

	testCases := []struct {
		desc           string
		jwks           string
		expectedErrMsg string
	}{
		{
			desc:           "Success 1: valid jwks",
			jwks:           validJwks,
			expectedErrMsg: "",
		},
		{
			desc:           "Error 1: missing key",
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
				assert.NotNil(t, decoder)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, decoder)
			}
		})
	}
}

func TestDecoderDecodeAllClaims(t *testing.T) {
	pubJwksBytes, err := os.ReadFile(filepath.Clean(decoderAuthJwks))
	require.NoError(t, err)

	decoder, err := NewJwtDecoder(string(pubJwksBytes))
	assert.Nil(t, err)
	assert.NotNil(t, decoder)

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMDg5NTI4NjEsIm5iZiI6MTU3NzgwMDg2MSwiaWF0IjoxNTc3ODAwODYxfQ.wV6z_kUjsKUebT7wcyTLqelUGhGfNA_82I_3lSsZOutegU4ct_652tp8enACcpgv2ZRmhbmCIC5w7PQn4rGLPx4Bdffjvf_4HvCqtvig0JV-0lmCpbNaadK93kiYNteZYFjLokLRKEHDt-uOoQbiWhc8DQBn7KbebLBRqp28HF28WL4-WVPDFsQ6H6pWL1RsXiuGyY8pI1y5b02t8-mte7CzrVx6uBgHPvfGgzWiTw4WpauxOxXUWTBIfK34OmPLb5sJQjrhM9RysE76j9703ptVfygTpokCcit-v_K3XlzQWw0T9sVfOu35mOS-NtXPJLB4-PK__gR60nANB-nNMsFf2Z1_ok44GAKE3an5Bi7cEaM-S5ZSbkq0rm6gbxEVZT5yjmJMNeecSgKc1dt_TAK5VMg5SJKzxXJ1DhLvKIB3rVNLyNfZVJT3mQW5NIBMjmfZad69_cu2TtS2_b8jOso7C7Vc3V-rB7MWYLsS47SDA46HFcAJvq7vsUHWM7POhoZEdSyN0cnpw0pWEOnhtpguJIw5XtrLQX00h2FytTAZEBnBUBYU-4AMQUBjK_FeEv5zDXSiXtR-iMs1YM0Qryhcw3Cx2EInkF8qt5AGNp6mYMwsFLJiv0RMO0CxxZ-08uZsng3Ekv8ewquKsXR5ZjSvtHO0vYSQEUcYvRUjELE"
	claim, err := decoder.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, "abc123", claim.AccountId)
	assert.Equal(t, "xyz234", claim.RealUserId)
	assert.Equal(t, "xyz345", claim.EffectiveUserId)
	assert.Equal(t, 2040, claim.ExpiresAt.Year())
}
