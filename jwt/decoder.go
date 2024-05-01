package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"sync"
	"time"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const (
	kidHeaderKey                   = "kid"
	algorithmHeaderKey             = "alg"
	webGatewayKid                  = "web-gateway"
	accountIDClaim                 = "accountId"
	realUserIDClaim                = "realUserId"
	effectiveUserIDClaim           = "effectiveUserId"
	defaultDecoderExpiration       = 60 * time.Minute
	defaultDecoderCleanupInterval  = 1 * time.Minute
	defaultDecoderRotationDuration = 30 * time.Second
)

type publicKey interface{} // Only ECDSA (perferred) and RSA public keys allowed

// DecoderJwksRetriever defines the function signature required to retrieve JWKS json.
type DecoderJwksRetriever func() string

// JwtDecoder can decode a jwt token string.
type JwtDecoder struct {
	fetchJwkKeys   DecoderJwksRetriever // func provided by clients of this library to supply a refreshed JWKS
	expiresWithin  time.Duration        // default is 60 minutes
	rotationWindow time.Duration        // default is 30 seconds

	mu          sync.Mutex // mutex to protect race conditions on jwks and jwksAddedAt
	jwks        jwk.Set    // the current JWK Set of jwt public keys
	jwksAddedAt time.Time  // the time when we fetched these keys (so we can refresh)
}

// NewJwtDecoder creates a new JwtDecoder with the set ECDSA and RSA public keys in the JWK string.
func NewJwtDecoder(fetchJWKS DecoderJwksRetriever, options ...JwtDecoderOption) (*JwtDecoder, error) {
	decoder := &JwtDecoder{
		fetchJwkKeys:   fetchJWKS,
		jwks:           nil,
		expiresWithin:  defaultDecoderExpiration,
		rotationWindow: defaultDecoderRotationDuration,
	}

	// Loop through our Decoder options and apply them
	for _, option := range options {
		option(decoder)
	}

	// call the getJWKS func to make sure its valid and we can parse the JWKS
	_, err := decoder.getCurrentJWKs()
	if err != nil {
		return nil, errors.Errorf("failed to load jwks: %w", err)
	}

	return decoder, nil
}

// Decode a jwt token string and return the Standard Culture Amp Claims.
func (d *JwtDecoder) Decode(tokenString string) (*StandardClaims, error) {
	claims := jwt.MapClaims{}
	err := d.DecodeWithCustomClaims(tokenString, claims)
	if err != nil {
		return nil, err
	}
	return newStandardClaims(claims), nil
}

// DecodeWithCustomClaims takes a jwt token string and populate the customClaims.
func (d *JwtDecoder) DecodeWithCustomClaims(tokenString string, customClaims jwt.Claims) error {
	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	validAlgs := []string{"RS256", "RS384", "RS512", "ES256", "ES384", "ES512"}

	// sample token string in the form "header.payload.signature"
	// eg. "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Eng Std: https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3253240053/JWT+Authentication

	// Exp
	// Expiry claim is currently MANDATORY, but until all producing services are reliably setting the Expiry claim,
	// we MAY still accept verified JWTs with no Expiry claim.
	// Nbf
	// NotBefore claim is currently MANDATORY, but until all producing services are reliably settings the NotBEfore claim,
	// we MAY still accept verificed JWT's with no NotBefore claim.
	token, err := jwt.ParseWithClaims(
		tokenString,
		customClaims,
		func(token *jwt.Token) (interface{}, error) {
			return d.useCorrectPublicKey(token)
		},
		jwt.WithValidMethods(validAlgs), // only keys with these "alg's" will be considered
		jwt.WithLeeway(10*time.Second),  // as per the JWT eng std: clock skew set to 10 seconds
		// jwt.WithExpirationRequired(),	// add this if we want to enforce that tokens MUST have an expiry
	)
	if err != nil || !token.Valid {
		return err
	}

	return nil
}

