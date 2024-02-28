package jwt_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cultureamp/ca-go/jwt"
)

const (
	webGatewayKid = "web-gateway"
)

func Example() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	claim, err := jwt.Decode(token)
	fmt.Printf("The decoded token is '%s %s %s %s' (err='%+v')\n", claim.AccountId, claim.RealUserId, claim.EffectiveUserId, claim.ExpiresAt.UTC().Format(time.RFC3339), err)

	// To create a specific instance of the encoder and decoder you can use the following
	privateKeyBytes, err := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
	encoder, err := jwt.NewJwtEncoder(string(privateKeyBytes), webGatewayKid)

	token, err = encoder.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	b, err := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
	decoder, err := jwt.NewJwtDecoder(string(b))

	claim, err = decoder.Decode(token)
	fmt.Printf("The decoded token is '%s %s %s %s' (err='%+v')\n", claim.AccountId, claim.RealUserId, claim.EffectiveUserId, claim.ExpiresAt.UTC().Format(time.RFC3339), err)

	// Output:
	// The encoded token is 'eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.B4yOj6cwkICFgIlqOCxV2nIrGS_u8O2zk22uqJW40dpmm0TD3rH57Fjq_TwNSIpx84tIfRUhA-FHfHu-ci0epurvJBcQ_nOG1IfRlxOjd1goZjxPPplddwelECQGCdAyqkoGHXy8YgTe0ZvupPijfRIVmgpJcznmQphqLIuIJhcFGnoruhp4NAxQfqyONQf1S5h2H57-vvmXnQk5tpdocXYC-MP3jFtmNukmdUWpsFlpr2Fclgy3d4opf2fDQzdC51vBpVl1DjKEngjGULtRo4jDy7VRKvrdHhNX25zeUQSsKyetlWARnn-O2RT_d7kYAbBBy195kqtplZ47QQjhptW8WBEfS8X0-wjOHM04gdW3p1iAJ4A88wYywy1T75zUMTH2iPiIHRilzwwPj5j4tWPiUCj__i8tQvLXIZVIIpV7jdP1yP9Kp_Vb2WV-DKy9osiImZotc_kAWxl5Jq6xqhKNAnRirWrwk1q_Z7KmPmnswC84Ao6h3Lqf728pR5NVQzFB2t5vWvFk-ocAx0gKNCGF0fug4PUS5t_M5WecFkLOrAx68fvRLfr7BA1JFAP6wPu4Alz0HbtixD1gUC6bHO4A8g7pb0lWoLE0a4hKkPnvrQjtV5ccjpVIj-4sgQLr9zIpYnPwxbzGg13DRGBPySKic7qrD4nBEktcep01q50' (err='<nil>')
	// The decoded token is 'abc123 xyz234 xyz345 2040-02-02T12:12:12Z' (err='<nil>')
	// The encoded token is 'eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.B4yOj6cwkICFgIlqOCxV2nIrGS_u8O2zk22uqJW40dpmm0TD3rH57Fjq_TwNSIpx84tIfRUhA-FHfHu-ci0epurvJBcQ_nOG1IfRlxOjd1goZjxPPplddwelECQGCdAyqkoGHXy8YgTe0ZvupPijfRIVmgpJcznmQphqLIuIJhcFGnoruhp4NAxQfqyONQf1S5h2H57-vvmXnQk5tpdocXYC-MP3jFtmNukmdUWpsFlpr2Fclgy3d4opf2fDQzdC51vBpVl1DjKEngjGULtRo4jDy7VRKvrdHhNX25zeUQSsKyetlWARnn-O2RT_d7kYAbBBy195kqtplZ47QQjhptW8WBEfS8X0-wjOHM04gdW3p1iAJ4A88wYywy1T75zUMTH2iPiIHRilzwwPj5j4tWPiUCj__i8tQvLXIZVIIpV7jdP1yP9Kp_Vb2WV-DKy9osiImZotc_kAWxl5Jq6xqhKNAnRirWrwk1q_Z7KmPmnswC84Ao6h3Lqf728pR5NVQzFB2t5vWvFk-ocAx0gKNCGF0fug4PUS5t_M5WecFkLOrAx68fvRLfr7BA1JFAP6wPu4Alz0HbtixD1gUC6bHO4A8g7pb0lWoLE0a4hKkPnvrQjtV5ccjpVIj-4sgQLr9zIpYnPwxbzGg13DRGBPySKic7qrD4nBEktcep01q50' (err='<nil>')
	// The decoded token is 'abc123 xyz234 xyz345 2040-02-02T12:12:12Z' (err='<nil>')
}
