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
	assert.Equal(t, "api-123-key", settings.DDApiKeyEnv)
	assert.Equal(t, "dd-endpoint", settings.DDLogEndpointEnv)
	assert.Equal(t, "us", settings.DDTypeEnv)
	assert.Equal(t, "dd-service", settings.DDServiceEnv)
	assert.Equal(t, "1.9.3", settings.DDVersionEnv)
	assert.Equal(t, "dd.host.com", settings.DDAgentHostEnv)
	assert.Equal(t, 4321, settings.DDTimeoutInMsEnv)
	assert.Equal(t, "abc", settings.DDSiteEnv)
	assert.Equal(t, "ERROR", settings.DDLogLevelEnv)
}
