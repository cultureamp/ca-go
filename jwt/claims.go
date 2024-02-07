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

	// the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
	ExpiresAt time.Time // default on Encode is +1 hour from now
	// the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
	IssuedAt time.Time // default on Encode is "now"
	// the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
	NotBefore time.Time // default on Encode is "now"
	// the `iss` (Issuer) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
	Issuer string // default on Encode is "ca-jwt-go"
	// the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	Subject string // default on Encode is "standard"
}

type encoderStandardClaims struct {
	AccountID       string `json:"accountId"`
	EffectiveUserID string `json:"effectiveUserId"`
	RealUserID      string `json:"realUserId"`
	jwt.RegisteredClaims
}

func newStandardClaims(claims jwt.MapClaims) *StandardClaims {
	// todo check for error and do what?
	accountId := getRawClaimString(claims, AccountIDClaim)
	realUserId := getRawClaimString(claims, RealUserIDClaim)
	effectiveUserId := getRawClaimString(claims, EffectiveUserIDClaim)
	expiryTime := getExpirationTime(claims)
	notBeforeTime := getNotBeforeTime(claims)
	issuedAtTime := getIssuedAtTime(claims)
	issuer := getIssuer(claims)
	subject := getSubject(claims)

	return &StandardClaims{
		AccountId:       accountId,
		RealUserId:      realUserId,
		EffectiveUserId: effectiveUserId,
		ExpiresAt:       expiryTime,
		NotBefore:       notBeforeTime,
		IssuedAt:        issuedAtTime,
		Issuer:          issuer,
		Subject:         subject,
	}
}

func newEncoderClaims(sc *StandardClaims) *encoderStandardClaims {
	claims := &encoderStandardClaims{
		AccountID:       sc.AccountId,
		EffectiveUserID: sc.EffectiveUserId,
		RealUserID:      sc.RealUserId,
	}

	// If Expiry is set, then use it, else set it 1 hour into the future
	if sc.ExpiresAt.IsZero() {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	} else {
		claims.ExpiresAt = jwt.NewNumericDate(sc.ExpiresAt)
	}

	if sc.IssuedAt.IsZero() {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
	} else {
		claims.IssuedAt = jwt.NewNumericDate(sc.IssuedAt)
	}

	if sc.NotBefore.IsZero() {
		claims.NotBefore = jwt.NewNumericDate(time.Now())
	} else {
		claims.NotBefore = jwt.NewNumericDate(sc.NotBefore)
	}

	if sc.Issuer == "" {
		claims.Issuer = "ca-go/jwt"
	}

	if sc.Subject == "" {
		claims.Subject = "standard"
	}

	return claims
}

func getExpirationTime(claims jwt.MapClaims) time.Time {
	// can return nil date with no error
	date, err := claims.GetExpirationTime()
	if err != nil || date == nil {
		return time.Time{}
	}

	return date.Time
}

func getNotBeforeTime(claims jwt.MapClaims) time.Time {
	// can return nil date with no error
	date, err := claims.GetNotBefore()
	if err != nil || date == nil {
		return time.Time{}
	}

	return date.Time
}

func getIssuedAtTime(claims jwt.MapClaims) time.Time {
	// can return nil date with no error
	date, err := claims.GetIssuedAt()
	if err != nil || date == nil {
		return time.Time{}
	}

	return date.Time
}

func getIssuer(claims jwt.MapClaims) string {
	issuer, err := claims.GetIssuer()
	if err != nil {
		return ""
	}

	return issuer
}

func getSubject(claims jwt.MapClaims) string {
	sub, err := claims.GetSubject()
	if err != nil {
		return ""
	}

	return sub
}

func getRawClaimString(claims jwt.MapClaims, key string) string {
	val, ok := claims[key].(string)
	if !ok {
		return ""
	}

	return val
}
