package jwt_test

import (
	"fmt"

	"github.com/cultureamp/ca-go/jwt"
)

func BasicExamples() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
	}

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	sc, err := jwt.Decode(token)
	fmt.Printf("The decode token is '%v' (err='%+v')\n", sc, err)
}
