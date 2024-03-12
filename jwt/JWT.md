# ca-go/jwt

The `jwt` package wraps JWT & JWKs `Encode` and `Decode` in a simple to use sington pattern that you can call directly. Only ECDSA and RSA public and private keys are currently supported (but this can easily be updated if needed in the future).

## Environment Variables

To use the package level methods `Encode` and `Decode` you MUST set these:

- AUTH_PUBLIC_JWK_KEYS = A JSON string containing the well known public keys for Decoding a token.
- AUTH_PRIVATE_KEY = The private RSA PEM key used for Encoding a token.
- AUTH_PRIVATE_KEY_ID = The "kid" (key_id) header to add to the token heading when Encoding.

## Managing Encoders and Decoders Yourself

While we recommend using the package level methods for their ease of use, you may desire to create and manage encoders or decoers yourself, which you can do by calling:

```
privKey := os.Getenv("AUTH_PRIVATE_KEY")
encoder, err := NewJwtEncoder(tprivKey, "kid")

jwkKeys := os.Getenv("AUTH_PUBLIC_JWK_KEYS")
decoder, err := NewJwtDecoder(jwkKeys)
```

## Claims

You MUST set the `Issuer`, `Subject`, and `Audience` claims along with the standard authentication values of `AccountId`, `RealUserId`, and `EffectiveUserId`.

- [Issuer](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1) `iss` claim.
- [Subject](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2) `sub` claim.
- [Audience](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3) `aud`claim.

Please read the [JWT Engineering Standard](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3253240053/JWT+Authentication) for more information and details.

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
		Issuer:          "name-of-the-encoder",
		Subject:         "name-of-this-jwt-token",
		Audience:        []string{"list-of-intended-decoders-1", "list-of-intended-decoders-2"},
	}

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	sc, err := jwt.Decode(token)
	fmt.Printf("The decode token is '%v' (err='%+v')\n", sc, err)
}
```
