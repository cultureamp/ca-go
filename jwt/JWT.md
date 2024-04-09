# ca-go/jwt

The `jwt` package wraps JWT & JWKs `Encode` and `Decode` in a simple to use sington pattern that you can call directly. Only ECDSA and RSA public and private keys are currently supported (but this can easily be updated if needed in the future).

## Environment Variables

To use the package level methods `Encode` and `Decode` you MUST set these:

- AUTH_PUBLIC_JWK_KEYS = A JSON string containing the well known public keys for Decoding a token.
- AUTH_PRIVATE_KEY = The private RSA PEM key used for Encoding a token.
- AUTH_PRIVATE_KEY_ID = The "kid" (key_id) header to add to the token heading when Encoding.

Failure to set these will result in a panic() at start up.

## Runtime Key Rotation

We want to be able to easily rotate keys without having to stop/start all services. The simplest approach is that the JWT Encoder and Decoder take a func() that clients must implement to provide a refreshed key.

If you are using the package level methods, then the `DefaultJwtEncoder` and `DefaultJwtDecoder` will check there environment variables every 60 minutes. So all you need to do is update them with `os.Setenv()` with new values if a key rotation occurs.

If you are managing Encoders and Decoders yourself, then you can provide a func of type `EncoderKeyRetriever` to the `NewJwtEncoder` constructor, and a func of type `DecoderJwksRetriever` to the `NewJwtDecoder`:

- type EncoderKeyRetriever func() (string, string)
- type DecoderJwksRetriever func() string


## Managing Encoders and Decoders Yourself

While we recommend using the package level methods for their ease of use, you may desire to create and manage encoders or decoers yourself, which you can do by calling:

```
func privateKeyRetriever() (string, string) {
	// todo: check if keys have rotated (eg. re-read secrets manager)

	privKey := secrets.Get("my-private-key")
	keyId := secrets.Get("my-private-key-id")
	return privKey, keyId
}

func jwksRetriever() string {
	// todo: check if keys have rotated (eg. call well-known URL)

	resp, err := http.Get("http://well-known-example.com/list.jwks")
	// todo: handle error (eg. retry)
	defer resp.Body.Close()
	jwkKeys, err := io.ReadAll(resp.Body)
	return jwkKeys
}

func main() {

encoder, err := NewJwtEncoder(privateKeyRetriever)
decoder, err := NewJwtDecoder(jwksRetriever)
}
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

## Testing and Mocks

During tests you can override the package level `DefaultJwtEncoder` and/or `DefaultJwtDecoder` with a mock that supports
the `Encoder` or `Decoder` interface.

- Encode(claims *StandardClaims) (string, error)
- Decode(tokenString string) (*StandardClaims, error)

```
import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/jwt"
	"github.com/stretchr/testify/mock"
)

func ExampleMocked_EncoderDecoder() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	mockEncDec := newMockedEncoderDecoder()
	mockEncDec.On("Encode", mock.Anything).Return("eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0", nil)
	mockEncDec.On("Decode", mock.Anything).Return(claims, nil)

	// Overwrite the Default package level encoder and decoder
	oldEncoder := jwt.DefaultJwtEncoder
	oldDecoder := jwt.DefaultJwtDecoder
	jwt.DefaultJwtEncoder = mockEncDec
	jwt.DefaultJwtDecoder = mockEncDec
	defer func() {
		jwt.DefaultJwtEncoder = oldEncoder
		jwt.DefaultJwtDecoder = oldDecoder
	}()

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	claim, err := jwt.Decode(token)
	fmt.Printf(
		"The decoded token is '%s %s %s %s %v %s %s' (err='%+v')\n",
		claim.AccountId, claim.RealUserId, claim.EffectiveUserId,
		claim.Issuer, claim.Subject, claim.Audience,
		claim.ExpiresAt.UTC().Format(time.RFC3339),
		err,
	)

	// Output:
	// The encoded token is 'eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0' (err='<nil>')
	// The decoded token is 'abc123 xyz234 xyz345 encoder-name test [decoder-name] 2040-02-02T12:12:12Z' (err='<nil>')
}

type mockedEncoderDecoder struct {
	mock.Mock
}

func newMockedEncoderDecoder() *mockedEncoderDecoder {
	return &mockedEncoderDecoder{}
}

func (m *mockedEncoderDecoder) Encode(claims *jwt.StandardClaims) (string, error) {
	args := m.Called(claims)
	output, _ := args.Get(0).(string)
	return output, args.Error(1)
}

// Decrypt on the test runner just returns the "encryptedStr" as the decrypted plainstr.
func (m *mockedEncoderDecoder) Decode(tokenString string) (*jwt.StandardClaims, error) {
	args := m.Called(tokenString)
	output, _ := args.Get(0).(*jwt.StandardClaims)
	return output, args.Error(1)
}
```
