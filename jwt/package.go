package jwt

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
)

// Encoder interface allows for mocking of the Encoder.
type Encoder interface {
	Encode(claims *StandardClaims) (string, error)
	EncodeWithCustomClaims(customClaims jwt.Claims) (string, error)
}

// Decoder interface allows for mocking of the Decoder.
type Decoder interface {
	Decode(tokenString string) (*StandardClaims, error)
	DecodeWithCustomClaims(tokenString string, customClaims jwt.Claims) error
}

var (
	// DefaultJwtEncoder used to package level methods.
	// This can be mocked during tests if required by supporting the Encoder interface.
	DefaultJwtEncoder Encoder = getEncoderInstance()
	// DefaultJwtDecoder used for package level methods.
	// This can be mocked during tests if required by supporting the Decoder interface.
	DefaultJwtDecoder Decoder = getDecoderInstance()
)

func getDecoderInstance() *JwtDecoder {
	decoder, err := NewJwtDecoder(jwksFromEnvVarRetriever)
	if err != nil {
		err := errors.Errorf("error loading default jwk decoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return decoder
}

func jwksFromEnvVarRetriever() string {
	jwkKeys, ok := os.LookupEnv("AUTH_PUBLIC_JWK_KEYS")
	if !ok || jwkKeys == "" {
		if !isTestMode() {
			err := errors.Errorf("missing AUTH_PUBLIC_JWK_KEYS environment variable - this should be set to a JWKS json string.")
			panic(err)
		}
		// If we are running inside a test, the make sure the DefaultJwtDecoder package level
		// instance doesn't panic with missing values.
		// test key only, not the production keys
		b, _ := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
		jwkKeys = string(b)
	}

	return jwkKeys
}

func getEncoderInstance() *JwtEncoder {
	encoder, err := NewJwtEncoder(privateKeyFromEnvVarRetriever)
	if err != nil {
		err := errors.Errorf("error loading jwk encoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return encoder
}

func privateKeyFromEnvVarRetriever() (string, string) {
	privKey, ok := os.LookupEnv("AUTH_PRIVATE_KEY")
	if !ok || privKey == "" {
		if !isTestMode() {
			err := errors.Errorf("missing AUTH_PRIVATE_KEY environment variable - this should be set to a private PEM key for this service.")
			panic(err)
		}
		// If we are running inside a test, the make sure the DefaultJwtEncoder package level
		// instance doesn't panic with missing values.
		// test key only, not the production key
		b, _ := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
		privKey = string(b)
	}

	keyId, ok := os.LookupEnv("AUTH_PRIVATE_KEY_ID")
	if !ok || keyId == "" {
		if !isTestMode() {
			err := errors.Errorf("missing AUTH_PRIVATE_KEY_ID environment variable - this should be set key_id for this service.")
			panic(err)
		}
		// test key_id to web-gateway only, not the production key_id
		keyId = webGatewayKid
	}

	return privKey, keyId
}

// Decode a jwt token string and return the Standard Culture Amp Claims.
func Decode(tokenString string) (*StandardClaims, error) {
	return DefaultJwtDecoder.Decode(tokenString)
}

// DecodeWithCustomClaims takes a jwt token string and populate the customClaims.
func DecodeWithCustomClaims(tokenString string, customClaims jwt.Claims) error {
	return DefaultJwtDecoder.DecodeWithCustomClaims(tokenString, customClaims)
}

// Encode the Standard Culture Amp Claims in a jwt token string.
func Encode(claims *StandardClaims) (string, error) {
	return DefaultJwtEncoder.Encode(claims)
}

// EncodeWithCustomClaims encodes the Custom Claims in a jwt token string.
func EncodeWithCustomClaims(customClaims jwt.Claims) (string, error) {
	return DefaultJwtEncoder.EncodeWithCustomClaims(customClaims)
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
