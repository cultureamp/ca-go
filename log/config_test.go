package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerConfigWithNoEnv(t *testing.T) {
	unsetEnvironmentVariables()

	config := newLoggerConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "", config.AppName)
	assert.Equal(t, "1.0.0", config.AppVersion)
	assert.Equal(t, "", config.AwsRegion)
	assert.Equal(t, "INFO", config.LogLevel)
	assert.Equal(t, "local", config.AwsAccountID)
	assert.Equal(t, "local", config.Farm)
	assert.Equal(t, "", config.Product)
}

func TestNewLoggerConfigWithEnv(t *testing.T) {
	unsetEnvironmentVariables()

	t.Setenv("APP", "test-app")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("AWS_REGION", "us-west")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("AWS_ACCOUNT_ID", "abc123")
	t.Setenv("FARM", "production")
	t.Setenv("PRODUCT", "performance")

	config := newLoggerConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "test-app", config.AppName)
	assert.Equal(t, "1.0.0", config.AppVersion)
	assert.Equal(t, "us-west", config.AwsRegion)
	assert.Equal(t, "DEBUG", config.LogLevel)
	assert.Equal(t, "abc123", config.AwsAccountID)
	assert.Equal(t, "production", config.Farm)
	assert.Equal(t, "performance", config.Product)
}

func unsetEnvironmentVariables() {
	os.Unsetenv("APP")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("AWS_ACCOUNT_ID")
	os.Unsetenv("FARM")
	os.Unsetenv("PRODUCT")
}
