# ca-go/jwt

The `jwt` package wraps JWT & JWKs `Encode` and `Decode` in a simple to use sington pattern that you can call directly. Only RSA public and private keys are currently supported (but this can easily be updated if needed in the future).

## Environment Variables

You can MUST set these:
- AUTH_PUBLIC_JWK_KEYS = A JSON string containing the well known public keys for Decoding a token.
- AUTH_PRIVATE_KEY = The private RSA PEM key used for Encoding a token.

You can OPTIONALLY set these:
- AUTH_PUBLIC_DEFAULT_KEY_ID = The default "kid" to use when no kid in present in the token (default to 'web-gateway`).
- AUTH_PRIVATE_KEY_ID = The default "kid" to add to the token heading when Encoding.

## Examples
```
package cago

import (
	"fmt"

	"github.com/cultureamp/ca-go/jwt"
)

func BasicExamples() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
	}

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	sc, err := jwt.Decode(token)
	fmt.Printf("The decode token is '%v' (err='%+v')\n", sc, err)
}
```
