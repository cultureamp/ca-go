package jwt

import (
	"time"
)

type JwtDecoderOption func(*JwtDecoder)

func WithDecoderCacheExpiry(defaultExpiration, cleanupInterval time.Duration) JwtDecoderOption {
	return func(decoder *JwtDecoder) {
		decoder.defaultExpiration = defaultExpiration
		decoder.cleanupInterval = cleanupInterval
	}
}
