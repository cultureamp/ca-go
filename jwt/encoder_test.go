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
	b, err := os.ReadFile(filepath.Clean(testRSAPrivateKey))
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

func TestEncoderRSAEncodeStandardClaims(t *testing.T) {
	// useful to create RS256 test tokens https://jwt.io/

	b1, err := os.ReadFile(filepath.Clean(testRSAPrivateKey))
	require.NoError(t, err)

	rsaEncoder, err := NewJwtEncoder(string(b1), webGatewayKid)
	assert.Nil(t, err)

	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	token, err := rsaEncoder.Encode(claims)
	// fmt.Printf("Token: '%s'", token)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	splitToken := strings.Split(token, ".")
	assert.Equal(t, 3, len(splitToken))

	header := splitToken[0]
	assert.Equal(t, "eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0", header)
	payload := splitToken[1]
	assert.Equal(t, "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ", payload)
	signature := splitToken[2]
	assert.Equal(t, "B4yOj6cwkICFgIlqOCxV2nIrGS_u8O2zk22uqJW40dpmm0TD3rH57Fjq_TwNSIpx84tIfRUhA-FHfHu-ci0epurvJBcQ_nOG1IfRlxOjd1goZjxPPplddwelECQGCdAyqkoGHXy8YgTe0ZvupPijfRIVmgpJcznmQphqLIuIJhcFGnoruhp4NAxQfqyONQf1S5h2H57-vvmXnQk5tpdocXYC-MP3jFtmNukmdUWpsFlpr2Fclgy3d4opf2fDQzdC51vBpVl1DjKEngjGULtRo4jDy7VRKvrdHhNX25zeUQSsKyetlWARnn-O2RT_d7kYAbBBy195kqtplZ47QQjhptW8WBEfS8X0-wjOHM04gdW3p1iAJ4A88wYywy1T75zUMTH2iPiIHRilzwwPj5j4tWPiUCj__i8tQvLXIZVIIpV7jdP1yP9Kp_Vb2WV-DKy9osiImZotc_kAWxl5Jq6xqhKNAnRirWrwk1q_Z7KmPmnswC84Ao6h3Lqf728pR5NVQzFB2t5vWvFk-ocAx0gKNCGF0fug4PUS5t_M5WecFkLOrAx68fvRLfr7BA1JFAP6wPu4Alz0HbtixD1gUC6bHO4A8g7pb0lWoLE0a4hKkPnvrQjtV5ccjpVIj-4sgQLr9zIpYnPwxbzGg13DRGBPySKic7qrD4nBEktcep01q50", signature)
}

func TestEncoderECDSAEncodeStandardClaims(t *testing.T) {
	// useful to create RS256 test tokens https://jwt.io/

	b1, err := os.ReadFile(filepath.Clean(testECDSAPrivateKey))
	require.NoError(t, err)

	ecdsaEncoder, err := NewJwtEncoder(string(b1), "ecdsa-test")
	assert.Nil(t, err)

	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	token, err := ecdsaEncoder.Encode(claims)
	// fmt.Printf("Token: '%s'", token)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	splitToken := strings.Split(token, ".")
	assert.Equal(t, 3, len(splitToken))

	header := splitToken[0]
	assert.Equal(t, "eyJhbGciOiJFUzUxMiIsImtpZCI6ImVjZHNhLXRlc3QiLCJ0eXAiOiJKV1QifQ", header)
	payload := splitToken[1]
	assert.Equal(t, "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ", payload)
	// This part of the signature changes each test run - not sure why!
	// signature := splitToken[2]
	// assert.Equal(t, "AHdNCgv1hBIa27cyDZtbal1Hv5dMP54YfBh5ZWXPlp2u9k_YslV5bibovguWmbDD640oN4OrMceqpkWTC8BD1IqAAeXqgoQSyAZB3AjopVKV66JK-_6TXeDVmckzA1x_2H-VFs1-6UXch7dWYdmS1vgxq0SXLa8Tf0sqeSyLJMy2l0z", signature)
}
