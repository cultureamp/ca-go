package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_redact(t *testing.T) {
	s := ""
	r := redact(s)
	assert.Equal(t, "******", r)

	s = "1234"
	r = redact(s)
	assert.Equal(t, "******", r)

	s = "123456"
	r = redact(s)
	assert.Equal(t, "******", r)

	s = "1234567"
	r = redact(s)
	assert.Equal(t, "******7", r)

	s = "12345678"
	r = redact(s)
	assert.Equal(t, "******78", r)

	s = "123456789"
	r = redact(s)
	assert.Equal(t, "******789", r)

	s = "1234567890"
	r = redact(s)
	assert.Equal(t, "******7890", r)

	s = "12345678901"
	r = redact(s)
	assert.Equal(t, "******8901", r)

	s = "123456789012"
	r = redact(s)
	assert.Equal(t, "******9012", r)

	s = "1234567890123"
	r = redact(s)
	assert.Equal(t, "******0123", r)
}
