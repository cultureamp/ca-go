package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecoderOptions(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	assert.Nil(t, err)
	validJwks := string(b)

	i := 0
	// This will be called each time the cache is refreshed and they we can assert i has been incremented the correct number of times below.
	jwks := func() string {
		i++
		return validJwks
	}

	decoder, err := NewJwtDecoder(jwks, WithDecoderCacheExpiry(100*time.Millisecond, 50*time.Millisecond))
	assert.Nil(t, err)
	assert.NotNil(t, decoder)

	token := "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS01MTIiLCJ0eXAiOiJKV1QifQ.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0.hQLuOe8qZHUgstYe0A0n4-Pww7ovlReyKiDR1Y02ltUnUlgbm9qpp-Ef6YNFuIKdHmS-ynQbDx5pbI36szsggzi80apNpI48cwSXshx82TwuU-_Z4wNBXu7MdPvbA5FdjhxCvRqaqhglsGJ6NofC1bP9awVyyy4j9LGfkVuVEXJQrVpdvEs8Ks-LxlWz7_9Cr7BrZcLuBJnujhe4CbdSudkrfeFl19EY3i1wH9OatGjfjwOSJVqv-ZLnn3QkaZmrQ1xwXTm3MlMUH3KSQjBn8h6vbqosIB5iHDFtqR11mLCgYExGHBpzFjM1d5NEmcTNLV9MtZ_qDZwG0wkgv9O4rXVQ0JfdXypMwhchED2Z45_mc2OiLidtKtDmeoE5g0Daq8YpM0ZpVRbXUFeYIZ1doQKUNsbWNdITmrjVOC3Zn8BecYPu1pC4Hk1y-ViArDzxlCMHA7Bua64BfzVuaJ8pBTEmbqMiZ9VujWcimCOtJ5yfCks_RPAhFYOErcqy3B56fmyYdIN__mKl7VvRDtBSiiPGCq07BUjGywaMoZIULbyXYSV4zs3hX_R4_o4asGiVWCZgn7k4pZzCJo_y2e-Mf85nYoRlyr1MXx7IM4srFQCgO-KTjDWL_TXqpMJU5zDzKyelrMFkc6EaMQ2KP_yBhOrh4UW-Pm7ghusox_-bV1U"
	for j := 0; j < 10; j++ {
		claim, err := decoder.Decode(token)
		assert.Nil(t, err)
		assert.NotNil(t, claim)
		time.Sleep(50 * time.Millisecond)
	}

	assert.Greater(t, i, 1)
}

func TestDecoderParsingOptions(t *testing.T) {
	token := "eyJhbGciOiJSUzUxMiIsImtpZCI6InJzYS01MTIiLCJ0eXAiOiJKV1QifQ.eyJhY2NvdW50SWQiOiJhYmMxMjMiLCJlZmZlY3RpdmVVc2VySWQiOiJ4eXozNDUiLCJyZWFsVXNlcklkIjoieHl6MjM0IiwiaXNzIjoiZW5jb2Rlci1uYW1lIiwic3ViIjoidGVzdCIsImF1ZCI6WyJkZWNvZGVyLW5hbWUiXSwiZXhwIjoyMjExNzk3NTMyLCJuYmYiOjE1ODA2MDg5MjIsImlhdCI6MTU4MDYwODkyMn0.hQLuOe8qZHUgstYe0A0n4-Pww7ovlReyKiDR1Y02ltUnUlgbm9qpp-Ef6YNFuIKdHmS-ynQbDx5pbI36szsggzi80apNpI48cwSXshx82TwuU-_Z4wNBXu7MdPvbA5FdjhxCvRqaqhglsGJ6NofC1bP9awVyyy4j9LGfkVuVEXJQrVpdvEs8Ks-LxlWz7_9Cr7BrZcLuBJnujhe4CbdSudkrfeFl19EY3i1wH9OatGjfjwOSJVqv-ZLnn3QkaZmrQ1xwXTm3MlMUH3KSQjBn8h6vbqosIB5iHDFtqR11mLCgYExGHBpzFjM1d5NEmcTNLV9MtZ_qDZwG0wkgv9O4rXVQ0JfdXypMwhchED2Z45_mc2OiLidtKtDmeoE5g0Daq8YpM0ZpVRbXUFeYIZ1doQKUNsbWNdITmrjVOC3Zn8BecYPu1pC4Hk1y-ViArDzxlCMHA7Bua64BfzVuaJ8pBTEmbqMiZ9VujWcimCOtJ5yfCks_RPAhFYOErcqy3B56fmyYdIN__mKl7VvRDtBSiiPGCq07BUjGywaMoZIULbyXYSV4zs3hX_R4_o4asGiVWCZgn7k4pZzCJo_y2e-Mf85nYoRlyr1MXx7IM4srFQCgO-KTjDWL_TXqpMJU5zDzKyelrMFkc6EaMQ2KP_yBhOrh4UW-Pm7ghusox_-bV1U"

	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	assert.Nil(t, err)
	validJwks := string(b)

	jwks := func() string {
		return validJwks
	}

	decoder, err := NewJwtDecoder(jwks)
	assert.Nil(t, err)
	assert.NotNil(t, decoder)

	// Decode with no MustMatch parsing options
	claim, err := decoder.Decode(token)
	assert.Nil(t, err)
	assert.NotNil(t, claim)

	// Mis-Matched aud claim returns error on Decode
	claim, err = decoder.Decode(token, MustMatchAudience("abc"))
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "token has invalid audience")

	// Matched aud claim returns success
	claim, err = decoder.Decode(token, MustMatchAudience("decoder-name"))
	assert.Nil(t, err)
	assert.NotNil(t, claim)

	// Mis-Matched iss claim returns error on Decode
	claim, err = decoder.Decode(token, MustMatchIssuer("abc"))
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "token has invalid issuer")

	// Matched iss claim returns success
	claim, err = decoder.Decode(token, MustMatchIssuer("encoder-name"))
	assert.Nil(t, err)
	assert.NotNil(t, claim)

	// Mis-Matched sub claim returns error on Decode
	claim, err = decoder.Decode(token, MustMatchSubject("abc"))
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "token has invalid subject")

	// Matched sub claim returns success
	claim, err = decoder.Decode(token, MustMatchSubject("test"))
	assert.Nil(t, err)
	assert.NotNil(t, claim)
}
