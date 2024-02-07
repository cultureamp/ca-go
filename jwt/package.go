package jwt

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	DefaultJwtDecoder *JwtDecoder = getDecoderInstance()
	DefaultJwtEncoder *JwtEncoder = getEncoderInstance()
)

func getDecoderInstance() *JwtDecoder {
	jwkKeys, ok := os.LookupEnv("AUTH_PUBLIC_JWK_KEYS")
	if !ok && isTestMode() {
		jwkKeys = `{ "keys" : [{
				"alg":"RS256",
				"kty":"RSA",
				"e":"AQAB",
				"kid":"web-gateway",
				"n":"zkzpPa8QB5JwYWJI1W3WmxnMwusvFZb-0EVY4Sko3C1zwBcY8P6NucHo1epXTO-rFQy8JPiSMyTBINkmDP0d1jfvJF_RDL8Gzi1_aM2mScsPxmXA7ftqHdvcaqP0aobuYNJSEk_3erM6iddBJwsKY5BNkzS-R9szsfCgnDdfN-9JvChpfrTvoOwI-vtsqpkgIgGB4uCeQ0CPvqZzsRMJyWouEt0Jj7huKXBOvDBuoZdInuh-2kzNpm9KEkdbB0wzhC57MnyA3ap0I-ES374utQGM1EbZfW68T0QU3t--Q7L7yQ4D8WjRLZw_WTS8amcLRYf0urb3yTmvQFA4ryhc25dBUF68xPrC2kETljf6SLtig2bWvr-TGqGiyLnqiPloSxeBtpZhWSBgH8KJ7iHjwCyT2dSMEhf-ouivT2rEn5wEP3joDPywBqywKs-hbJrOB_x9cg4dGqERuljvW02tMGHu1JTK8tb23wWl8_5RSPHGetM526G3MW8r8hJ4mPHASPzQ2jWM_XhHtvLOg4_0V3CczMe93e6ilWkxala1hnZA180lOFoOOscdQmcH7LbOjkH6Iwb_9lc0Ez6n2tcfuY9p1aujcsJ5uQNBJtoX4kOSTM7LfUJa88ZbUkOeJ9AHhCe9xqaAS-W0LJYR00-JZcsaZz31F2DSFMmOWLUCVZ8",
				"use":"sig"
			}]}`
	}

	decoder, err := NewJwtDecoder(jwkKeys)
	if err != nil {
		err := fmt.Errorf("error loading default jwk decoder, maybe missing env vars: err='%w'\n", err)
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
