package jwt_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cultureamp/ca-go/jwt"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

const (
	webGatewayKid = "web-gateway"
)

func Example() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	claim, err := jwt.Decode(token)
	fmt.Printf(
		"The decoded token is '%s %s %s %s %v %s %s' (err='%+v')\n",
		claim.AccountId, claim.RealUserId, claim.EffectiveUserId,
		claim.Issuer, claim.Subject, claim.Audience,
		claim.ExpiresAt.UTC().Format(time.RFC3339),
		err,
	)

	// To create a specific instance of the encoder and decoder you can use the following
	privateKeyBytes, err := os.ReadFile(filepath.Clean("./testKeys/jwt-rsa256-test-webgateway.key"))
	encoder, err := jwt.NewJwtEncoder(func() (string, string) { return string(privateKeyBytes), webGatewayKid })

	token, err = encoder.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	b, err := os.ReadFile(filepath.Clean("./testKeys/development.jwks"))
	decoder, err := jwt.NewJwtDecoder(func() string { return string(b) })

	claim, err = decoder.Decode(token)
	fmt.Printf(
		"The decoded token is '%s %s %s %s %v %s %s' (err='%+v')\n",
		claim.AccountId, claim.RealUserId, claim.EffectiveUserId,
		claim.Issuer, claim.Subject, claim.Audience,
		claim.ExpiresAt.UTC().Format(time.RFC3339),
		err,
	)

	// Output:
	// The encoded token is 'eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0.CH_UIzR_W1275ffAUES0EzsHNRYZyBbrLsKQBbfJ6DpsLW3HAxH5RSjzXL_yCGTrbcHytTYLIZKhN37lC9BZdhkxZtR9bMqqGu4K0zHNtztoC5u1P7kc81FX_dPi9aiR7B4hruSfOFHoWM1A_D_i55qPAJlB0LRFf4nwX9FIWt2IIMwSGUcxfjFYE7MKTlzP3heCYNVzIxLD5g5gcoIyttmltiD_bBvObvExuDsJSlxwrAYvKc2cpIsh1MZ1x16uhG-du2_YdfSK6Ykd6aAvVpq3IGkb99SKS3xUsCV3JkSDRIcWMKzPhEh_huDV4Z3AA3jA4sWvR20WOqzaW3dRAoYIYL7kP92PrXX8m0EtLPAlX471POgNREWqdmxrbdkZcYNHqrmHcAsMRPMXcZ15tH8_-jIDUvGpNbcetgmQRjcpLtyniN_Ag4kGoPhYzGLx6122DEBrYf0Os5TQcRAzAoSF1n_43hsfmuGw00ey3ye5siJle7LN8EHUAXjegrpC7WTFF_eIsOtkuXTJx6OMmuggRvlMaCughYP6IvoIXD7ME0DnzmuvANID9yo-X8DJpMiWbZ2_edCE7dmuqxIZOqJmTolswQs1p0hzFyaX5SrEgcGjHxwTpuCYfaQ7qrbz2D_OQfXbglbk4e8Hm63bGmmz9bKV4KDBVPJO1zOGLtM' (err='<nil>')
	// The decoded token is 'abc123 xyz234 xyz345 encoder-name test [decoder-name] 2040-02-02T12:12:12Z' (err='<nil>')
	// The encoded token is 'eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0.CH_UIzR_W1275ffAUES0EzsHNRYZyBbrLsKQBbfJ6DpsLW3HAxH5RSjzXL_yCGTrbcHytTYLIZKhN37lC9BZdhkxZtR9bMqqGu4K0zHNtztoC5u1P7kc81FX_dPi9aiR7B4hruSfOFHoWM1A_D_i55qPAJlB0LRFf4nwX9FIWt2IIMwSGUcxfjFYE7MKTlzP3heCYNVzIxLD5g5gcoIyttmltiD_bBvObvExuDsJSlxwrAYvKc2cpIsh1MZ1x16uhG-du2_YdfSK6Ykd6aAvVpq3IGkb99SKS3xUsCV3JkSDRIcWMKzPhEh_huDV4Z3AA3jA4sWvR20WOqzaW3dRAoYIYL7kP92PrXX8m0EtLPAlX471POgNREWqdmxrbdkZcYNHqrmHcAsMRPMXcZ15tH8_-jIDUvGpNbcetgmQRjcpLtyniN_Ag4kGoPhYzGLx6122DEBrYf0Os5TQcRAzAoSF1n_43hsfmuGw00ey3ye5siJle7LN8EHUAXjegrpC7WTFF_eIsOtkuXTJx6OMmuggRvlMaCughYP6IvoIXD7ME0DnzmuvANID9yo-X8DJpMiWbZ2_edCE7dmuqxIZOqJmTolswQs1p0hzFyaX5SrEgcGjHxwTpuCYfaQ7qrbz2D_OQfXbglbk4e8Hm63bGmmz9bKV4KDBVPJO1zOGLtM' (err='<nil>')
	// The decoded token is 'abc123 xyz234 xyz345 encoder-name test [decoder-name] 2040-02-02T12:12:12Z' (err='<nil>')
}

