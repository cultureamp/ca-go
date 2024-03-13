package jwt

import (
	"time"
)

type JwtEncoderOption func(*JwtEncoder)

func WithEncoderCacheExpiry(defaultExpiration, cleanupInterval time.Duration) JwtEncoderOption {
	return func(encoder *JwtEncoder) {
		encoder.defaultExpiration = defaultExpiration
		encoder.cleanupInterval = cleanupInterval
	}
}
