package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetEnvString(t *testing.T) {
	t.Setenv("TEST_STRING", "string")

	val := GetString("should_not_exist_env_var", "fallback")
	assert.Equal(t, "fallback", val)

	val = GetString("TEST_STRING", "fallback")
	assert.Equal(t, "string", val)
}

func Test_GetEnvInt(t *testing.T) {
	t.Setenv("TEST_INT", "123")

	val := GetInt("should_not_exist_env_var", 42)
	assert.Equal(t, 42, val)

	val = GetInt("TEST_INT", 6)
	assert.Equal(t, 123, val)
}

func Test_GetEnvIntFailure(t *testing.T) {
	t.Setenv("TEST_INT", "abc")

	val := GetInt("TEST_INT", 6)
	assert.Equal(t, 6, val)
}

func Test_GetEnvBool(t *testing.T) {
	t.Setenv("TEST_BOOL", "true")

	val := GetBool("should_not_exist_env_var", false)
	assert.False(t, val)

	val = GetBool("TEST_BOOL", false)
	assert.True(t, val)
}

func Test_GetEnvBoolFailure(t *testing.T) {
	t.Setenv("TEST_BOOL", "abc")

	val := GetBool("TEST_BOOL", false)
	assert.False(t, val)
}
