package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerConfigWithNoEnvVarSet(t *testing.T) {
	t.Setenv("APP", "")
	t.Setenv("APP_VERSION", "")
	t.Setenv("AWS_REGION", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("AWS_ACCOUNT_ID", "")
	t.Setenv("FARM", "")
	t.Setenv("PRODUCT", "")
	t.Setenv("QUIET_MODE", "")

	config := NewLoggerConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "", config.AppName)
	assert.Equal(t, "1.0.0", config.AppVersion)
	assert.Equal(t, "", config.AwsRegion)
	assert.Equal(t, "INFO", config.LogLevel)
	assert.Equal(t, "local", config.AwsAccountID)
	assert.Equal(t, "local", config.Farm)
	assert.Equal(t, "", config.Product)
	assert.Equal(t, true, config.Quiet)
}

func TestNewLoggerConfigWithEnvVarSet(t *testing.T) {
	t.Setenv("APP", "test-app")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("AWS_REGION", "us-west")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("AWS_ACCOUNT_ID", "abc123")
	t.Setenv("FARM", "production")
	t.Setenv("PRODUCT", "performance")
	t.Setenv("QUIET_MODE", "NO")

	config := NewLoggerConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "test-app", config.AppName)
	assert.Equal(t, "1.0.0", config.AppVersion)
	assert.Equal(t, "us-west", config.AwsRegion)
	assert.Equal(t, "DEBUG", config.LogLevel)
	assert.Equal(t, "abc123", config.AwsAccountID)
	assert.Equal(t, "production", config.Farm)
	assert.Equal(t, "performance", config.Product)
	assert.Equal(t, false, config.Quiet)
}
