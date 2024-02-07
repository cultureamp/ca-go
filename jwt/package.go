package jwt

import (
	"fmt"
	"os"
)

var (
	DefaultJwtDecoder *JwtDecoder = getDecoderInstance()
	DefaultJwtEncoder *JwtEncoder = getEncoderInstance()
)

func getDecoderInstance() *JwtDecoder {
	jwkKeys := os.Getenv("AUTH_PUBLIC_JWK_KEYS")

	decoder, err := NewJwtDecoder(jwkKeys)
	if err != nil {
		err := fmt.Errorf("error loading jwk decoder, maybe missing env vars: err='%w'\n", err)
		panic(err)
	}

	return decoder
}

func getEncoderInstance() *JwtEncoder {
	privKey := os.Getenv("AUTH_PRIVATE_KEY")
	keyId := os.Getenv("AUTH_PRIVATE_KEY_ID")

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
