package jwt

import (
	"crypto/elliptic"
	"sync"
	"time"

	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
)

const (
	privCacheKey                  = "encoder_priv_key"
	defaultEncoderExpiration      = 60 * time.Minute
	defaultEncoderCleanupInterval = 1 * time.Minute

	invalidKey = iota
	rsaKey
	ecdsaKey512
	ecdsaKey384
	ecdsaKey256
)

// Only ECDSA (perferred) and RSA public keys allowed.
type privateKey interface{}

// EncoderKeyRetriever defines the function signature required to retrieve private PEM key.
type EncoderKeyRetriever func() (string, string)

type encoderPrivateKey struct {
	privateSigningKey privateKey
	keyType           int
	kid               string
}

// JwtEncoder can encode a claim to a jwt token string.
type JwtEncoder struct {
	fetchPrivateKey   EncoderKeyRetriever // func provided by clients of this library to supply a refreshed private key and kid
	mu                sync.Mutex          // mutex to protect cache.Get/Set race condition
	cache             *cache.Cache        // memory cache holding the encoderPrivateKey struct
	defaultExpiration time.Duration       // default is 60 minutes
	cleanupInterval   time.Duration       // default is every 1 minute
}

// NewJwtEncoder creates a new JwtEncoder.
func NewJwtEncoder(fetchPrivateKey EncoderKeyRetriever, options ...JwtEncoderOption) (*JwtEncoder, error) {
	encoder := &JwtEncoder{
		fetchPrivateKey:   fetchPrivateKey,
		defaultExpiration: defaultEncoderExpiration,
		cleanupInterval:   defaultEncoderCleanupInterval,
	}

	// Loop through our Encoder options and apply them
	for _, option := range options {
		option(encoder)
	}

	encoder.cache = cache.New(encoder.defaultExpiration, encoder.cleanupInterval)

	// call the fetchPrivateKey func to make sure the private key is valid
	_, err := encoder.loadPrivateKey()
	if err != nil {
		return nil, errors.Errorf("failed to load private key: %w", err)
	}

	return encoder, nil
}

// Encode the Standard Culture Amp Claims in a jwt token string.
func (e *JwtEncoder) Encode(claims *StandardClaims) (string, error) {
	registerClaims := newEncoderClaims(claims)

	return e.EncodeWithCustomClaims(registerClaims)
}

// EncodeWithCustomClaims encodes the Custom Claims in a jwt token string.
func (e *JwtEncoder) EncodeWithCustomClaims(customClaims jwt.Claims) (string, error) {
	var token *jwt.Token

	// Will check cache and re-fetch if expired
	encodingKey, err := e.loadPrivateKey()
	if err != nil {
		return "", errors.Errorf("failed to load private key: %w", err)
	}

	switch encodingKey.keyType {
	case ecdsaKey512:
		token = jwt.NewWithClaims(jwt.SigningMethodES512, customClaims)
	case ecdsaKey384:
		token = jwt.NewWithClaims(jwt.SigningMethodES384, customClaims)
	case ecdsaKey256:
		token = jwt.NewWithClaims(jwt.SigningMethodES256, customClaims)
	case rsaKey:
		token = jwt.NewWithClaims(jwt.SigningMethodRS512, customClaims)
	default:
		return "", errors.Errorf("Only ECDSA and RSA private keys are supported")
	}

	if encodingKey.kid != "" {
		token.Header[kidHeaderKey] = encodingKey.kid
	}

	return token.SignedString(encodingKey.privateSigningKey)
}

func (e *JwtEncoder) loadPrivateKey() (*encoderPrivateKey, error) {
	// First chech cache, if its there then great, use it!
	if key, ok := e.getCachedPrivateKey(); ok {
		return key, nil
	}

	// Only allow one thread to refetch, parse and update the cache
	e.mu.Lock()
	defer e.mu.Unlock()

	// check the cache again in case another go routine just updated it
	if key, ok := e.getCachedPrivateKey(); ok {
		return key, nil
	}

	// Call client retriever func
	privateKey, kid := e.fetchPrivateKey()

	// check its valid by parsing the PEM key
	encodingKey, err := e.parsePrivateKey(privateKey, kid)
	if err != nil {
		return nil, err
	}

	// Add back into the cache
	err = e.cache.Add(privCacheKey, encodingKey, cache.DefaultExpiration)
	return encodingKey, err
}

func (e *JwtEncoder) getCachedPrivateKey() (*encoderPrivateKey, bool) {
	// First chech cache, if its there then great, use it!
	obj, found := e.cache.Get(privCacheKey)
	if !found {
		return nil, false
	}

	key, ok := obj.(*encoderPrivateKey)
	return key, ok
}

func (e *JwtEncoder) parsePrivateKey(privKey string, kid string) (*encoderPrivateKey, error) {
	privatePEMKey := []byte(privKey)

	ecdaPrivateKey, err := jwt.ParseECPrivateKeyFromPEM(privatePEMKey)
	if err == nil {
		edcdaKey := &encoderPrivateKey{
			privateSigningKey: ecdaPrivateKey,
			kid:               kid,
		}
		switch ecdaPrivateKey.Curve {
		case elliptic.P256():
			edcdaKey.keyType = ecdsaKey256
		case elliptic.P384():
			edcdaKey.keyType = ecdsaKey384
		default:
			edcdaKey.keyType = ecdsaKey512
		}

		return edcdaKey, nil
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEMKey)
	if err == nil {
		return &encoderPrivateKey{
			privateSigningKey: rsaPrivateKey,
			keyType:           rsaKey,
			kid:               kid,
		}, nil
	}
	// add other key types in the future

	return nil, errors.Errorf("invalid private key: only ECDSA and RSA private keys are supported")
}
