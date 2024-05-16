package jwt

import (
	"time"
)

// JwtDecoderOption function signature for added JWT Decoder options.
type JwtDecoderOption func(*JwtDecoder)

// WithDecoderJwksExpiry sets the JwtDecoder JWKs expiry time.Duration
// defaultExpiration defaults to 60 minutes.
func WithDecoderJwksExpiry(expiry time.Duration) JwtDecoderOption {
	return func(decoder *JwtDecoder) {
		decoder.expiresWithin = expiry
	}
}

// WithDecoderRotateWindow sets the JWKS rotation window to an time.Duration.
func WithDecoderRotateWindow(rotate time.Duration) JwtDecoderOption {
	return func(decoder *JwtDecoder) {
		decoder.rotationWindow = rotate
	}
}
