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

	dispatcher := func() string { return string(b) }
	expiresIn := 100 * time.Millisecond
	rotatesIn := 100 * time.Millisecond

	jwk := newJWKSet(dispatcher, expiresIn, rotatesIn)
	assert.NotNil(t, jwk)

	set, err := jwk.get()
	assert.Nil(t, err)
	assert.NotNil(t, set)
}
