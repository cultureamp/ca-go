package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSentrySettings(t *testing.T) {
	settings := newSentrySettings()
	assert.NotNil(t, settings)
}

func TestSentrySettings(t *testing.T) {
	t.Setenv(SentryDsnEnv, "sentry.dsn.com")
	t.Setenv(SentryFlushTimeoutInMsEnv, "1234")

	settings := newSentrySettings()
	assert.Equal(t, "sentry.dsn.com", settings.SentryDSNEnv)
	assert.Equal(t, 1234, settings.SentryFlushInMsEnv)
}
