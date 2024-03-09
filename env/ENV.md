# ca-go/env

The `env` package provides access to common environment values . The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

The package creates default structs containing common environment variable values.


__Note__: The env package does NOT support redacting, so be mindful about logging any sensitive setting information.

## Environment Variables

Here is the list of supported environment variables currently supported:

### Common
- AppNameEnv        = "APP"
- AppVerEnv         = "APP_VERSION"
- AppEnvironmentEnv = "APP_ENV"
- AppFarmEnv        = "FARM"
- AppFarmLegacyEnv  = "APP_ENV"
- ProductEnv        = "PRODUCT"

### AWS
- AwsProfileEnv   = "AWS_PROFILE"
- AwsRegionEnv    = "AWS_REGION"
- AwsAccountIDEnv = "AWS_ACCOUNT_ID"
- AwsXrayEnv      = "XRAY_LOGGING"

### Logging
- LogLevelEnv = "LOG_LEVEL"

### Datadog
- DatadogAPIEnv         = "DD_API_KEY"
- DatadogLogEndpointEnv = "DD_LOG_ENDPOINT"
- DatadogEnvironmentEnv = "DD_ENV"
- DatadogServiceEnv     = "DD_SERVICE"
- DatadogVersionEnv     = "DD_VERSION"
- DatadogAgentHostEnv   = "DD_AGENT_HOST"
- DatadogStatsdPortEnv  = "DD_DOGSTATSD_PORT"
- DatadogTimeoutEnv     = "DD_TIMEOUT"
- DatadogSiteEnv        = "DD_SITE"
- DatadogLogLevelEnv    = "DD_LOG_LEVEL"

### SEntry
- SentryDsnEnv              = "SENTRY_DSN"
- SentryFlushTimeoutInMsEnv = "SENTRY_FLUSH_TIMEOUT_IN_MS"


## Methods

### Common
func AppName() string
func AppVersion() string [Default: "1.0.0"]
func AppEnv() string
func Farm() string
func ProductSuite() string
func IsProduction() bool
func IsRunningInAWS() bool
func IsRunningLocal() bool


### AWS
func AwsProfile() string
func AwsRegion() string
func AwsAccountID() string
func IsXrayTracingEnabled() bool

### Logging
func LogLevel() string

### Datadog
func DatadogApiKey() string
func DatadogLogEndpoint() string
func DatadogEnv() string
func DatadogService() string
func DatadogVersion() string
func DatadogAgentHost()
func DatadogStatsDPort()
func DatadogTimeoutInMs() int [Default: 500]
func DatadogSite() string
func DatadogLogLevel()

### Sentry
func SentryDSN() string
func SentryFlushTimeoutInMs() int

## Examples
```
package cago

import (
	"testing"

	"github.com/cultureamp/ca-go/env"
)

func SettingsExample(t *testing.T) {
	// call env methods to retrieve environment values

	appName := env.GetAppName()
	appVer := env.GetAppVersion()
	appEnv := env.GetAppEnv()
	farm := env.GetFarm()
	product := env.GetProductSuite()
	isProd := env.IsProduction()
	isAWS := env.IsRunningInAWS()
	isLocal := env.IsRunningLocal()
	awsProfile := env.GetAwsProfile()
	awsRegion := env.GetAwsRegion()
	awsAccountID := env.GetAwsAccountID()
	xrayEnabled := env.IsXrayTracingEnabled()
	logLevel := env.GetLogLevel()
	ddApiKey := env.GetDatadogApiKey()
	ddEndpoint := env.GetDatadogLogEndpoint()
	ddEnv := env.GetDatadogEnv()
	ddService := env.GetDatadogService()
	ddVersion := env.GetDatadogVersion()
	ddAgentHost := env.GetDatadogAgentHost()
	ddStatsDPort := env.GetDatadogStatsDPort()
	ddTimeout := env.GetDatadogTimeoutInMs()
	ddSite := env.GetDatadogSite()
	ddLogLevel := env.GetDatadogLogLevel()
	sentryDSN := env.GetSentryDSN()
	sentryFlushTimeout := env.GetSentryFlushTimeoutInMs()
}
```

## Testing and Mocks

During tests you can override the package level `DefaultXYZSettings` with a mock that supports the specific set of environemnt  interface.

```
import (
	"testing"

	"github.com/cultureamp/ca-go/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockPackageLevelMethods(t *testing.T) {
	// 1. set up your mock
	mock := new(mockSettings)
	mock.On("GetAppEnv").Return("abc")

	// 2. override the package level DefaultAWSSecrets.Client with your mock
	oldSettings := env.DefaultCommonSettings
	defer func() { env.DefaultCommonSettings = oldSettings }()
	env.DefaultCommonSettings = mock

	// 3. call the package methods which will call you mock
	app := env.AppEnv()
	assert.Equal(t, "abc", app)
	mock.AssertExpectations(t)
}
```