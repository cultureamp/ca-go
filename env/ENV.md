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

	appName := env.AppName()
	appVer := env.AppVersion()
	appEnv := env.AppEnv()
	farm := env.Farm()
	product := env.ProductSuite()
	isProd := env.IsProduction()
	isAWS := env.IsRunningInAWS()
	isLocal := env.IsRunningLocal()
	awsProfile := env.AwsProfile()
	awsRegion := env.AwsRegion()
	awsAccountID := env.AwsAccountID()
	xrayEnabled := env.IsXrayTracingEnabled()
	logLevel := env.LogLevel()
	ddApiKey := env.DatadogApiKey()
	ddEndpoint := env.DatadogLogEndpoint()
	ddEnv := env.DatadogEnv()
	ddService := env.DatadogService()
	ddVersion := env.DatadogVersion()
	ddAgentHost := env.DatadogAgentHost()
	ddStatsDPort := env.DatadogStatsDPort()
	ddTimeout := env.DatadogTimeoutInMs()
	ddSite := env.DatadogSite()
	ddLogLevel := env.DatadogLogLevel()
	sentryDSN := env.SentryDSN()
	sentryFlushTimeout := env.SentryFlushTimeoutInMs()
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
	mock.On("AppEnv").Return("abc")

	// 2. override the package level DefaultCommonSettings (or which ever settings you like) with your mock
	oldSettings := env.DefaultCommonSettings
	defer func() { env.DefaultCommonSettings = oldSettings }()
	env.DefaultCommonSettings = mock

	// 3. call the package methods which will call you mock
	app := env.AppEnv()
	assert.Equal(t, "abc", app)
	mock.AssertExpectations(t)
}
```