func (d *JwtDecoder) useCorrectPublicKey(token *jwt.Token) (publicKey, error) {
	if token == nil {
		return nil, errors.Errorf("failed to decode: missing token")
	}

	// Eng Std: https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3253240053/JWT+Authentication
	// Perferred is ECDSA, but is RSA accepted
	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected signing method - only ecdsa or rsa supported: %v", token.Header[algorithmHeaderKey])
		}
	}

	kidHeader, found := token.Header[kidHeaderKey]
	if !found {
		// no kid header but its MANDATORY
		return nil, errors.Errorf("failed to decode: missing key_id (kid) header")
	}

	kid, ok := kidHeader.(string)
	if !ok {
		// kid header isn't a string?!
		return nil, errors.Errorf("failed to decode: invalid key_id (kid) header")
	}

	// check if kid exists in the JWK Set
	key, err := d.lookupKeyID(kid)
	if err != nil {
		// if the JWKS is at least d.rotationWindow (default is 30 seconds) old
		if d.safeRefreshJWKs() {
			// then its safe to refresh the JWKS to check if a new key has been added/rotated
			key, err = d.lookupKeyID(kid)
		}
	}

	return key, err
}

// lookupKeyID returns the public key in the JWKS that matches the "kid".
func (d *JwtDecoder) lookupKeyID(kid string) (publicKey, error) {
	// check cache and possibly fetch new JWKS if cache has expired
	jwkSet, err := d.getCurrentJWKs()
	if err != nil {
		return nil, errors.Errorf("failed to load jwks: %w", err)
	}

	key, found := jwkSet.LookupKeyID(kid)
	if found {
		// Found a match, so use this key
		return d.getPublicKey(key)
	}
	return nil, errors.Errorf("failed to decode: no matching key_id (kid) header for: %s", kid)
}

// safeRefreshJWKs is ONLY called when a "kid" is missing from the JWK Set.
// The purpose of this method is to remove the JWK Set IF it is older than 30 secs.
func (d *JwtDecoder) safeRefreshJWKs() bool {
	// Only allow one thread to update the jwks
	d.mu.Lock()
	defer d.mu.Unlock()

	freshness := time.Since(d.jwksAddedAt)
	if freshness > d.rotationWindow {
		// only rotate keys if we haven't updated in the last 30 secs to stop a bunch of requests heading to
		// FushionAuth if key is missing
		d.jwks = nil
		return true
	}

	return false
}

// getCurrentJWKs will check if the JWKS have expired.
// If not, then it returns it.
// Otherwise, it re-fetches, parses, and updates the decoder.
func (d *JwtDecoder) getCurrentJWKs() (jwk.Set, error) {
	// Only allow one thread to update the jwks
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.jwks != nil && time.Now().Before(d.jwksAddedAt.Add(d.expiresWithin)) {
		// we have jwks and it hasn't expired yet, so all good!
		return d.jwks, nil
	}

	// Call client retriever func
	jwkKeys := d.fetchJwkKeys()

	// Parse all new JWKs JSON keys and make sure its valid
	jwkSet, err := d.parseJWKs(jwkKeys)
	if err != nil {
		return nil, err
	}

	// update with latest values
	d.jwksAddedAt = time.Now()
	d.jwks = jwkSet
	return d.jwks, nil
}

func (decoder *JwtDecoder) parseJWKs(jwks string) (jwk.Set, error) {
	if jwks == "" {
		// If no jwks json, then returm empty map
		return nil, errors.Errorf("missing jwks")
	}

	// 1. Parse the jwks JSON string to an iterable set
	return jwk.ParseString(jwks)
}

func (d *JwtDecoder) getPublicKey(key jwk.Key) (publicKey, error) {
	var rawkey interface{}
	err := key.Raw(&rawkey)
	if err != nil {
		return nil, errors.Errorf("failed to decode: bad public key in jwks")
	}

	// If the JWKS contains the full key (Private AND Public) then only return the public one
	// ECDSA & RSA keys only.
	// NOTE: this should never happen in production - but does in the unit tests
	if ecdsa, ok := rawkey.(*ecdsa.PrivateKey); ok {
		return &ecdsa.PublicKey, nil
	}
	if rsa, ok := rawkey.(*rsa.PrivateKey); ok {
		return &rsa.PublicKey, nil
	}

	return rawkey, err
}
