package jwt

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	DefaultJwtDecoder *JwtDecoder = getDecoderInstance()
	DefaultJwtEncoder *JwtEncoder = getEncoderInstance()
)

func getDecoderInstance() *JwtDecoder {
	jwkKeys := os.Getenv("AUTH_PUBLIC_JWK_KEYS")
	if jwkKeys == "" && isTestMode() {
		// test key only, not the production keys
		b, _ := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
		jwkKeys = string(b)
	}

	decoder, err := NewJwtDecoder(jwkKeys)
	if err != nil {
		err := fmt.Errorf("error loading default jwk decoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return decoder
}

func getEncoderInstance() *JwtEncoder {
	keyId := os.Getenv("AUTH_PRIVATE_KEY_ID")
	if keyId == "" {
		keyId = webGatewayKid
	}

	privKey := os.Getenv("AUTH_PRIVATE_KEY")
	if privKey == "" && isTestMode() {
		// test key only, not the production key
		b, _ := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
		privKey = string(b)
	}

	encoder, err := NewJwtEncoder(privKey, keyId)
	if err != nil {
		err := fmt.Errorf("error loading jwk encoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return encoder
}

// Decode a jwt token string and return the Standard Culture Amp Claims.
func Decode(tokenString string) (*StandardClaims, error) {
	return DefaultJwtDecoder.Decode(tokenString)
}

// Encode the Standard Culture Amp Claims in a jwt token string.
func Encode(claims *StandardClaims) (string, error) {
	return DefaultJwtEncoder.Encode(claims)
}

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		strings.Contains(argZero, "__debug_bin") || // vscode debug binary
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}
