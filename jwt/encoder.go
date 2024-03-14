package jwt

import (
	"crypto/elliptic"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
)

const (
	privCacheKey = "encoder_priv_key"

	invalidKey = iota
	rsaKey
	ecdsaKey512
	ecdsaKey384
	ecdsaKey256
)

type (
	privateKey          interface{} // Only ECDSA (perferred) and RSA public keys allowed
	EncoderKeyRetriever func() (string, string)
)

type encoderPrivateKey struct {
	privateSigningKey privateKey
	keyType           int
	kid               string
}

// JwtEncoder can encode a claim to a jwt token string.
type JwtEncoder struct {
	fetchPrivateKey EncoderKeyRetriever

	mu sync.Mutex

	// memory cache holding the privateKey
	cache             *cache.Cache
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

// NewJwtEncoder creates a new JwtEncoder.
func NewJwtEncoder(fetchPrivateKey EncoderKeyRetriever, options ...JwtEncoderOption) (*JwtEncoder, error) {
	encoder := &JwtEncoder{
		fetchPrivateKey:   fetchPrivateKey,
		defaultExpiration: 60 * time.Minute,
		cleanupInterval:   10 * time.Minute,
	}

	// Loop through our Encoder options and apply them
	for _, option := range options {
		option(encoder)
	}

	encoder.cache = cache.New(encoder.defaultExpiration, encoder.cleanupInterval)

	// call the fetchPrivateKey func to make sure the private key is valid
	_, err := encoder.getPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return encoder, nil
}

// Encode the Standard Culture Amp Claims in a jwt token string.
func (e *JwtEncoder) Encode(claims *StandardClaims) (string, error) {
	var token *jwt.Token

	registerClaims := newEncoderClaims(claims)

	// Will check cache and re-fetch if expired
	encodingKey, err := e.getPrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to load private key: %w", err)
	}

	switch encodingKey.keyType {
	case ecdsaKey512:
		token = jwt.NewWithClaims(jwt.SigningMethodES512, registerClaims)
	case ecdsaKey384:
		token = jwt.NewWithClaims(jwt.SigningMethodES384, registerClaims)
	case ecdsaKey256:
		token = jwt.NewWithClaims(jwt.SigningMethodES256, registerClaims)
	case rsaKey:
		token = jwt.NewWithClaims(jwt.SigningMethodRS512, registerClaims)
	default:
		return "", fmt.Errorf("Only ECDSA and RSA private keys are supported")
	}

	if encodingKey.kid != "" {
		token.Header[kidHeaderKey] = encodingKey.kid
	}

	return token.SignedString(encodingKey.privateSigningKey)
}

func (e *JwtEncoder) getPrivateKey() (*encoderPrivateKey, error) {
	// First chech cache, if its there then great, use it!
	obj, found := e.cache.Get(privCacheKey)
	if found {
		key, ok := obj.(*encoderPrivateKey)
		if !ok {
			return nil, fmt.Errorf("internal error: cache key does not point to private key")
		}

		return key, nil
	}

	// The cache has expired the key

	// Only allow one thread to fetch, parse and update the cache
	e.mu.Lock()
	defer e.mu.Unlock()

	// Time to fetch new ones
	privateKey, kid := e.fetchPrivateKey()

	encodingKey, err := e.parsePrivateKey(privateKey, kid)
	if err != nil {
		return nil, err
	}

	// Add back into the cache
	err = e.cache.Add(privCacheKey, encodingKey, cache.DefaultExpiration)
	return encodingKey, err
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

	return nil, fmt.Errorf("invalid private key: only ECDSA and RSA private keys are supported")
}
