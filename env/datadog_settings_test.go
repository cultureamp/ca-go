package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDatadogSettings(t *testing.T) {
	settings := newDatadogSettings()
	assert.NotNil(t, settings)
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

	settings := newDatadogSettings()
	assert.Equal(t, "api-123-key", settings.DD_ApiKey)
	assert.Equal(t, "dd-endpoint", settings.DD_LogEndpoint)
	assert.Equal(t, "us", settings.DD_Env)
	assert.Equal(t, "dd-service", settings.DD_Service)
	assert.Equal(t, "1.9.3", settings.DD_Version)
	assert.Equal(t, "dd.host.com", settings.DD_AgentHost)
	assert.Equal(t, 4321, settings.DD_TimeoutInMs)
	assert.Equal(t, "abc", settings.DD_Site)
	assert.Equal(t, "ERROR", settings.DD_LogLevel)
}
