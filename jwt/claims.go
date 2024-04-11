package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// StandardClaims represent the standard Culture Amp JWT claims.
type StandardClaims struct {
	AccountId       string // uuid
	RealUserId      string // uuid
	EffectiveUserId string // uuid

	// Optional claims

	// the `iss` (Issuer) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
	Issuer string
	// the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	Subject string
	// the `aud` (Audience) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
	Audience []string
	// the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
	ExpiresAt time.Time // default on Encode is +1 hour from now
	// the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
	NotBefore time.Time // default on Encode is "now"
	// the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
	IssuedAt time.Time // default on Encode is "now"
}

type encoderStandardClaims struct {
	AccountID       string `json:"accountId"`
	EffectiveUserID string `json:"effectiveUserId"`
	RealUserID      string `json:"realUserId"`
	jwt.RegisteredClaims
}

func newStandardClaims(claims jwt.MapClaims) *StandardClaims {
	std := &StandardClaims{}

	std.AccountId = std.getCustomString(claims, accountIDClaim)
	std.RealUserId = std.getCustomString(claims, realUserIDClaim)
	std.EffectiveUserId = std.getCustomString(claims, effectiveUserIDClaim)

	std.Issuer = std.getString(claims.GetIssuer)
	std.Subject = std.getString(claims.GetSubject)
	std.Audience = std.getStringList(claims.GetAudience)
	std.ExpiresAt = std.getTime(claims.GetExpirationTime)
	std.NotBefore = std.getTime(claims.GetNotBefore)
	std.IssuedAt = std.getTime(claims.GetIssuedAt)

	return std
}

func (sc *StandardClaims) getTime(f func() (*jwt.NumericDate, error)) time.Time {
	// can return nil date with no error
	date, err := f()
	if err != nil || date == nil {
		return time.Time{}
	}

	return date.Time
}

func (sc *StandardClaims) getString(f func() (string, error)) string {
	s, err := f()
	if err != nil {
		return ""
	}

	return s
}

func (sc *StandardClaims) getStringList(f func() (jwt.ClaimStrings, error)) []string {
	s, err := f()
	if err != nil {
		return []string{""}
	}

	return s
}

func (sc *StandardClaims) getCustomString(claims jwt.MapClaims, key string) string {
	val, ok := claims[key].(string)
	if !ok {
		return ""
	}

	return val
}

func newEncoderClaims(sc *StandardClaims) *encoderStandardClaims {
	claims := &encoderStandardClaims{
		AccountID:       sc.AccountId,
		EffectiveUserID: sc.EffectiveUserId,
		RealUserID:      sc.RealUserId,
	}

	claims.Issuer = sc.Issuer
	claims.Subject = sc.Subject
	claims.Audience = sc.Audience
	now := time.Now()
	claims.ExpiresAt = claims.correctTime(sc.ExpiresAt, now.Add(10*time.Minute))
	claims.NotBefore = claims.correctTime(sc.NotBefore, now)
	claims.IssuedAt = claims.correctTime(sc.IssuedAt, now)

	return claims
}

func (esc *encoderStandardClaims) correctTime(t time.Time, def time.Time) *jwt.NumericDate {
	if t.IsZero() {
		return jwt.NewNumericDate(def)
	}

	return jwt.NewNumericDate(t)
}
