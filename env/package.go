package env

type Settings interface {
	GetAppName() string
	GetAppVersion() string
	GetAppEnv() string
	GetFarm() string
	GetProductSuite() string
	GetAwsProfile() string
	GetAwsRegion() string
	GetAwsAccountID() string
	IsXrayTracingEnabled() bool
	GetLogLevel() string
	GetDatadogApiKey() string
	GetDatadogLogEndpoint() string
	GetDatadogEnv() string
	GetDatadogService() string
	GetDatadogVersion() string
	GetDatadogAgentHost() string
	GetDatadogStatsDPort() int
	GetDatadogTimeoutInMs() int
	GetDatadogSite() string
	GetDatadogLogLevel() string
	GetSentryDSN() string
	GetSentryFlushTimeoutInMs() int
	IsProduction() bool
	IsRunningInAWS() bool
	IsRunningLocal() bool
}

var DefaultSettings Settings = getInstance()

func getInstance() *standardSettings {
	return newSettings()
}

// GetAppName returns the application name from the "APP" environment variable.
func AppName() string {
	return DefaultSettings.GetAppName()
}

// GetAppVersion returns the application version from the "APP_VER" environment variable.
// Default: "1.0.0".
func AppVersion() string {
	return DefaultSettings.GetAppVersion()
}

// GetAppEnv returns the application environment from the "APP_ENV" environment variable.
// Examples: "development", "production".
func AppEnv() string {
	return DefaultSettings.GetAppEnv()
}

// GetFarm returns the farm running the application from the "FARM" environment variable.
// Examples: "local", "dolly", "production".
func Farm() string {
	return DefaultSettings.GetFarm()
}

// GetProductSuite returns the product suite this application belongs to from the "PRODUCT" environment variable.
// Examples: "engagement", "performance".
func ProductSuite() string {
	return DefaultSettings.GetProductSuite()
}

// GetAwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func AwsProfile() string {
	return DefaultSettings.GetAwsProfile()
}

// GetAwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func AwsRegion() string {
	return DefaultSettings.GetAwsRegion()
}

// GetAwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
func AwsAccountID() string {
	return DefaultSettings.GetAwsAccountID()
}

// IsXrayTracingEnabled returns "true" if the "XRAY_LOGGING" environment variable is turned on.
func IsXrayTracingEnabled() bool {
	return DefaultSettings.IsXrayTracingEnabled()
}

// GetLogLevel returns the "LOG_LEVEL" environment variable.
// Examples: "DEBUG, "INFO", "WARN", "ERROR".
func LogLevel() string {
	return DefaultSettings.GetLogLevel()
}

// GetDatadogApiKey returns the "DD_API_KEY" environment variable.
func DatadogApiKey() string {
	return DefaultSettings.GetDatadogApiKey()
}

// GetDatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func DatadogLogEndpoint() string {
	return DefaultSettings.GetDatadogLogEndpoint()
}

// GetDatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func DatadogEnv() string {
	return DefaultSettings.GetDatadogEnv()
}

// GetDatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func DatadogService() string {
	return DefaultSettings.GetDatadogService()
}

// GetDatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func DatadogVersion() string {
	return DefaultSettings.GetDatadogVersion()
}

// GetDatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func DatadogAgentHost() string {
	return DefaultSettings.GetDatadogAgentHost()
}

// GetDatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func DatadogStatsDPort() int {
	return DefaultSettings.GetDatadogStatsDPort()
}

// GetDatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func DatadogTimeoutInMs() int {
	return DefaultSettings.GetDatadogTimeoutInMs()
}

// GetDatadogSite returns the "DD_SITE" environment variable.
func DatadogSite() string {
	return DefaultSettings.GetDatadogSite()
}

// GetDatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func DatadogLogLevel() string {
	return DefaultSettings.GetDatadogLogLevel()
}

// GetSentryDSN returns the "SENTRY_DSN" environment variable.
func SentryDSN() string {
	return DefaultSettings.GetSentryDSN()
}

// GetSentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func SentryFlushTimeoutInMs() int {
	return DefaultSettings.GetSentryFlushTimeoutInMs()
}

// IsProduction returns true if "APP_ENV" == "production".
func IsProduction() bool {
	return DefaultSettings.IsProduction()
}

// IsRunningInAWS returns true if "APP_ENV" != "local".
func IsRunningInAWS() bool {
	return DefaultSettings.IsRunningInAWS()
}

// IsRunningLocal returns true if FARM" == "local".
func IsRunningLocal() bool {
	return DefaultSettings.IsRunningLocal()
}
