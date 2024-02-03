package jwt

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

const (
	AuthBearerPrefix   = "Bearer "
	KidHeaderKey       = "kid"
	AlgorithmHeaderKey = "alg"
	SignatureHeaderKey = "sig"
	PerformCoreKid     = "perform-core"

	PublicKeyType = "RSA PUBLIC KEY"

	AccountIDClaim       = "accountId"
	RealUserIDClaim      = "realUserId"
	EffectiveUserIDClaim = "effectiveUserId"
)

// PublicRSAKeyMap "keyid => PublicKey".
type publicRSAKeyMap map[string]*rsa.PublicKey

type JwtDecoder struct {
	publicPEMKey *rsa.PublicKey  // Web Gateway does not provide a kid header
	jwkPEMKeys   publicRSAKeyMap // Optional jwt's signed by other services or Fusion Auth (via JWKS)
}

func NewJwtDecoder(pubKey string, perfCoreKey string, jwkKeys string) (*JwtDecoder, error) {
	decoder := &JwtDecoder{}
	decoder.jwkPEMKeys = make(publicRSAKeyMap)

	// 1. Get default (web-gateway) public key. This is current MANDATORY, but when its included in the JWKS then
	// change the error handling here to be optional
	// REVISIT: Remove this as soon as web-gateway key in JWK
	publicKey, err := decoder.getPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	decoder.publicPEMKey = publicKey

	// 2. Get the the Perform-Core public key
	// REVISIT: Remove this as soon as perf-core key in JWK
	perfKey, err := decoder.getPublicKey(perfCoreKey)
	if err != nil {
		return nil, err
	}
	decoder.jwkPEMKeys["perform-core"] = perfKey

	// 3. Parse JWKs JSON keys if supplied
	rsaPublicKeyMap, err := decoder.parseJWKs(context.Background(), jwkKeys)
	if err != nil {
		return nil, err
	}

	// 4. Upsert into machine keys with "kid" as the key (may overwrite settings.JwtPublicMachineKeys)
	for key, val := range rsaPublicKeyMap {
		decoder.jwkPEMKeys[key] = val
	}

	return decoder, nil
}

func (d *JwtDecoder) Decode(tokenString string) (*StandardClaims, error) {
	payload := &StandardClaims{}

	claims, err := d.decodeClaims(tokenString)
	if err != nil {
		return payload, err
	}

	return NewStandardClaims(claims), nil
}

func (decoder *JwtDecoder) decodeClaims(tokenString string) (jwtgo.MapClaims, error) {
	// sample token string in the form "header.payload.signature"
	// eg. "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// Expiry claim is current OPTIONAL (set jwt.WithExpirationRequired() below if we want to make it mandatory)
	// If the token includes an expiry claim, then it is honoured and the time is checked correctly and will return error if expired
	// If the toekn does not include an expiry clain, then the time is not checked and it will not return an error
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return decoder.useCorrectPublicKey(token)
		},
		// jwt.WithLeeway(10 * time.Second), // add this if we want to add some leeway for clock scew across systems
		// jwt.WithExpirationRequired(),     // add this if we want to enforce that tokens MUST have an expiry
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("missing claims in jwt token")
	}

	return claims, nil
}

func (d *JwtDecoder) useCorrectPublicKey(token *jwt.Token) (*rsa.PublicKey, error) {
	if token == nil {
		return d.publicPEMKey, nil
	}

	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header[AlgorithmHeaderKey])
	}

	kid, found := token.Header[KidHeaderKey]
	if !found {
		return d.publicPEMKey, nil
	}

	key, found := d.jwkPEMKeys[kid.(string)]
	if !found {
		return nil, fmt.Errorf("missing key for kid header: %s", kid)
	}

	return key, nil
}

func (d *JwtDecoder) getPublicKey(key string) (*rsa.PublicKey, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func (decoder *JwtDecoder) parseJWKs(ctx context.Context, jwks string) (publicRSAKeyMap, error) {
	rsaKeys := make(publicRSAKeyMap)

	if jwks == "" {
		// If no jwks json, then returm empty map
		return rsaKeys, nil
	}

	// 1. Parse the jwks JSON string and extra the jwks set
	jwkset, err := jwk.ParseString(jwks)
	if err != nil {
		return rsaKeys, err
	}

	// 2. Enumerate the set
	for it := jwkset.Iterate(ctx); it.Next(ctx); {
		pair := it.Pair()
		key, ok := pair.Value.(jwk.Key)
		if !ok {
			// the jwks Set value isn't valid (for some reason) - just skip it
			continue
		}

		usage := key.KeyUsage()
		if usage != SignatureHeaderKey {
			// we aren't interested if it isn't a "sig"
			continue
		}

		kid := key.KeyID()

		var rsa rsa.PublicKey
		if err := key.Raw(&rsa); err != nil {
			// We only support RSA keys currently so skip if not a RSA public key
			continue
		}

		pubKeyBytes, err := x509.MarshalPKIXPublicKey(&rsa)
		if err != nil {
			return rsaKeys, err
		}

		pubKeyPEM := pem.EncodeToMemory(
			&pem.Block{
				Type:  PublicKeyType,
				Bytes: pubKeyBytes,
			},
		)

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyPEM)
		if err != nil {
			return rsaKeys, err
		}

		// 3. add public key to the map, overwriting if already exists
		rsaKeys[kid] = publicKey
	}

	// 4. return all the valid RSA keys
	return rsaKeys, nil
}
