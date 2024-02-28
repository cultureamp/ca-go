package jwt

import (
	"crypto/elliptic"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ecdsaKey512 = iota
	ecdsaKey384
	ecdsaKey256
	rsaKey
)

type privateKey interface{} // Only ECDSA (perferred) and RSA public keys allowed

// JwtEncoder can encode a claim to a jwt token string.
type JwtEncoder struct {
	privateSigningKey privateKey
	keyType           int
	kid               string
}

// NewJwtEncoder creates a new JwtEncoder.
func NewJwtEncoder(privateKey string, kid string) (*JwtEncoder, error) {
	encoder := &JwtEncoder{}
	privatePEMKey := []byte(privateKey)

	ecdaPrivateKey, err := jwt.ParseECPrivateKeyFromPEM(privatePEMKey)
	if err == nil {
		encoder.privateSigningKey = ecdaPrivateKey
		encoder.kid = kid
		switch ecdaPrivateKey.Curve {
		case elliptic.P256():
			encoder.keyType = ecdsaKey256
		case elliptic.P384():
			encoder.keyType = ecdsaKey384
		default:
			encoder.keyType = ecdsaKey512
		}
		return encoder, nil
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEMKey)
	if err == nil {
		encoder.privateSigningKey = rsaPrivateKey
		encoder.keyType = rsaKey
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
	case ecdsaKey512:
		token = jwt.NewWithClaims(jwt.SigningMethodES512, registerClaims)
	case ecdsaKey384:
		token = jwt.NewWithClaims(jwt.SigningMethodES384, registerClaims)
	case ecdsaKey256:
		token = jwt.NewWithClaims(jwt.SigningMethodES256, registerClaims)
	case rsaKey:
		token = jwt.NewWithClaims(jwt.SigningMethodRS512, registerClaims)
	default:
		return "", fmt.Errorf("Only ECDSA and RSA private keys are supported")
	}

	if e.kid != "" {
		token.Header[kidHeaderKey] = e.kid
	}

	return token.SignedString(e.privateSigningKey)
}
