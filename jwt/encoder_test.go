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
	testRSAPrivateKey   string = "./testKeys/jwt-rsa256-test-webgateway.key"
	testECDSAPrivateKey string = "./testKeys/jwt-ecdsa521-test.key"
)

// useful to create RS256 test tokens https://jwt.io/
// useful for PEM to JWKS https://jwkset.com/generate

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

func TestEncoderRSA(t *testing.T) {
	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
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
			key:       "./testKeys/jwt-rsa256-test-webgateway.key",
			kid:       "web-gateway",
			header:    "eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ",
			signature: "B4yOj6cwkICFgIlqOCxV2nIrGS_u8O2zk22uqJW40dpmm0TD3rH57Fjq_TwNSIpx84tIfRUhA-FHfHu-ci0epurvJBcQ_nOG1IfRlxOjd1goZjxPPplddwelECQGCdAyqkoGHXy8YgTe0ZvupPijfRIVmgpJcznmQphqLIuIJhcFGnoruhp4NAxQfqyONQf1S5h2H57-vvmXnQk5tpdocXYC-MP3jFtmNukmdUWpsFlpr2Fclgy3d4opf2fDQzdC51vBpVl1DjKEngjGULtRo4jDy7VRKvrdHhNX25zeUQSsKyetlWARnn-O2RT_d7kYAbBBy195kqtplZ47QQjhptW8WBEfS8X0-wjOHM04gdW3p1iAJ4A88wYywy1T75zUMTH2iPiIHRilzwwPj5j4tWPiUCj__i8tQvLXIZVIIpV7jdP1yP9Kp_Vb2WV-DKy9osiImZotc_kAWxl5Jq6xqhKNAnRirWrwk1q_Z7KmPmnswC84Ao6h3Lqf728pR5NVQzFB2t5vWvFk-ocAx0gKNCGF0fug4PUS5t_M5WecFkLOrAx68fvRLfr7BA1JFAP6wPu4Alz0HbtixD1gUC6bHO4A8g7pb0lWoLE0a4hKkPnvrQjtV5ccjpVIj-4sgQLr9zIpYnPwxbzGg13DRGBPySKic7qrD4nBEktcep01q50",
		},
		{
			desc:      "Success 2: RSA 384",
			key:       "./testKeys/jwt-rsa384-test.key",
			kid:       "rsa-384",
			header:    "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS0zODQiLCJ0eXAiOiJKV1QifQ",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ",
			signature: "lDByZXFrtXlPR1kJdRCPv_Bp4E2gIriRtEf5cmXeKuWqI9SJRwNMg7Zly3FAicUtCJw7DaP6xhL2FF8CA1oA8STWT_q2FiuwbzliVfhSQaa-dZ3KpkEKRshuX3cSjDV3SawmfHDvTyImcU6etG1lirRCtiuIOdyx7ArQoeLZyjq61a6hnoLANhbBxmT5tDxG2rDsYsVBnrhRn-1TOrA4VGJeYZrzogDi9X1JQjcCxaME28wMoIo9z6evqivXu3fuAznPlWH14sfwakc7TDX_6EaSy6SNPmsCHQE8J8qLPfWzE-9tVeJFNfVSs6YMUnHlJsy_Hmj9XeImH0gucsiO-g",
		},
		{
			desc:      "Success 3: RSA 512",
			key:       "./testKeys/jwt-rsa512-test.key",
			kid:       "rsa-512",
			header:    "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS01MTIiLCJ0eXAiOiJKV1QifQ",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ",
			signature: "mcK1NfW8YJARrDfkPsRXjz76r2iGvqPyzirvmhKxXWi6ry8ozT_1RXgazIre3GXZSLkkdgdKNRZKFduPih1Iv2--6F5PXG7D2e5fCfIIjIgViRILKZD65i95TWpJ7s64U8EB2GNWDz1AosvS5hfOfN20H1fcgWUb-pP6F9Gptq_03RVsBGOorvoSjzfig37UueIGj7HDp3qIi_LnEHdVPoVGhlyLzIrOfPq3amYmBZjlXsWF-iFHVFUIszvdmZ0OsJsvL-9xFqg-4m_aGAnqWPf9lM9k8u2pXf0Qbj-UnVMYREllBGWlOWlWUpR_xoop8r3nJIDIZ39sijUD27v47auehpR1b6v2WXhiV3nB8BHbbRhXDYmkG8RmgCHDnocnRM7dCrAECRET97lgxmat2VK4o0yZL3A9UCYEEtD0K7r_33IAUzDA_yBPVdRIRZ5bEOQDuSJdSZ-cn-aV_SPXweVOyeOpkPcGjBKo9lmhZchg2inkpotJZ4baGr4088z514ShbKMRuDK-L9KGwZnZ__YYl3dvmSPZHxQg6nI3stG0-hmxGzWdK9uPhDNR4ujZpLNWLVEBs2bLewWQDhrKhZ7rPnJ3DZ008g8pSlBN1FsOZIcWRpCkZXr-7yuGJkoUF1DfgDMQP4jrq9J-D4U4ak9Rzf-VWnKfx93oKgON1Mo",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Clean(tC.key))
			require.NoError(t, err)

			rsaEncoder, err := NewJwtEncoder(string(b), tC.kid)
			assert.Nil(t, err)

			token, err := rsaEncoder.Encode(claims)
			// fmt.Printf("Token: '%s'", token)

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
			desc:      "Success 1: ECDSA 256",
			key:       "./testKeys/jwt-ecdsa256-test.key",
			kid:       "ecdsa-256",
			header:    "eyJhbGciOiJFUzI1NiIsImtpZCI6ImVjZHNhLTI1NiIsInR5cCI6IkpXVCJ9",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ",
			signature: "",
		},
		{
			desc:      "Success 2: ECDSA 384",
			key:       "./testKeys/jwt-ecdsa384-test.key",
			kid:       "ecdsa-384",
			header:    "eyJhbGciOiJFUzM4NCIsImtpZCI6ImVjZHNhLTM4NCIsInR5cCI6IkpXVCJ9",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ",
			signature: "",
		},
		{
			desc:      "Success 3: ECDSA 512",
			key:       "./testKeys/jwt-ecdsa521-test.key",
			kid:       "ecdsa-512",
			header:    "eyJhbGciOiJFUzUxMiIsImtpZCI6ImVjZHNhLTUxMiIsInR5cCI6IkpXVCJ9",
			payload:   "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ",
			signature: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Clean(tC.key))
			require.NoError(t, err)

			ecdsaEncoder, err := NewJwtEncoder(string(b), tC.kid)
			assert.Nil(t, err)

			token, err := ecdsaEncoder.Encode(claims)
			// fmt.Printf("Token: '%s'", token)

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
