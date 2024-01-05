package env_test

import (
	"testing"

	"github.com/cultureamp/ca-go/env"
	"github.com/stretchr/testify/assert"
)

func TestSettingsExample(t *testing.T) {
	// call env methods to retrieve environment values

	// get the application name
	appName := env.AppName()
	assert.Equal(t, "", appName)

	appVer := env.AppVersion()
	assert.Equal(t, "1.0.0", appVer)

	appEnv := env.AppEnv()
	assert.Equal(t, "development", appEnv)

	farm := env.Farm()
	assert.Equal(t, "local", farm)

	product := env.ProductSuite()
	assert.Equal(t, "", product)

	awsProfile := env.AwsProfile()
	assert.Equal(t, "default", awsProfile)

	awsRegion := env.AwsRegion()
	assert.Equal(t, "", awsRegion)

	awsAccountID := env.AwsAccountID()
	assert.Equal(t, "", awsAccountID)

	xrayEnabled := env.IsXrayTracingEnabled()
	assert.Equal(t, true, xrayEnabled)

	logLevel := env.LogLevel()
	assert.Equal(t, "INFO", logLevel)

	ddApiKey := env.DatadogApiKey()
	assert.Equal(t, "", ddApiKey)

	ddEndpoint := env.DatadogLogEndpoint()
	assert.Equal(t, "", ddEndpoint)

	ddEnv := env.DatadogEnv()
	assert.Equal(t, "development", ddEnv)

	ddService := env.DatadogService()
	assert.Equal(t, "", ddService)

	ddVersion := env.DatadogVersion()
	assert.Equal(t, "1.0.0", ddVersion)

	ddAgentHost := env.DatadogAgentHost()
	assert.Equal(t, "", ddAgentHost)

	ddStatsDPort := env.DatadogStatsDPort()
	assert.Equal(t, 8125, ddStatsDPort)

	ddTimeout := env.DatadogTimeoutInMs()
	assert.Equal(t, 500, ddTimeout)

	ddSite := env.DatadogSite()
	assert.Equal(t, "", ddSite)

	ddLogLevel := env.DatadogLogLevel()
	assert.Equal(t, "INFO", ddLogLevel)

	sentryDSN := env.SentryDSN()
	assert.Equal(t, "", sentryDSN)

	sentryFlushTimeout := env.SentryFlushTimeoutInMs()
	assert.Equal(t, 100, sentryFlushTimeout)
}

func TestHelperExample(t *testing.T) {
	// call env methods to retrieve environment values

	isProd := env.IsProduction()
	assert.Equal(t, false, isProd)

	isAWS := env.IsRunningInAWS()
	assert.Equal(t, false, isAWS)

	isLocal := env.IsRunningLocal()
	assert.Equal(t, true, isLocal)
}
