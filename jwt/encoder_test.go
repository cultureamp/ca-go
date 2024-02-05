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
	encoderAuthKey       string = "./testKeys/jwt.rs256.key.development.pem"
	encoderSecondAuthKey string = "./testKeys/jwt.rs256.key.development.extra_1.pem"
	encoderThirdAuthKey  string = "./testKeys/jwt.rs256.key.development.extra_2.pem"
)

func TestNewEncoderSuccess(t *testing.T) {
	privateKeyBytes, err := os.ReadFile(filepath.Clean(encoderAuthKey))
	require.NoError(t, err)

	encoder, err := NewJwtEncoder(string(privateKeyBytes), "")
	assert.Nil(t, err)
	assert.NotNil(t, encoder)
}

func TestEncoderEncodeStandardClaims(t *testing.T) {
	privateKeyBytes, err := os.ReadFile(filepath.Clean(encoderAuthKey))
	require.NoError(t, err)

	encoder, err := NewJwtEncoder(string(privateKeyBytes), "")
	assert.Nil(t, err)

	claims := &StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Expiry:          time.Unix(2208952861, 0), //  1/1/2040
	}

	token, err := encoder.Encode(claims)
	// fmt.Printf("Token: '%s'", token)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	splitToken := strings.Split(token, ".")
	assert.Equal(t, 3, len(splitToken))

	header := splitToken[0]
	assert.Equal(t, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9", header)

	// Note: Hard to test payload/signature because there is a elapsed time fields that always change...
	payload := splitToken[1]
	assert.Equal(t, "eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0Ii", payload[:94])
	// signature := splitToken[2]
	// assert.Equal(t, "PdNemOn7tIYNEvbY9FxVRY5p7tBRghMb4pBv59FVSO4PKCIpaNbg9IcxVzXyH1wvL1rD4j4jRnZJtT5-w4HAYC4YYpED9d_CqQVRPz4geaN-ZHzWjsAsp_B3XQCDpgMWmJUWepjkUVENGwE7Mg00lg4vfeeYOPFKt-goHnecpqin3XcgmbfaRi9CghS7611GM_CsCTjsFV_tCsR9f0l0sIf9dysGBSN0IT0zxG7Ro5UOagf98ZGPaZz8c1nrGC_1vO3zqMhSv_r8REqL13t8CulfxTZrDL4ARohzu9DqtVjG8sj3AXSjLGBPD10sGwvnwfn6lC7gJcBLOhJAezODF0AOgJll3Q7RgFHDb1-T0vphwegXaApsxW_Zm_LE2PpoQjPnmWpCkIlXTpZtHlmYy8xtjJQTpkw8_xDoD1-LsKZXxPGfzAbzE9kD7wBG8-hddT_mbk3nSPFV_t7M0WjCgmEAfhmS8QrCRUve3X8bdBByYZzFfe7Pmq9s2Ib_XQgaxH8r8Of9WKMb0JJ5ZH7X9w1utXgT4JkVhpQPsNoKSqprbNCcN5qqKfxr1tYA3AN2CQOJaSb6nYY-f1oxrYIgxiB-t2sP7nMKwn9Mtozy_jX5hU3s2fkRTFNC6slg7PYPMVuxVdMW72_fON35-UYf6dbNUnDwt3x99s6d0OioXik", signature)
}
