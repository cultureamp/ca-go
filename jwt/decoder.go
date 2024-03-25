package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/patrickmn/go-cache"
)

const (
	kidHeaderKey                  = "kid"
	algorithmHeaderKey            = "alg"
	signatureHeaderKey            = "sig"
	webGatewayKid                 = "web-gateway"
	accountIDClaim                = "accountId"
	realUserIDClaim               = "realUserId"
	effectiveUserIDClaim          = "effectiveUserId"
	jwksCacheKey                  = "decoder_jwks_key"
	defaultDecoderExpiration      = 60 * time.Minute
	defaultDecoderCleanupInterval = 1 * time.Minute
)

type publicKey interface{} // Only ECDSA (perferred) and RSA public keys allowed

type DecoderJwksRetriever func() string

// JwtDecoder can decode a jwt token string.
type JwtDecoder struct {
	fetchJwkKeys      DecoderJwksRetriever // func provided by clients of this library to supply a refreshed JWKS
	mu                sync.Mutex           // mutex to protect cache.Get/Set race condition
	cache             *cache.Cache         // memory cache holding the jwk.Set
	defaultExpiration time.Duration        // default is 60 minutes
	cleanupInterval   time.Duration        // default is every 1 minute
}

// NewJwtDecoder creates a new JwtDecoder with the set ECDSA and RSA public keys in the JWK string.
func NewJwtDecoder(fetchJWKS DecoderJwksRetriever, options ...JwtDecoderOption) (*JwtDecoder, error) {
	decoder := &JwtDecoder{
		fetchJwkKeys:      fetchJWKS,
		defaultExpiration: defaultDecoderExpiration,
		cleanupInterval:   defaultDecoderCleanupInterval,
	}

	// Loop through our Decoder options and apply them
	for _, option := range options {
		option(decoder)
	}

	decoder.cache = cache.New(decoder.defaultExpiration, decoder.cleanupInterval)

	// call the getJWKS func to make sure its valid and we can parse the JWKS
	_, err := decoder.getJWKSet()
	if err != nil {
		return nil, fmt.Errorf("failed to load jwks: %w", err)
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
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return decoder.useCorrectPublicKey(token)
		},
		jwt.WithValidMethods(validAlgs), // only keys with these "alg's" will be considered
		jwt.WithLeeway(10*time.Second),  // as per the JWT eng std: clock skew set to 10 seconds
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

	// check cache and possibly fetch new JWKS
	jwkSet, err := d.getJWKSet()
	if err != nil {
		return nil, fmt.Errorf("failed to load jwks: %w", err)
	}

	key, found := jwkSet.LookupKeyID(kid)
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

func (d *JwtDecoder) getJWKSet() (jwk.Set, error) {
	// First check cache, if its there then great, use it!
	obj, found := d.cache.Get(jwksCacheKey)
	if found {
		jwks, ok := obj.(jwk.Set)
		if !ok {
			return nil, fmt.Errorf("internal error: cache key does not point to jwk.Set")
		}

		return jwks, nil
	}

	// The cache has expired the keys
	return d.refetchJWKSet()
}

func (d *JwtDecoder) refetchJWKSet() (jwk.Set, error) {
	// Only allow one thread to fetch, parse and update the cache
	d.mu.Lock()
	defer d.mu.Unlock()

	// check the cache again in case another go routine just updated it
	obj, found := d.cache.Get(jwksCacheKey)
	if found {
		jwks, ok := obj.(jwk.Set)
		if !ok {
			return nil, fmt.Errorf("internal error: cache key does not point to jwk.Set")
		}

		return jwks, nil
	}

	// Call client retriever func
	jwkKeys := d.fetchJwkKeys()

	// Parse all new JWKs JSON keys and make sure its valid
	jwkSet, err := d.parseJWKs(jwkKeys)
	if err != nil {
		return nil, err
	}

	// Add back into the cache
	err = d.cache.Add(jwksCacheKey, jwkSet, cache.DefaultExpiration)
	return jwkSet, err
}

func (decoder *JwtDecoder) parseJWKs(jwks string) (jwk.Set, error) {
	if jwks == "" {
		// If no jwks json, then returm empty map
		return nil, fmt.Errorf("missing jwks")
	}

	// 1. Parse the jwks JSON string to an iterable set
	return jwk.ParseString(jwks)
}
