package jwt

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwkSet(t *testing.T) {
	b, err := os.ReadFile(filepath.Clean(testAuthJwks))
	require.Nil(t, err)

	count := 0
	dispatcher := func() string {
		count++
		return string(b)
	}

	expiresIn := 200 * time.Millisecond
	rotatesIn := 100 * time.Millisecond

	// 1. test constructor
	jwk := newJWKSet(dispatcher, expiresIn, rotatesIn)
	assert.NotNil(t, jwk)

	// 2. test get works ok
	set, err := jwk.Get()
	assert.Nil(t, err)
	assert.NotNil(t, set)
	assert.Equal(t, 1, count)

	// 3. check refresh returns the current set
	set, err = jwk.Refresh()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to refresh jwks as just recently updated")
	assert.NotNil(t, set)
	assert.Equal(t, 1, count)

	// Now wait so that the refresh window is reached
	time.Sleep(100 * time.Millisecond)

	// 4. check refresh returns new set
	set, err = jwk.Refresh()
	assert.Nil(t, err)
	assert.NotNil(t, set)
	assert.Equal(t, 2, count)
}
