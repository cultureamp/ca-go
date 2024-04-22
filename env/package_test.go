package env_test

import (
	"testing"

	"github.com/cultureamp/ca-go/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultSettings(t *testing.T) {
	appName := env.AppName()
	assert.Equal(t, "<unknown>", appName)

	appVer := env.AppVersion()
	assert.Equal(t, "1.0.0", appVer)

	appEnv := env.AppEnv()
	assert.Equal(t, "development", appEnv)

	farm := env.Farm()
	assert.Equal(t, "local", farm)

	product := env.ProductSuite()
	assert.Equal(t, "<unknown>", product)

	awsProfile := env.AwsProfile()
	assert.Equal(t, "default", awsProfile)

	awsRegion := env.AwsRegion()
	assert.Equal(t, "dev", awsRegion)

	awsAccountID := env.AwsAccountID()
	assert.Equal(t, "<unknown>", awsAccountID)

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

func TestIsHelpers(t *testing.T) {
	isProd := env.IsProduction()
	assert.Equal(t, false, isProd)

	isAWS := env.IsRunningInAWS()
	assert.Equal(t, false, isAWS)

	isLocal := env.IsRunningLocal()
	assert.Equal(t, true, isLocal)
}

func TestMockPackageLevelMethods(t *testing.T) {
	// 1. set up your mock
	mock := new(mockSettings)
	mock.On("AppEnv").Return("abc")

	// 2. override the package level DefaultAWSSecrets.Client with your mock
	oldSettings := env.DefaultCommonSettings
	defer func() { env.DefaultCommonSettings = oldSettings }()
	env.DefaultCommonSettings = mock

	// 3. call the package methods which will call you mock
	app := env.AppEnv()
	assert.Equal(t, "abc", app)
	mock.AssertExpectations(t)
}

type mockSettings struct {
	mock.Mock
}

func (_m *mockSettings) AppName() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) AppVersion() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) AppEnv() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) Farm() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) ProductSuite() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetAwsProfile() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetAwsRegion() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetAwsAccountID() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) IsXrayTracingEnabled() bool {
	args := _m.Called()
	output, _ := args.Get(0).(bool)
	return output
}

func (_m *mockSettings) GetLogLevel() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogApiKey() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogLogEndpoint() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogEnv() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogService() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogVersion() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogAgentHost() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogStatsDPort() int {
	args := _m.Called()
	output, _ := args.Get(0).(int)
	return output
}

func (_m *mockSettings) GetDatadogTimeoutInMs() int {
	args := _m.Called()
	output, _ := args.Get(0).(int)
	return output
}

func (_m *mockSettings) GetDatadogSite() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetDatadogLogLevel() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetSentryDSN() string {
	args := _m.Called()
	output, _ := args.Get(0).(string)
	return output
}

func (_m *mockSettings) GetSentryFlushTimeoutInMs() int {
	args := _m.Called()
	output, _ := args.Get(0).(int)
	return output
}

func (_m *mockSettings) IsProduction() bool {
	args := _m.Called()
	output, _ := args.Get(0).(bool)
	return output
}

func (_m *mockSettings) IsRunningInAWS() bool {
	args := _m.Called()
	output, _ := args.Get(0).(bool)
	return output
}

func (_m *mockSettings) IsRunningLocal() bool {
	args := _m.Called()
	output, _ := args.Get(0).(bool)
	return output
}
