package jwt

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	testAuthJwks       string = "./testKeys/development.jwks"
	testAuthPrivateKey string = "./testKeys/jwt-rsa256-test-webgateway.key"
)

var (
	DefaultJwtDecoder *JwtDecoder = getDecoderInstance()
	DefaultJwtEncoder *JwtEncoder = getEncoderInstance()
)

func getDecoderInstance() *JwtDecoder {
	keyId, ok := os.LookupEnv("AUTH_PUBLIC_DEFAULT_KEY_ID")
	if !ok {
		keyId = WebGatewayKid
	}

	jwkKeys, ok := os.LookupEnv("AUTH_PUBLIC_JWK_KEYS")
	if !ok && isTestMode() {
		// test key only, not the production key
		b, _ := os.ReadFile(filepath.Clean(testAuthJwks))
		jwkKeys = string(b)
		keyId = WebGatewayKid
	}

	decoder, err := NewJwtDecoderWithDefaultKid(jwkKeys, keyId)
	if err != nil {
		err := fmt.Errorf("error loading default jwk decoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return decoder
}

func getEncoderInstance() *JwtEncoder {
	keyId, ok := os.LookupEnv("AUTH_PRIVATE_KEY_ID")
	if !ok {
		keyId = WebGatewayKid
	}

	privKey, ok := os.LookupEnv("AUTH_PRIVATE_KEY")
	if !ok && isTestMode() {
		// test key only, not the production key
		b, _ := os.ReadFile(filepath.Clean(testAuthPrivateKey))
		privKey = string(b)
		keyId = WebGatewayKid
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
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}
