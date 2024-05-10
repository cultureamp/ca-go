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
	_, err = jwk.Refresh()
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	time.Sleep(100 * time.Millisecond)

	// 4. check refresh returns new set
	newJwks, err := jwk.Refresh()
	assert.Nil(t, err)
	assert.NotNil(t, newJwks)
	assert.Equal(t, 2, count)
}
