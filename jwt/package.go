package jwt

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Encoder interface allows for mocking of the Encoder.
type Encoder interface {
	Encode(claims *StandardClaims) (string, error)
}

// Decoder interface allows for mocking of the Decoder.
type Decoder interface {
	Decode(tokenString string) (*StandardClaims, error)
}

var (
	DefaultJwtEncoder Encoder = getEncoderInstance()
	DefaultJwtDecoder Decoder = getDecoderInstance()
)

func getDecoderInstance() *JwtDecoder {
	decoder, err := NewJwtDecoder(jwksFromEnvVarRetriever)
	if err != nil {
		err := fmt.Errorf("error loading default jwk decoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return decoder
}

func jwksFromEnvVarRetriever() string {
	jwkKeys := os.Getenv("AUTH_PUBLIC_JWK_KEYS")

	if isTestMode() {
		// If we are running inside a test, the make sure the DefaultJwtDecoder package level
		// instance doesn't panic with missing values.
		if jwkKeys == "" {
			// test key only, not the production keys
			b, _ := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
			jwkKeys = string(b)
		}
	}

	return jwkKeys
}

func getEncoderInstance() *JwtEncoder {
	encoder, err := NewJwtEncoder(privateKeyFromEnvVarRetriever)
	if err != nil {
		err := fmt.Errorf("error loading jwk encoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return encoder
}

func privateKeyFromEnvVarRetriever() (string, string) {
	privKey := os.Getenv("AUTH_PRIVATE_KEY")
	keyId := os.Getenv("AUTH_PRIVATE_KEY_ID")

	if isTestMode() {
		// If we are running inside a test, the make sure the DefaultJwtEncoder package level
		// instance doesn't panic with missing values.
		if privKey == "" {
			// test key only, not the production key
			b, _ := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
			privKey = string(b)
		}

		if keyId == "" {
			keyId = webGatewayKid
		}
	}

	return privKey, keyId
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

