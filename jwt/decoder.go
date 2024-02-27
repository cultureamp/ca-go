package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const (
	kidHeaderKey       = "kid"
	algorithmHeaderKey = "alg"
	signatureHeaderKey = "sig"

	webGatewayKid = "web-gateway"

	publicKeyType = "RSA PUBLIC KEY"

	accountIDClaim       = "accountId"
	realUserIDClaim      = "realUserId"
	effectiveUserIDClaim = "effectiveUserId"
)

type (
	publicKey    interface{}          // Only ECDSA (perferred) and RSA public keys allowed
	publicKeyMap map[string]publicKey // "keyid => Public ECDSA/RSA Key".
)

// JwtDecoder can decode a jwt token string.
type JwtDecoder struct {
	jwkKeys publicKeyMap // public jwt's
}

// NewJwtDecoder creates a new JwtDecoder with the set ECDSA and RSA public keys in the JWK string.
func NewJwtDecoder(jwkKeys string) (*JwtDecoder, error) {
	decoder := &JwtDecoder{}
	decoder.jwkKeys = make(publicKeyMap)

	// 1. Parse all JWKs JSON keys
	publicKeyMap, err := decoder.parseJWKs(context.Background(), jwkKeys)
	if err != nil {
		return nil, err
	}

	// 2. Upsert into machine keys with "kid" as the key
	for key, val := range publicKeyMap {
		decoder.jwkKeys[key] = val
	}

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
	// Expiry claim is currently MANDATORY.
	// If the token includes an expiry claim, then the time is checked correctly and will return error if expired.
	// If the token does not include an expiry claim then returns an error.
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return decoder.useCorrectPublicKey(token)
		},
		jwt.WithLeeway(30*time.Second), // add this if we want to add some leeway for clock scew across systems
		jwt.WithExpirationRequired(),   // add this if we want to enforce that tokens MUST have an expiry
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

	kid, found := token.Header[kidHeaderKey]
	if !found {
		// no kid header but its MANDATORY
		return nil, fmt.Errorf("failed to decode: missing key_id (kid) header")
	}

	key, found := d.jwkKeys[kid.(string)]
	if found {
		// Found a match, so use this key
		return key, nil
	}

	// Didn't find a matching kid
	return nil, fmt.Errorf("failed to decode: no matching key_id (kid) header for: %s", kid)
}

func (decoder *JwtDecoder) parseJWKs(ctx context.Context, jwks string) (publicKeyMap, error) {
	rsaKeys := make(publicKeyMap)

	if jwks == "" {
		// If no jwks json, then returm empty map
		return rsaKeys, fmt.Errorf("missing jwks")
	}

	// 1. Parse the jwks JSON string to an iterable set
	jwkset, err := jwk.ParseString(jwks)
	if err != nil {
		return rsaKeys, err
	}

	// 2. Enumerate the set
	for it := jwkset.Keys(ctx); it.Next(ctx); {
		pair := it.Pair()
		key, ok := pair.Value.(jwk.Key)
		if !ok {
			// the jwks Set value isn't valid (for some reason) - just skip it
			continue
		}

		usage := key.KeyUsage()
		if usage != signatureHeaderKey {
			// we aren't interested if it isn't a "sig"
			continue
		}

		kid := key.KeyID()

		var rsa rsa.PublicKey
		if err := key.Raw(&rsa); err != nil {
			// We only support RSA keys currently so skip if not a RSA public key
			continue
		}

		pubKeyBytes, err := x509.MarshalPKIXPublicKey(&rsa)
		if err != nil {
			return rsaKeys, err
		}

		pubKeyPEM := pem.EncodeToMemory(
			&pem.Block{
				Type:  publicKeyType,
				Bytes: pubKeyBytes,
			},
		)

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyPEM)
		if err != nil {
			return rsaKeys, err
		}

		// 3. add public key to the map, overwriting if already exists
		rsaKeys[kid] = publicKey
	}

	// 4. return all the valid RSA keys
	return rsaKeys, nil
}
