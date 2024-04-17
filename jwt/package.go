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
	DefaultJwtEncoder Encoder = nil
	// DefaultJwtDecoder used for package level methods.
	// This can be mocked during tests if required by supporting the Decoder interface.
	DefaultJwtDecoder Decoder = nil
)

// Decode a jwt token string and return the Standard Culture Amp Claims.
func Decode(tokenString string) (*StandardClaims, error) {
	err := mustHaveDefaultJwtDecoder()
	if err != nil {
		return nil, err
	}
	return DefaultJwtDecoder.Decode(tokenString)
}

// DecodeWithCustomClaims takes a jwt token string and populate the customClaims.
func DecodeWithCustomClaims(tokenString string, customClaims jwt.Claims) error {
	err := mustHaveDefaultJwtDecoder()
	if err != nil {
		return err
	}
	return DefaultJwtDecoder.DecodeWithCustomClaims(tokenString, customClaims)
}

// Encode the Standard Culture Amp Claims in a jwt token string.
func Encode(claims *StandardClaims) (string, error) {
	err := mustHaveDefaultJwtEncoder()
	if err != nil {
		return "", err
	}
	return DefaultJwtEncoder.Encode(claims)
}

// EncodeWithCustomClaims encodes the Custom Claims in a jwt token string.
func EncodeWithCustomClaims(customClaims jwt.Claims) (string, error) {
	err := mustHaveDefaultJwtEncoder()
	if err != nil {
		return "", err
	}
	return DefaultJwtEncoder.EncodeWithCustomClaims(customClaims)
}

func mustHaveDefaultJwtDecoder() error {
	if DefaultJwtDecoder != nil {
		return nil // its set so we are good to go
	}

	decoder, err := NewJwtDecoder(jwksFromEnvVarRetriever)
	if err != nil {
		return errors.Errorf("error loading default jwk decoder, maybe missing env vars: err='%w'\n", err)
	}

	DefaultJwtDecoder = decoder
	return nil
}

func jwksFromEnvVarRetriever() string {
	jwkKeys, ok := os.LookupEnv("AUTH_PUBLIC_JWK_KEYS")
	if !ok || jwkKeys == "" {
		if isTestMode() {
			// If we are running inside a test, the make sure the DefaultJwtDecoder package level
			// instance loads these test keys be default.
			b, _ := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
			jwkKeys = string(b)
		}
	}

	return jwkKeys
}

func mustHaveDefaultJwtEncoder() error {
	if DefaultJwtEncoder != nil {
		return nil // its set so we are good to go
	}

	encoder, err := NewJwtEncoder(privateKeyFromEnvVarRetriever)
	if err != nil {
		return errors.Errorf("error loading jwk encoder, maybe missing env vars: err='%w'\n", err)
	}

	DefaultJwtEncoder = encoder
	return nil
}

func privateKeyFromEnvVarRetriever() (string, string) {
	privKey, ok := os.LookupEnv("AUTH_PRIVATE_KEY")
	if !ok || privKey == "" {
		if isTestMode() {
			// If we are running inside a test, the make sure the DefaultJwtEncoder package level
			// instance loads this test key by default.
			b, _ := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
			privKey = string(b)
		}
	}

	keyId, ok := os.LookupEnv("AUTH_PRIVATE_KEY_ID")
	if !ok || keyId == "" {
		if isTestMode() {
			// test key_id to web-gateway only, not the production key_id
			keyId = webGatewayKid
		}
	}

	return privKey, keyId
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
