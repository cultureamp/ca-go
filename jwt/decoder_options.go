package jwt

import (
	"time"
)

// JwtDecoderOption function signature for added JWT Decoder options.
type JwtDecoderOption func(*JwtDecoder)

// WithDecoderJwksExpiry sets the JwtDecoder JWKs cache expiry time.
// defaultExpiration defaults to 60 minutes.
func WithDecoderJwksExpiry(expiry time.Duration) JwtDecoderOption {
	return func(decoder *JwtDecoder) {
		decoder.expiresWithin = expiry
	}
}

func WithDecoderRotateWindow(rotate time.Duration) JwtDecoderOption {
	return func(decoder *JwtDecoder) {
		decoder.rotationWindow = rotate
	}
}
