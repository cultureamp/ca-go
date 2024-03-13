package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogSettings(t *testing.T) {
	settings := newLogSettings()
	assert.NotNil(t, settings)
}

func TestLogSettings(t *testing.T) {
	t.Setenv(LogLevelEnv, "WARN")

	settings := newLogSettings()
	assert.Equal(t, "WARN", settings.LSLogLevel)
}
