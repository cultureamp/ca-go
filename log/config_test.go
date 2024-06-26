package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerConfigWithNoEnvVarSet(t *testing.T) {
	t.Setenv(AppNameEnv, "")
	t.Setenv(AppVerEnv, "")
	t.Setenv(AwsRegionEnv, "")
	t.Setenv(LogLevelEnv, "")
	t.Setenv(AwsAccountIDEnv, "")
	t.Setenv(AppFarmEnv, "")
	t.Setenv(ProductEnv, "")
	t.Setenv(LogQuietModeEnv, "")

	config, err := NewLoggerConfig()
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "<unknown>", config.AppName)
	assert.Equal(t, "1.0.0", config.AppVersion)
	assert.Equal(t, "dev", config.AwsRegion)
	assert.Equal(t, "INFO", config.LogLevel)
	assert.Equal(t, "development", config.AwsAccountID)
	assert.Equal(t, "local", config.Farm)
	assert.Equal(t, "<unknown>", config.Product)
	assert.Equal(t, false, config.Quiet)
}

func TestNewLoggerConfigWithEnvVarSet(t *testing.T) {
	t.Setenv(AppNameEnv, "test-app")
	t.Setenv(AppVerEnv, "2.1.2")
	t.Setenv(AwsRegionEnv, "us-west")
	t.Setenv(LogLevelEnv, "DEBUG")
	t.Setenv(AwsAccountIDEnv, "abc123")
	t.Setenv(AppFarmEnv, "production")
	t.Setenv(ProductEnv, "performance")
	t.Setenv(LogQuietModeEnv, "true")

	config, err := NewLoggerConfig()
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test-app", config.AppName)
	assert.Equal(t, "2.1.2", config.AppVersion)
	assert.Equal(t, "us-west", config.AwsRegion)
	assert.Equal(t, "DEBUG", config.LogLevel)
	assert.Equal(t, "abc123", config.AwsAccountID)
	assert.Equal(t, "production", config.Farm)
	assert.Equal(t, "performance", config.Product)
	assert.Equal(t, true, config.Quiet)
}
