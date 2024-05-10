package jwt

import (
	"sync"
	"time"

	"github.com/go-errors/errors"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// jwkSet manages the life-cycle of a jwk.Set().
type jwkSet struct {
	dispatcher     DecoderJwksRetriever // func provided by clients of this library to supply a refreshed JWKS
	expiresWithin  time.Duration
	rotationWindow time.Duration

	mu          sync.RWMutex // mutex to protect race conditions on jwks and jwksAddedAt
	jwks        jwk.Set      // the current JWK Set of jwt public keys
	jwksAddedAt time.Time    // the time when we fetched these keys (so we can refresh)
}

// newJWKSet creates a new jwkSet.
func newJWKSet(dispatcher DecoderJwksRetriever, expiresWithin time.Duration, rotationWindow time.Duration) *jwkSet {
	return &jwkSet{
		dispatcher:     dispatcher,
		expiresWithin:  expiresWithin,
		rotationWindow: rotationWindow,
		jwks:           nil,
	}
}

func (c *jwkSet) Get() (jwk.Set, error) { //nolint:ireturn
	if !c.expired() {
		// we have jwks and it hasn't expired yet, so all good!
		return c.jwks, nil
	}

	jwks, err := c.fetch()
	if err != nil {
		return c.jwks, err
	}

	c.jwksAddedAt = time.Now()
	c.jwks = jwks
	return jwks, nil
}

func (c *jwkSet) Refresh() (jwk.Set, error) { //nolint:ireturn
	if !c.canRefresh() {
		// we can't refresh (ie. get new jwks yet)
		return c.jwks, nil
	}

	jwks, err := c.fetch()
	if err != nil {
		return c.jwks, err
	}

	c.jwksAddedAt = time.Now()
	c.jwks = jwks
	return jwks, nil
}

func (c *jwkSet) expired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.jwks == nil {
		return true
	}

	now := time.Now()
	expiresAt := c.jwksAddedAt.Add(c.expiresWithin)
	return now.After(expiresAt)
}

func (c *jwkSet) canRefresh() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.jwks == nil {
		return true
	}

	freshness := time.Since(c.jwksAddedAt)
	return freshness > c.rotationWindow
}

func (c *jwkSet) fetch() (jwk.Set, error) { //nolint:ireturn
	// Only allow one thread to update the jwks
	c.mu.Lock()
	defer c.mu.Unlock()

	// Call client retriever func
	jwkKeys := c.dispatcher()

	// Parse all new JWKs JSON keys and make sure its valid
	jwkSet, err := c.parse(jwkKeys)
	if err != nil {
		return nil, err
	}

	return jwkSet, nil
}

func (c *jwkSet) parse(jwks string) (jwk.Set, error) { //nolint:ireturn
	if jwks == "" {
		// If no jwks json, then returm empty map
		return nil, errors.Errorf("missing jwks")
	}

	// 1. Parse the jwks JSON string to an iterable set
	return jwk.ParseString(jwks)
}
