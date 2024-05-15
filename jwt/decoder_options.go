package jwt

import (
	"time"
)

// JwtDecoderOption function signature for added JWT Decoder options.
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

type DecoderParserOption func(*decoderParser)

type decoderParser struct {
	expectedAud string
	expectedIss string
	expectedSub string
}

func newDecoderParser() *decoderParser {
	return &decoderParser{}
}

// MustMatchAudience configures the validator to require the specified audience in
// the `aud` claim. Validation will fail if the audience is not listed in the
// token or the `aud` claim is missing.
func MustMatchAudience(aud string) DecoderParserOption {
	return func(p *decoderParser) {
		p.expectedAud = aud
	}
}

// MustMatchIssuer configures the validator to require the specified issuer in the
// `iss` claim. Validation will fail if a different issuer is specified in the
// token or the `iss` claim is missing.
func MustMatchIssuer(iss string) DecoderParserOption {
	return func(p *decoderParser) {
		p.expectedIss = iss
	}
}

// MustMatchSubject configures the validator to require the specified subject in the
// `sub` claim. Validation will fail if a different subject is specified in the
// token or the `sub` claim is missing.
func MustMatchSubject(sub string) DecoderParserOption {
	return func(p *decoderParser) {
		p.expectedSub = sub
	}
}
