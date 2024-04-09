package jwt

import (
	"time"
)

// JwtEncoderOption function signature for added JWT Encoder options.
type JwtEncoderOption func(*JwtEncoder)

// WithEncoderCacheExpiry sets the JwtEncoder private key cache expiry time.
// defaultExpiration defaults to 60 minutes.
// cleanupInterval defaults to every 1 minute.
// For no expiry (not recommended for production) use:
// defaultExpiration to NoExpiration (ie. time.Duration = -1).
func WithEncoderCacheExpiry(defaultExpiration, cleanupInterval time.Duration) JwtEncoderOption {
	return func(encoder *JwtEncoder) {
		encoder.defaultExpiration = defaultExpiration
		encoder.cleanupInterval = cleanupInterval
	}
}
