package jwt

import (
	"time"
)

type JwtDecoderOption func(*JwtDecoder)

// WithDecoderCacheExpiry sets the JwtDecoder JWKs cache expiry time.
// defaultExpiration defaults to 60 minutes.
// cleanupInterval defaults to every 1 minute.
// For no expiry (not recommended for production) use:
// defaultExpiration to NoExpiration (ie. time.Duration = -1).
func WithDecoderCacheExpiry(defaultExpiration, cleanupInterval time.Duration) JwtDecoderOption {
	return func(decoder *JwtDecoder) {
		decoder.defaultExpiration = defaultExpiration
		decoder.cleanupInterval = cleanupInterval
	}
}
