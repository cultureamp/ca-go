package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettings(t *testing.T) {
	settings := newSettings()
	assert.NotNil(t, settings)
}

func TestGlobalSettings(t *testing.T) {
	t.Setenv(AppNameEnv, "ca-go-unit-tests")
	t.Setenv(AppVerEnv, "1.2.3")
	t.Setenv(AppEnvironmentEnv, "local")
	t.Setenv(AppFarmEnv, "local")
	t.Setenv(ProductEnv, "standard_library")

	settings := newSettings()
	assert.Equal(t, "ca-go-unit-tests", settings.App)
	assert.Equal(t, "1.2.3", settings.AppVersion)
	assert.Equal(t, "local", settings.AppEnv)
	assert.Equal(t, "local", settings.Farm)
	assert.Equal(t, "standard_library", settings.Product)
}

func TestAwsSettings(t *testing.T) {
	t.Setenv(AwsProfileEnv, "dev")
	t.Setenv(AwsRegionEnv, "us-west-1")
	t.Setenv(AwsAccountIDEnv, "123456789")
	t.Setenv(AwsXrayEnv, "true")

	settings := newSettings()
	assert.Equal(t, "dev", settings.AwsProfile)
	assert.Equal(t, "us-west-1", settings.AwsRegion)
	assert.Equal(t, "123456789", settings.AwsAccountID)
	assert.Equal(t, true, settings.XrayLogging)
}

func TestLogSettings(t *testing.T) {
	t.Setenv(LogLevelEnv, "WARN")

	settings := newSettings()
	assert.Equal(t, "WARN", settings.LogLevel)
}

func TestAuthZSettings(t *testing.T) {
	t.Setenv(AuthzClientTimeoutEnv, "123")
	t.Setenv(AuthzCacheDurationEnv, "456")
	t.Setenv(AuthzDialerTimeoutEnv, "789")
	t.Setenv(AuthzTLSTimeoutEnv, "10")

	settings := newSettings()
	assert.Equal(t, 123, settings.AuthzClientTimeoutInMs)
	assert.Equal(t, 456, settings.AuthzCacheDurationInSec)
	assert.Equal(t, 789, settings.AuthzDialerTimeoutInMs)
	assert.Equal(t, 10, settings.AuthzTLSTimeoutInMs)
}

func TestCacheSettings(t *testing.T) {
	t.Setenv(CacheDurationEnv, "1234")

	settings := newSettings()
	assert.Equal(t, 1234, settings.CacheDurationInSec)
}

func TestDatadogSettings(t *testing.T) {
	t.Setenv(DatadogAPIEnv, "api-123-key")
	t.Setenv(DatadogLogEndpointEnv, "dd-endpoint")
	t.Setenv(DatadogEnvironmentEnv, "us")
	t.Setenv(DatadogServiceEnv, "dd-service")
	t.Setenv(DatadogVersionEnv, "1.9.3")
	t.Setenv(DatadogAgentHostEnv, "dd.host.com")
	t.Setenv(DatadogStatsdPortEnv, "8888")
	t.Setenv(DatadogTimeoutEnv, "4321")
	t.Setenv(DatadogSiteEnv, "abc")
	t.Setenv(DatadogLogLevelEnv, "ERROR")

	settings := newSettings()
	assert.Equal(t, "api-123-key", settings.DatadogApiKey)
	assert.Equal(t, "dd-endpoint", settings.DatadogLogEndpoint)
	assert.Equal(t, "us", settings.DatadogEnv)
	assert.Equal(t, "dd-service", settings.DatadogService)
	assert.Equal(t, "1.9.3", settings.DatadogVersion)
	assert.Equal(t, "dd.host.com", settings.DatadogAgentHost)
	assert.Equal(t, 4321, settings.DatadogTimeoutInMs)
	assert.Equal(t, "abc", settings.DatadogSite)
	assert.Equal(t, "ERROR", settings.DatadogLogLevel)
}

func TestSentrySettings(t *testing.T) {
	t.Setenv(SentryDsnEnv, "sentry.dsn.com")
	t.Setenv(SentryFlushTimeoutInMsEnv, "1234")

	settings := newSettings()
	assert.Equal(t, "sentry.dsn.com", settings.SentryDSN)
	assert.Equal(t, 1234, settings.SentryFlushInMs)
}

func TestSettingsHelpers(t *testing.T) {
	t.Setenv(AppEnvironmentEnv, "production")
	settings := newSettings()
	isProd := settings.isProduction()
	assert.True(t, isProd)

	t.Setenv(AppEnvironmentEnv, "dev")
	settings = newSettings()
	isProd = settings.isProduction()
	assert.False(t, isProd)
}

func Test_Settings_Env_IsAws_IsLocal(t *testing.T) {
	defer os.Unsetenv(AppFarmEnv)

	t.Setenv(AppFarmEnv, "local")
	settings := newSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "local", settings.Farm)
	assert.True(t, settings.isRunningLocal())
	assert.False(t, settings.isRunningInAWS())

	t.Setenv(AppFarmEnv, "falcon")
	settings = newSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "falcon", settings.Farm)
	assert.False(t, settings.isRunningLocal())
	assert.True(t, settings.isRunningInAWS())

	t.Setenv(AppFarmEnv, "production")
	settings = newSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "production", settings.Farm)
	assert.False(t, settings.isRunningLocal())
	assert.True(t, settings.isRunningInAWS())
}
