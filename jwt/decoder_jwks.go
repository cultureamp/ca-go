package jwt

import (
	"sync"
	"time"

	"github.com/go-errors/errors"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// jwkFetcher manages the life-cycle of a jwk.Set().
type jwkFetcher struct {
	dispatcher     DecoderJwksRetriever // func provided by clients of this library to supply a refreshed JWKS
	expiresWithin  time.Duration
	rotationWindow time.Duration

	mu          sync.RWMutex // mutex to protect race conditions on jwks and jwksAddedAt
	jwks        jwk.Set      // the current JWK Set of jwt public keys
	jwksAddedAt time.Time    // the time when we fetched these keys (so we can refresh)
}

// newJWKSet creates a new jwkSet.
func newJWKSet(dispatcher DecoderJwksRetriever, expiresWithin time.Duration, rotationWindow time.Duration) *jwkFetcher {
	return &jwkFetcher{
		dispatcher:     dispatcher,
		expiresWithin:  expiresWithin,
		rotationWindow: rotationWindow,
		jwks:           nil,
	}
}

func (f *jwkFetcher) Get() (jwk.Set, error) {
	if !f.expired() {
		// we have jwks and it hasn't expired yet, so all good!
		return f.jwks, nil
	}

	jwks, err := f.fetch()
	if err != nil {
		return f.jwks, err
	}

	f.jwksAddedAt = time.Now()
	f.jwks = jwks
	return jwks, nil
}

func (f *jwkFetcher) Refresh() (jwk.Set, error) {
	if !f.canRefresh() {
		// we can't refresh (ie. get new jwks yet) as we just updated recently
		return f.jwks, errors.Errorf("failed to refresh jwks as just recently updated")
	}

	jwks, err := f.fetch()
	if err != nil {
		return f.jwks, err
	}

	f.jwksAddedAt = time.Now()
	f.jwks = jwks
	return jwks, nil
}

func (f *jwkFetcher) expired() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.jwks == nil {
		return true
	}

	now := time.Now()
	expiresAt := f.jwksAddedAt.Add(f.expiresWithin)
	return now.After(expiresAt)
}

func (f *jwkFetcher) canRefresh() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.jwks == nil {
		return true
	}

	freshness := time.Since(f.jwksAddedAt)
	return freshness > f.rotationWindow
}

func (f *jwkFetcher) fetch() (jwk.Set, error) {
	// Only allow one thread to update the jwks
	f.mu.Lock()
	defer f.mu.Unlock()

	// Call client retriever func
	jwkKeys := f.dispatcher()

	// Parse all new JWKs JSON keys and make sure its valid
	jwkSet, err := f.parse(jwkKeys)
	if err != nil {
		return nil, err
	}

	return jwkSet, nil
}

func (f *jwkFetcher) parse(jwks string) (jwk.Set, error) {
	if jwks == "" {
		// If no jwks json, then returm empty map
		return nil, errors.Errorf("missing jwks")
	}

	// 1. Parse the jwks JSON string to an iterable set
	return jwk.ParseString(jwks)
}
