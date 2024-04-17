package jwt

import (
	"os"

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

	decoder, err := NewJwtDecoder(func() string { return os.Getenv("AUTH_PUBLIC_JWK_KEYS") })
	if err != nil {
		return errors.Errorf("error loading default jwk decoder, maybe missing env vars: err='%w'\n", err)
	}

	DefaultJwtDecoder = decoder
	return nil
}

func mustHaveDefaultJwtEncoder() error {
	if DefaultJwtEncoder != nil {
		return nil // its set so we are good to go
	}

	encoder, err := NewJwtEncoder(func() (string, string) {
		privKey := os.Getenv("AUTH_PRIVATE_KEY")
		keyId := os.Getenv("AUTH_PRIVATE_KEY_ID")
		return privKey, keyId
	})
	if err != nil {
		return errors.Errorf("error loading jwk encoder, maybe missing env vars: err='%w'\n", err)
	}

	DefaultJwtEncoder = encoder
	return nil
}
