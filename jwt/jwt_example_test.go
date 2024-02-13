package jwt_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cultureamp/ca-go/jwt"
)

func Example() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	sc, err := jwt.Decode(token)
	fmt.Printf("The decoded token is '%v' (err='%+v')\n", sc, err)

	// To create a specific instance of the encoder and decoder you can use the following
	privateKeyBytes, err := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
	encoder, err := jwt.NewJwtEncoder(string(privateKeyBytes), "web-gateway")

	token, err = encoder.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	b, err := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
	decoder, err := jwt.NewJwtDecoderWithDefaultKid(string(b), "web-gateway")

	claim, err := decoder.Decode(token)
	fmt.Printf("The decoded token is '%v' (err='%+v')\n", claim, err)

	// Output:
	// The encoded token is 'eyJhbGciOiJSUzI1NiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.IYPu_PGUO7lpy_wTSObA4S-n9HQUwPf2kTG2AzvSFUwfz994SHZOazYL7CyiRqqhIndIt5R4CQ3cXY7_Lok_wgBQ-U4FAciJw0Fx9tawJIEqwVeL10P4w0h5OIU21E7jeNmlcLOO57QN-ip7hc_--zyAFVKV5qjlbemuHWWpeUGu62SsdHr4J33O6hR8ubTyfXVF7wxKhNM4hCdM7PNanP9OOyAgEWxhwutURiA1nJsATwDf6QKNceGpqkb5A31PvFdfPHoktY4u6e4feBt2KjYJ1xy9opDlllFOEIwTw4nuksQk4q3437bGtfoQkC_CTGO83YTX5GHs70rxu_AubBxCazqSxqMwagiekkpgKZd6d0g7u5F5K8QImRJsore3oHNDAuVg7pbZmH9sApFN_bJhonOkECoPeeF5oYLSLHOXjN7CakvAsmCW01_ENPVXXO2E1yObzwmsY28_Ox5r_jC6XugGdXVfco6l1Oqbxb0ogG6BbOngYEZwVMbEO5qsBnUtBfr0nNUjFKIYCYXdpoeT_bxlt8GI4H2cMAb6FGa_XIEd60fJGazgAk9axA61xHEnqxgUyZv5PEL908zPBRvcNGpQuMsDpGOXTOQ_fgJO1IRBx4VwWcobzKbOyRNarTNwQZH0OY13HMMnFoiPjk8U0fWkJdj1ujobTQYYtz0' (err='<nil>')
	// The decoded token is '&{abc123 xyz234 xyz345 2040-02-02 23:12:12 +1100 AEDT 2020-02-02 13:02:02 +1100 AEDT 2020-02-02 13:02:02 +1100 AEDT ca-go/jwt standard}' (err='<nil>')
	// The encoded token is 'eyJhbGciOiJSUzI1NiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiY2EtZ28vand0Iiwic3ViIjoic3RhbmRhcmQiLCJleHAiOjIyMTE3OTc1MzIsIm5iZiI6MTU4MDYwODkyMiwiaWF0IjoxNTgwNjA4OTIyfQ.IYPu_PGUO7lpy_wTSObA4S-n9HQUwPf2kTG2AzvSFUwfz994SHZOazYL7CyiRqqhIndIt5R4CQ3cXY7_Lok_wgBQ-U4FAciJw0Fx9tawJIEqwVeL10P4w0h5OIU21E7jeNmlcLOO57QN-ip7hc_--zyAFVKV5qjlbemuHWWpeUGu62SsdHr4J33O6hR8ubTyfXVF7wxKhNM4hCdM7PNanP9OOyAgEWxhwutURiA1nJsATwDf6QKNceGpqkb5A31PvFdfPHoktY4u6e4feBt2KjYJ1xy9opDlllFOEIwTw4nuksQk4q3437bGtfoQkC_CTGO83YTX5GHs70rxu_AubBxCazqSxqMwagiekkpgKZd6d0g7u5F5K8QImRJsore3oHNDAuVg7pbZmH9sApFN_bJhonOkECoPeeF5oYLSLHOXjN7CakvAsmCW01_ENPVXXO2E1yObzwmsY28_Ox5r_jC6XugGdXVfco6l1Oqbxb0ogG6BbOngYEZwVMbEO5qsBnUtBfr0nNUjFKIYCYXdpoeT_bxlt8GI4H2cMAb6FGa_XIEd60fJGazgAk9axA61xHEnqxgUyZv5PEL908zPBRvcNGpQuMsDpGOXTOQ_fgJO1IRBx4VwWcobzKbOyRNarTNwQZH0OY13HMMnFoiPjk8U0fWkJdj1ujobTQYYtz0' (err='<nil>')
	// The decoded token is '&{abc123 xyz234 xyz345 2040-02-02 23:12:12 +1100 AEDT 2020-02-02 13:02:02 +1100 AEDT 2020-02-02 13:02:02 +1100 AEDT ca-go/jwt standard}' (err='<nil>')
}
