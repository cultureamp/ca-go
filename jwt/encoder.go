package jwt

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

// JwtEncoder can encode a claim to a jwt token string.
type JwtEncoder struct {
	privatePEMKey *rsa.PrivateKey
	kid           string
}

// NewJwtEncoder creates a new JwtEncoder.
func NewJwtEncoder(privateKey string, kid string) (*JwtEncoder, error) {
	privatePEMKey := []byte(privateKey)
	pemKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEMKey)
	if err != nil {
		return nil, err
	}

	return &JwtEncoder{
		privatePEMKey: pemKey,
		kid:           kid,
	}, nil
}

// Encode the Standard Culture Amp Claims in a jwt token string.
func (e *JwtEncoder) Encode(claims *StandardClaims) (string, error) {
	registerClaims := newEncoderClaims(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, registerClaims)
	if e.kid != "" {
		token.Header[KidHeaderKey] = e.kid
	}
	return token.SignedString(e.privatePEMKey)
}
