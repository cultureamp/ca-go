package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// StandardClaims represent the standard Culture Amp JWT claims.
type StandardClaims struct {
	AccountId       string    // uuid
	RealUserId      string    // uuid
	EffectiveUserId string    // uuid
	Expiry          time.Time // note: the jwt decoder enforces expiry for you
}

type encoderStandardClaims struct {
	AccountID       string `json:"accountId"`
	EffectiveUserID string `json:"effectiveUserId"`
	RealUserID      string `json:"realUserId"`
	jwt.RegisteredClaims
}

func newStandardClaims(claims jwt.MapClaims) *StandardClaims {
	// todo check for error and do what?
	accountId, _ := getRawClaimString(claims, AccountIDClaim)
	realUserId, _ := getRawClaimString(claims, RealUserIDClaim)
	effectiveUserId, _ := getRawClaimString(claims, EffectiveUserIDClaim)
	expiryTime, _ := getExpirationTime(claims)

	return &StandardClaims{
		AccountId:       accountId,
		RealUserId:      realUserId,
		EffectiveUserId: effectiveUserId,
		Expiry:          expiryTime,
	}
}

func newEncoderClaims(sc *StandardClaims) *encoderStandardClaims {
	claims := &encoderStandardClaims{
		AccountID:       sc.AccountId,
		EffectiveUserID: sc.EffectiveUserId,
		RealUserID:      sc.RealUserId,
	}

	// If Expiry is set, then use it, else set it 1 hour into the future
	if sc.Expiry.IsZero() {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	} else {
		claims.ExpiresAt = jwt.NewNumericDate(sc.Expiry)
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())
	claims.Issuer = "ca-go/jwt"
	claims.Subject = "standard"

	return claims
}

func getExpirationTime(claims jwt.MapClaims) (time.Time, error) {
	// can return nil date with no error
	date, err := claims.GetExpirationTime()
	if err != nil || date == nil {
		return time.Time{}, err
	}

	return date.Time, nil
}

func getRawClaimString(claims jwt.MapClaims, key string) (string, error) {
	val, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("missing %s in jwt token", key)
	}

	return val, nil
}