func ExampleDefaultJwtEncoder() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	mockEncDec := newMockedEncoderDecoder()
	mockEncDec.On("Encode", mock.Anything).Return("eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0", nil)

	// Overwrite the Default package level encoder and decoder
	oldEncoder := jwt.DefaultJwtEncoder
	jwt.DefaultJwtEncoder = mockEncDec
	defer func() {
		jwt.DefaultJwtEncoder = oldEncoder
	}()

	// Encode this claim with the default "web-gateway" key and add the kid to the token header
	token, err := jwt.Encode(claims)
	fmt.Printf("The encoded token is '%s' (err='%v')\n", token, err)

	// Output:
	// The encoded token is 'eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0' (err='<nil>')
}

func ExampleDefaultJwtDecoder() {
	claims := &jwt.StandardClaims{
		AccountId:       "abc123",
		RealUserId:      "xyz234",
		EffectiveUserId: "xyz345",
		Issuer:          "encoder-name",
		Subject:         "test",
		Audience:        []string{"decoder-name"},
		ExpiresAt:       time.Unix(2211797532, 0), //  2/2/2040
		IssuedAt:        time.Unix(1580608922, 0), // 1/1/2020
		NotBefore:       time.Unix(1580608922, 0), // 1/1/2020
	}

	mockEncDec := newMockedEncoderDecoder()
	mockEncDec.On("Decode", mock.Anything, mock.Anything).Return(claims, nil)

	// Overwrite the Default package level encoder and decoder
	oldDecoder := jwt.DefaultJwtDecoder
	jwt.DefaultJwtDecoder = mockEncDec
	defer func() {
		jwt.DefaultJwtDecoder = oldDecoder
	}()

	// Decode it back again using the key that matches the kid header using the default JWKS JSON keys
	claim, err := jwt.Decode("eyJhbGciOiJSUzUxMiIsImtpZCI6IndlYi1nYXRld2F5IiwidHlwIjoiSldUIn0")
	fmt.Printf(
		"The decoded token is '%s %s %s %s %v %s %s' (err='%+v')\n",
		claim.AccountId, claim.RealUserId, claim.EffectiveUserId,
		claim.Issuer, claim.Subject, claim.Audience,
		claim.ExpiresAt.UTC().Format(time.RFC3339),
		err,
	)

	// Output:
	// The decoded token is 'abc123 xyz234 xyz345 encoder-name test [decoder-name] 2040-02-02T12:12:12Z' (err='<nil>')
}

type mockedEncoderDecoder struct {
	mock.Mock
}

func newMockedEncoderDecoder() *mockedEncoderDecoder {
	return &mockedEncoderDecoder{}
}

func (m *mockedEncoderDecoder) Encode(claims *jwt.StandardClaims) (string, error) {
	args := m.Called(claims)
	output, _ := args.Get(0).(string)
	return output, args.Error(1)
}

func (m *mockedEncoderDecoder) EncodeWithCustomClaims(customClaims gojwt.Claims) (string, error) {
	args := m.Called(customClaims)
	output, _ := args.Get(0).(string)
	return output, args.Error(1)
}

func (m *mockedEncoderDecoder) Decode(tokenString string, options ...jwt.DecoderParserOption) (*jwt.StandardClaims, error) {
	args := m.Called(tokenString, options)
	output, _ := args.Get(0).(*jwt.StandardClaims)
	return output, args.Error(1)
}

func (m *mockedEncoderDecoder) DecodeWithCustomClaims(tokenString string, customClaims gojwt.Claims, options ...jwt.DecoderParserOption) error {
	args := m.Called(tokenString, customClaims, options)
	return args.Error(0)
}
