package jwt

import (
	"time"
)

// EncoderOption function signature for added JWT Encoder options.
type EncoderOption func(*StandardEncoder)

// WithEncoderCacheExpiry sets the JwtEncoder private key cache expiry time.
// defaultExpiration defaults to 60 minutes.
// cleanupInterval defaults to every 1 minute.
// For no expiry (not recommended for production) use:
// defaultExpiration to NoExpiration (ie. time.Duration = -1).
func WithEncoderCacheExpiry(defaultExpiration, cleanupInterval time.Duration) EncoderOption {
	return func(encoder *StandardEncoder) {
		encoder.defaultExpiration = defaultExpiration
		encoder.cleanupInterval = cleanupInterval
	}
}
