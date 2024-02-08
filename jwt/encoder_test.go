package jwt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	encoderAuthKey string = "./testKeys/jwt-rsa256-test-webgateway.key"
)

func TestNewEncoder(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthPrivateKey))
	require.NoError(t, err)
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
			kid:            "web-gateway",
			expectedErrMsg: "",
		},
		{
			desc:           "Error 1: missing key",
			key:            "",
			kid:            "",
			expectedErrMsg: "invalid key",
		},
		{
			desc:           "Error 2: bad key",
			key:            "bad key",
			kid:            "bad",
			expectedErrMsg: "invalid key",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			encoder, err := NewJwtEncoder(tC.key, tC.kid)
			if tC.expectedErrMsg != "" {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tC.expectedErrMsg)
				assert.NotNil(t, encoder)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, encoder)
			}
		})
	}
}

func TestEncoderEncodeStandardClaims(t *testing.T) {
	privateKeyBytes, err := os.ReadFile(filepath.Clean(testAuthPrivateKey))
	require.NoError(t, err)

	encoder, err := NewJwtEncoder(string(privateKeyBytes), WebGatewayKid)
	assert.Nil(t, err)

	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		ExpiresAt:       time.Unix(2208952861, 0), //  1/1/2040
		IssuedAt:        time.Unix(1577800861, 0), // 1/1/2020
		NotBefore:       time.Unix(1577800861, 0), // 1/1/2020
	}

	token, err := encoder.Encode(claims)
	// fmt.Printf("Token: '%s'", token)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	splitToken := strings.Split(token, ".")
	assert.Equal(t, 3, len(splitToken))

	header := splitToken[0]
	assert.Equal(t, "eyJhbGciOiJSUzI1NiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0", header)
	payload := splitToken[1]
	assert.Equal(t, "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMDg5NTI4NjEsIm5iZiI6MTU3NzgwMDg2MSwiaWF0IjoxNTc3ODAwODYxfQ", payload)
	signature := splitToken[2]
	assert.Equal(t, "wV6z_kUjsKUebT7wcyTLqelUGhGfNA_82I_3lSsZOutegU4ct_652tp8enACcpgv2ZRmhbmCIC5w7PQn4rGLPx4Bdffjvf_4HvCqtvig0JV-0lmCpbNaadK93kiYNteZYFjLokLRKEHDt-uOoQbiWhc8DQBn7KbebLBRqp28HF28WL4-WVPDFsQ6H6pWL1RsXiuGyY8pI1y5b02t8-mte7CzrVx6uBgHPvfGgzWiTw4WpauxOxXUWTBIfK34OmPLb5sJQjrhM9RysE76j9703ptVfygTpokCcit-v_K3XlzQWw0T9sVfOu35mOS-NtXPJLB4-PK__gR60nANB-nNMsFf2Z1_ok44GAKE3an5Bi7cEaM-S5ZSbkq0rm6gbxEVZT5yjmJMNeecSgKc1dt_TAK5VMg5SJKzxXJ1DhLvKIB3rVNLyNfZVJT3mQW5NIBMjmfZad69_cu2TtS2_b8jOso7C7Vc3V-rB7MWYLsS47SDA46HFcAJvq7vsUHWM7POhoZEdSyN0cnpw0pWEOnhtpguJIw5XtrLQX00h2FytTAZEBnBUBYU-4AMQUBjK_FeEv5zDXSiXtR-iMs1YM0Qryhcw3Cx2EInkF8qt5AGNp6mYMwsFLJiv0RMO0CxxZ-08uZsng3Ekv8ewquKsXR5ZjSvtHO0vYSQEUcYvRUjELE", signature)
}
