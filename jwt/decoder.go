package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"time"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const (
	webgatewayKid                  = "web-gateway"
	kidHeaderKey                   = "kid"
	algorithmHeaderKey             = "alg"
	accountIDClaim                 = "accountId"
	realUserIDClaim                = "realUserId"
	effectiveUserIDClaim           = "effectiveUserId"
	defaultDecoderExpiration       = 60 * time.Minute
	defaultDecoderRotationDuration = 30 * time.Second
	defaultDecoderLeeway           = 10 * time.Second
)

type publicKey interface{} // Only ECDSA (perferred) and RSA public keys allowed

// DecoderJwksRetriever defines the function signature required to retrieve JWKS json.
type DecoderJwksRetriever func() string

// StandardDecoder can decode a jwt token string.
type StandardDecoder struct {
	dispatcher     DecoderJwksRetriever // func provided by clients of this library to supply the current JWKS
	expiresWithin  time.Duration        // default is 60 minutes
	rotationWindow time.Duration        // default is 30 seconds
	jwks           *jwkFetcher          // manages the life cycle of a JWK Set
}

// NewDecoder creates a new JwtDecoder with the set ECDSA and RSA public keys in the JWK string.
func NewDecoder(fetchJWKS DecoderJwksRetriever, options ...DecoderOption) (*StandardDecoder, error) {
	decoder := &StandardDecoder{
		dispatcher:     fetchJWKS,
		jwks:           nil,
		expiresWithin:  defaultDecoderExpiration,
		rotationWindow: defaultDecoderRotationDuration,
	}

	// Loop through our Decoder options and apply them
	for _, option := range options {
		option(decoder)
	}

	decoder.jwks = newJWKSet(fetchJWKS, decoder.expiresWithin, decoder.rotationWindow)

	// call the get to make sure its valid and we can parse the JWKS
	_, err := decoder.jwks.Get()
	if err != nil {
		return nil, errors.Errorf("failed to load jwks: %w", err)
	}

	return decoder, nil
}

// Decode a jwt token string and return the Standard Culture Amp Claims.
func (d *StandardDecoder) Decode(tokenString string, options ...DecoderParserOption) (*StandardClaims, error) {
	claims := jwt.MapClaims{}
	err := d.DecodeWithCustomClaims(tokenString, claims, options...)
	if err != nil {
		return nil, err
	}
	return newStandardClaims(claims), nil
}

// DecodeWithCustomClaims takes a jwt token string and populate the customClaims.
func (d *StandardDecoder) DecodeWithCustomClaims(tokenString string, customClaims jwt.Claims, options ...DecoderParserOption) error {
	// sample token string in the form "header.payload.signature"
	// eg. "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	token, err := jwt.ParseWithClaims(
		tokenString,
		customClaims,
		func(token *jwt.Token) (interface{}, error) {
			return d.useCorrectPublicKey(token)
		},
		d.enforceParsingOptions(options...)...,
	)
	if err != nil || !token.Valid {
		return err
	}

	return nil
}

func (d *StandardDecoder) enforceParsingOptions(options ...DecoderParserOption) []jwt.ParserOption {
	// Eng Std: https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3253240053/JWT+Authentication
	var opts []jwt.ParserOption

	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	validAlgs := []string{"RS256", "RS384", "RS512", "ES256", "ES384", "ES512"}
	opts = append(opts,
		jwt.WithValidMethods(validAlgs),      // only keys with these "alg's" will be considered
		jwt.WithLeeway(defaultDecoderLeeway), // as per the JWT eng std: clock skew set to 10 seconds

		// Exp
		// Expiry claim is currently MANDATORY, but until all producing services are reliably setting the Expiry claim,
		// we MAY still accept verified JWTs with no Expiry claim.
		// jwt.WithExpirationRequired(),

		// Nbf
		// If the NotBefore claim is set it will automatically be enforced.
		// Note: There is no parsing option for this.
	)

	// Loop through any client provided parsing options and apply them
	parserOptions := newDecoderParser()
	for _, option := range options {
		option(parserOptions)
	}

	if parserOptions.expectedAud != "" {
		opts = append(opts, jwt.WithAudience(parserOptions.expectedAud))
	}

	if parserOptions.expectedIss != "" {
		opts = append(opts, jwt.WithIssuer(parserOptions.expectedIss))
	}

	if parserOptions.expectedSub != "" {
		opts = append(opts, jwt.WithSubject(parserOptions.expectedSub))
	}

	return opts
}

func (d *StandardDecoder) useCorrectPublicKey(token *jwt.Token) (publicKey, error) {
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
		// Currently, the web-gateway doesn't encode the "kid" in its tokens (we should fix that)
		// So until we do, instead of returning an error here, we default to trying the "web-gateway" kid in the JWKs.

		kidHeader = webgatewayKid
		// return nil, errors.Errorf("failed to decode: missing key_id (kid) header")
	}

	kid, ok := kidHeader.(string)
	if !ok {
		// kid header isn't a string?!
		return nil, errors.Errorf("failed to decode: invalid key_id (kid) header")
	}

	// check if kid exists in the JWK Set
	return d.lookupKeyID(kid)
}

// lookupKeyID returns the public key in the JWKS that matches the "kid".
func (d *StandardDecoder) lookupKeyID(kid string) (publicKey, error) {
	// check cache and possibly fetch new JWKS if cache has expired
	jwkSet, err := d.jwks.Get()
	if err != nil {
		return nil, errors.Errorf("failed to load jwks: %w", err)
	}

	// set if the kid exists in the set
	key, found := jwkSet.LookupKeyID(kid)
	if found {
		// Found a match, so use this key!
		return d.getPublicKey(key)
	}

	return d.tryRefreshedLookupKeyID(kid)
}

func (d *StandardDecoder) tryRefreshedLookupKeyID(kid string) (publicKey, error) {
	// If the jwks aren't "fresh" and we are being asked for a kid we don't have
	// then get a new jwks and try again. This can occur when a new key has been
	// added or rotated and we haven't got the latest copy.
	// The "canRefresh" check is important here, as for bad kid's we don't want
	// blast the client (which in turn might blast Secrets Manager or FushionAuth)
	// with a huge number of requests over and over again.
	jwkSet, err := d.jwks.Refresh()
	if err != nil {
		// we didn't refresh, or we did but we failed to parse it
		return nil, errors.Errorf("failed to decode: no matching key_id (kid) header for: %s. err: %w", kid, err)
	}

	key, found := jwkSet.LookupKeyID(kid)
	if found {
		// Found a match, so use this key
		return d.getPublicKey(key)
	}

	return nil, errors.Errorf("failed to decode: no matching key_id (kid) header for: %s", kid)
}

func (d *StandardDecoder) getPublicKey(key jwk.Key) (publicKey, error) {
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
