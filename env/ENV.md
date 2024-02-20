# ca-go/env

The `env` package provides access to common environment values . The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

There are no new settings to create or pass around, instead there is a singleton setting created in the package that you can call directly.

The env package provides access to many of the common field required for logging in line with our [logging standard](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard).
However, the [log](../log) package will add these fields by default, to avoid duplicating fields please see the documentation for the [log](../log) package [here](../log/LOGGER.md)

Note: The env package does NOT support redacting, so be mindful about logging any sensitive setting information.

## Environment Variables

Here is the list of supported environment variables currently supported:
- AppNameEnv        = "APP"
- AppVerEnv         = "APP_VERSION"
- AppEnvironmentEnv = "APP_ENV"
- AppFarmEnv        = "FARM"
- AppFarmLegacyEnv  = "APP_ENV"
- ProductEnv        = "PRODUCT"
- AwsProfileEnv   = "AWS_PROFILE"
- AwsRegionEnv    = "AWS_REGION"
- AwsAccountIDEnv = "AWS_ACCOUNT_ID"
- AwsXrayEnv      = "XRAY_LOGGING"
- LogLevelEnv = "LOG_LEVEL"
- AuthzClientTimeoutEnv = "AUTHZ_CLIENT_TIMEOUT_IN_MS"
- AuthzCacheDurationEnv = "AUTHZ_CACHE_DURATION_IN_SEC"
- AuthzDialerTimeoutEnv = "AUTHZ_DIALER_TIMEOUT_IN_MS"
- AuthzTLSTimeoutEnv    = "AUTHZ_TLS_TIMEOUT_IN_MS"
- CacheDurationEnv = "CACHE_DURATION_IN_SEC"
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
- SentryDsnEnv              = "SENTRY_DSN"
- SentryFlushTimeoutInMsEnv = "SENTRY_FLUSH_TIMEOUT_IN_MS"


## Methods

func AppName() string
func AppVersion() string [Default: "1.0.0"]
func AppEnv() string
func Farm() string
func ProductSuite() string
func AwsProfile() string
func AwsRegion() string
func AwsAccountID() string
func IsXrayTracingEnabled() bool
func LogLevel() string
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
func SentryDSN() string
func SentryFlushTimeoutInMs() int
func IsProduction() bool
func IsRunningInAWS() bool
func IsRunningLocal() bool


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
	isProd := env.IsProduction()
	isAWS := env.IsRunningInAWS()
	isLocal := env.IsRunningLocal()
}
```
