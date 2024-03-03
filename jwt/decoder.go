package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const (
	kidHeaderKey         = "kid"
	algorithmHeaderKey   = "alg"
	signatureHeaderKey   = "sig"
	webGatewayKid        = "web-gateway"
	accountIDClaim       = "accountId"
	realUserIDClaim      = "realUserId"
	effectiveUserIDClaim = "effectiveUserId"
)

type (
	publicKey interface{} // Only ECDSA (perferred) and RSA public keys allowed
)

// JwtDecoder can decode a jwt token string.
type JwtDecoder struct {
	jwkSet jwk.Set // public jwt's
}

// NewJwtDecoder creates a new JwtDecoder with the set ECDSA and RSA public keys in the JWK string.
func NewJwtDecoder(jwkKeys string) (*JwtDecoder, error) {
	decoder := &JwtDecoder{}

	// 1. Parse all JWKs JSON keys
	jwkSet, err := decoder.parseJWKs(jwkKeys)
	if err != nil {
		return nil, err
	}

	decoder.jwkSet = jwkSet
	return decoder, nil
}

// Decode a jwt token string and return the Standard Culture Amp Claims.
func (d *JwtDecoder) Decode(tokenString string) (*StandardClaims, error) {
	payload := &StandardClaims{}

	claims, err := d.decodeClaims(tokenString)
	if err != nil {
		return payload, err
	}

	return newStandardClaims(claims), nil
}

func (decoder *JwtDecoder) decodeClaims(tokenString string) (jwt.MapClaims, error) {
	// sample token string in the form "header.payload.signature"
	// eg. "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Eng Std: https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3253240053/JWT+Authentication
	// Expiry claim is currently MANDATORY, but until all producing services are reliably setting the Expiry claim,
	// we MAY still accept verified JWTs with no Expiry claim.
	// So:
	// If the token includes an expiry claim, then the time is checked correctly and will return error if expired.
	// If the token does not include an expiry claim then ignore and just test that verification is valid.
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return decoder.useCorrectPublicKey(token)
		},
		jwt.WithLeeway(10*time.Second), // as per the JWT eng std: clock skew set to 10 seconds
		// jwt.WithExpirationRequired(),	// add this if we want to enforce that tokens MUST have an expiry
	)
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("missing claims in jwt token")
	}

	return claims, nil
}

func (d *JwtDecoder) useCorrectPublicKey(token *jwt.Token) (publicKey, error) {
	if token == nil {
		return nil, fmt.Errorf("failed to decode: missing token")
	}

	// Eng Std: https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3253240053/JWT+Authentication
	// Perferred is ECDSA, but is RSA accepted
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method - only ecdsa or rsa supported: %v", token.Header[algorithmHeaderKey])
		}
	}

	kidHeader, found := token.Header[kidHeaderKey]
	if !found {
		// no kid header but its MANDATORY
		return nil, fmt.Errorf("failed to decode: missing key_id (kid) header")
	}

	kid, ok := kidHeader.(string)
	if !ok {
		// kid header isn't a string?!
		return nil, fmt.Errorf("failed to decode: invalid key_id (kid) header")
	}

	key, found := d.jwkSet.LookupKeyID(kid)
	if found {
		// Found a match, so use this key
		var rawkey interface{}
		err := key.Raw(&rawkey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: bad public key in jwks")
		}

		// If the JWKS contains the full key (Private AND Public) then check for that for both ECDSA & RSA
		// NOTE: this should never happen in PRPD - but does in the unit tests
		if ecdsa, ok := rawkey.(*ecdsa.PrivateKey); ok {
			return &ecdsa.PublicKey, nil
		}
		if rsa, ok := rawkey.(*rsa.PrivateKey); ok {
			return &rsa.PublicKey, nil
		}

		return rawkey, err
	}

	// Didn't find a matching kid
	return nil, fmt.Errorf("failed to decode: no matching key_id (kid) header for: %s", kid)
}

func (decoder *JwtDecoder) parseJWKs(jwks string) (jwk.Set, error) {
	if jwks == "" {
		// If no jwks json, then returm empty map
		return nil, fmt.Errorf("missing jwks")
	}

	// 1. Parse the jwks JSON string to an iterable set
	return jwk.ParseString(jwks)
}
