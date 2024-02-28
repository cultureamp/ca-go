package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type privateKey interface{} // Only ECDSA (perferred) and RSA public keys allowed

// JwtEncoder can encode a claim to a jwt token string.
type JwtEncoder struct {
	privatePEMKey privateKey
	keyType       int
	kid           string
}

// NewJwtEncoder creates a new JwtEncoder.
func NewJwtEncoder(privateKey string, kid string) (*JwtEncoder, error) {
	encoder := &JwtEncoder{}
	privatePEMKey := []byte(privateKey)

	pemECDSAKey, err := jwt.ParseECPrivateKeyFromPEM(privatePEMKey)
	if err == nil {
		encoder.privatePEMKey = pemECDSAKey
		encoder.keyType = ECDSA512
		encoder.kid = kid
		return encoder, nil
	}

	pemRSAKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEMKey)
	if err == nil {
		encoder.privatePEMKey = pemRSAKey
		encoder.keyType = RSA512
		encoder.kid = kid
		return encoder, nil
	}
	// add other key types in the future

	return nil, fmt.Errorf("invalid private key: only ECDSA and RSA private keys are supported")
}

// Encode the Standard Culture Amp Claims in a jwt token string.
func (e *JwtEncoder) Encode(claims *StandardClaims) (string, error) {
	var token *jwt.Token

	registerClaims := newEncoderClaims(claims)

	switch e.keyType {
	case ECDSA512:
		token = jwt.NewWithClaims(jwt.SigningMethodES512, registerClaims)
	case RSA512:
		token = jwt.NewWithClaims(jwt.SigningMethodRS512, registerClaims)
	default:
		return "", fmt.Errorf("Only ECDSA and RSA private keys are supported")
	}

	if e.kid != "" {
		token.Header[kidHeaderKey] = e.kid
	}

	return token.SignedString(e.privatePEMKey)
}
