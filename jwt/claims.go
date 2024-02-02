package jwt

import (
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

type StandardClaims struct {
	AccountId       string    // uuid
	RealUserId      string    // uuid
	EffectiveUserId string    // uuid
	Expiry          time.Time // note: the jwt decoder enforces expiry for you
}

func NewStandardClaims(claims *jwtgo.MapClaims) *StandardClaims {
	// todo copy claims
	return &StandardClaims{}
} 

func (sc *StandardClaims) toRegisteredClaims() *jwtgo.RegisteredClaims {
	// todo copy claims
	return &jwtgo.RegisteredClaims{}
}
