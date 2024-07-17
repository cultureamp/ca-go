package jwt

import (
	"time"
)

// DecoderOption function signature for adding JWT Decoder options.
type DecoderOption func(*StandardDecoder)

// WithDecoderJwksExpiry sets the JwtDecoder JWKs expiry time.Duration
// defaultExpiration defaults to 60 minutes.
func WithDecoderJwksExpiry(expiry time.Duration) DecoderOption {
	return func(decoder *StandardDecoder) {
		decoder.expiresWithin = expiry
	}
}

// WithDecoderRotateWindow sets the JWKS rotation window to an time.Duration.
func WithDecoderRotateWindow(rotate time.Duration) DecoderOption {
	return func(decoder *StandardDecoder) {
		decoder.rotationWindow = rotate
	}
}

// DecoderParserOption function signature for adding JWT Decoder Parsing options.
type DecoderParserOption func(*decoderParser)

type decoderParser struct {
	expectedAud string
	expectedIss string
	expectedSub string
}

func newDecoderParser() *decoderParser {
	return &decoderParser{}
}

// MustMatchAudience configures the jwt parser to require the specified audience in
// the `aud` claim. Validation will fail if the audience is not listed in the
// token or the `aud` claim is missing.
func MustMatchAudience(aud string) DecoderParserOption {
	return func(p *decoderParser) {
		p.expectedAud = aud
	}
}

// MustMatchIssuer configures the jwt parser to require the specified issuer in the
// `iss` claim. Validation will fail if a different issuer is specified in the
// token or the `iss` claim is missing.
func MustMatchIssuer(iss string) DecoderParserOption {
	return func(p *decoderParser) {
		p.expectedIss = iss
	}
}

// MustMatchSubject configures the jwt parser to require the specified subject in the
// `sub` claim. Validation will fail if a different subject is specified in the
// token or the `sub` claim is missing.
func MustMatchSubject(sub string) DecoderParserOption {
	return func(p *decoderParser) {
		p.expectedSub = sub
	}
}
