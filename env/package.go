package env

import (
	"flag"
	"os"
	"strings"
)

// DefaultCommonSettings is the package level instance of CommonSettings.
var DefaultCommonSettings CommonSettings = getCommonInstance()

func getCommonInstance() *commonSettings {
	if isTestMode() {
		err := os.Setenv(AppNameEnv, "test-app")
		if err != nil {
			panic(err)
		}
	}
	return newCommonSettings()
}

// AppName returns the application name from the "APP" environment variable.
func AppName() string {
	return DefaultCommonSettings.AppName()
}

// AppVersion returns the application version from the "APP_VER" environment variable.
// Default: "1.0.0".
func AppVersion() string {
	return DefaultCommonSettings.AppVersion()
}

// AppEnv returns the application environment from the "APP_ENV" environment variable.
// Examples: "development", "production".
func AppEnv() string {
	return DefaultCommonSettings.AppEnv()
}

// Farm returns the farm running the application from the "FARM" environment variable.
// Examples: "local", "dolly", "production".
func Farm() string {
	return DefaultCommonSettings.Farm()
}

// ProductSuite returns the product suite this application belongs to from the "PRODUCT" environment variable.
// Examples: "engagement", "performance".
func ProductSuite() string {
	return DefaultCommonSettings.ProductSuite()
}

// IsProduction returns true if "APP_ENV" == "production".
func IsProduction() bool {
	return DefaultCommonSettings.IsProduction()
}

// IsRunningInAWS returns true if "APP_ENV" != "local".
func IsRunningInAWS() bool {
	return DefaultCommonSettings.IsRunningInAWS()
}

// IsRunningLocal returns true if FARM" == "local".
func IsRunningLocal() bool {
	return DefaultCommonSettings.IsRunningLocal()
}

// DefaultAWSSettings is the package level instance of AWSSettings.
var DefaultAWSSettings AWSSettings = getAWSInstance()

func getAWSInstance() *awsSettings {
	if isTestMode() {
		err := os.Setenv(AwsRegionEnv, "dev")
		if err != nil {
			panic(err)
		}
	}
	return newAWSSettings()
}

// GetAwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func AwsProfile() string {
	return DefaultAWSSettings.AwsProfile()
}

// GetAwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func AwsRegion() string {
	return DefaultAWSSettings.AwsRegion()
}

// GetAwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
func AwsAccountID() string {
	return DefaultAWSSettings.AwsAccountID()
}

// IsXrayTracingEnabled returns "true" if the "XRAY_LOGGING" environment variable is turned on.
func IsXrayTracingEnabled() bool {
	return DefaultAWSSettings.IsXrayTracingEnabled()
}

// DefaultLogSettings is the package level instance of LoggingSettings.
var DefaultLogSettings LogSettings = getLogInstance()

func getLogInstance() *logSettings {
	return newLogSettings()
}

// GetLogLevel returns the "LOG_LEVEL" environment variable.
// Examples: "DEBUG, "INFO", "WARN", "ERROR".
func LogLevel() string {
	return DefaultLogSettings.LogLevel()
}

// DefaultDatadogSettings is the package level instance of DatadogSettings.
var DefaultDatadogSettings DatadogSettings = getDatadogInstance()

func getDatadogInstance() *datadogSettings {
	return newDatadogSettings()
}

// GetDatadogApiKey returns the "DD_API_KEY" environment variable.
func DatadogApiKey() string {
	return DefaultDatadogSettings.DatadogApiKey()
}

// GetDatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func DatadogLogEndpoint() string {
	return DefaultDatadogSettings.DatadogLogEndpoint()
}

// GetDatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func DatadogEnv() string {
	return DefaultDatadogSettings.DatadogEnv()
}

// GetDatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func DatadogService() string {
	return DefaultDatadogSettings.DatadogService()
}

// GetDatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func DatadogVersion() string {
	return DefaultDatadogSettings.DatadogVersion()
}

// GetDatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func DatadogAgentHost() string {
	return DefaultDatadogSettings.DatadogAgentHost()
}

// GetDatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func DatadogStatsDPort() int {
	return DefaultDatadogSettings.DatadogStatsDPort()
}

// GetDatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func DatadogTimeoutInMs() int {
	return DefaultDatadogSettings.DatadogTimeoutInMs()
}

// GetDatadogSite returns the "DD_SITE" environment variable.
func DatadogSite() string {
	return DefaultDatadogSettings.DatadogSite()
}

// GetDatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func DatadogLogLevel() string {
	return DefaultDatadogSettings.DatadogLogLevel()
}

// DefaultSentrySettings is the package level instance of SentrySettings.
var DefaultSentrySettings SentrySettings = getSentryInstance()

func getSentryInstance() *sentrySettings {
	return newSentrySettings()
}

// GetSentryDSN returns the "SENTRY_DSN" environment variable.
func SentryDSN() string {
	return DefaultSentrySettings.SentryDSN()
}

// GetSentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func SentryFlushTimeoutInMs() int {
	return DefaultSentrySettings.SentryFlushTimeoutInMs()
}

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}
