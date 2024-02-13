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

func TestNewEncoder(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testDefaultAuthPrivateKey))
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
				assert.Nil(t, encoder)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, encoder)
			}
		})
	}
}

func TestEncoderEncodeStandardClaims(t *testing.T) {
	// useful to create RS256 test tokens https://jwt.io/

	privateKeyBytes, err := os.ReadFile(filepath.Clean(testDefaultAuthPrivateKey))
	require.NoError(t, err)

	encoder, err := NewJwtEncoder(string(privateKeyBytes), webGatewayKid)
	assert.Nil(t, err)

	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
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
	assert.Equal(t, "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ", payload)
	signature := splitToken[2]
	assert.Equal(t, "IYPu_PGUO7lpy_wTSObA4S-n9HQUwPf2kTG2AzvSFUwfz994SHZOazYL7CyiRqqhIndIt5R4CQ3cXY7_Lok_wgBQ-U4FAciJw0Fx9tawJIEqwVeL10P4w0h5OIU21E7jeNmlcLOO57QN-ip7hc_--zyAFVKV5qjlbemuHWWpeUGu62SsdHr4J33O6hR8ubTyfXVF7wxKhNM4hCdM7PNanP9OOyAgEWxhwutURiA1nJsATwDf6QKNceGpqkb5A31PvFdfPHoktY4u6e4feBt2KjYJ1xy9opDlllFOEIwTw4nuksQk4q3437bGtfoQkC_CTGO83YTX5GHs70rxu_AubBxCazqSxqMwagiekkpgKZd6d0g7u5F5K8QImRJsore3oHNDAuVg7pbZmH9sApFN_bJhonOkECoPeeF5oYLSLHOXjN7CakvAsmCW01_ENPVXXO2E1yObzwmsY28_Ox5r_jC6XugGdXVfco6l1Oqbxb0ogG6BbOngYEZwVMbEO5qsBnUtBfr0nNUjFKIYCYXdpoeT_bxlt8GI4H2cMAb6FGa_XIEd60fJGazgAk9axA61xHEnqxgUyZv5PEL908zPBRvcNGpQuMsDpGOXTOQ_fgJO1IRBx4VwWcobzKbOyRNarTNwQZH0OY13HMMnFoiPjk8U0fWkJdj1ujobTQYYtz0", signature)
}
