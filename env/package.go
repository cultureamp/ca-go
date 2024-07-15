package env

// DefaultCommonSettings is the package level instance of CommonSettings.
var DefaultCommonSettings CommonSettings = getCommonInstance()

func getCommonInstance() *commonSettings {
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
	return newAWSSettings()
}

// AwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func AwsProfile() string {
	return DefaultAWSSettings.AwsProfile()
}

// AwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func AwsRegion() string {
	return DefaultAWSSettings.AwsRegion()
}

// AwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
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

// LogLevel returns the "LOG_LEVEL" environment variable.
// Examples: "DEBUG, "INFO", "WARN", "ERROR".
func LogLevel() string {
	return DefaultLogSettings.LogLevel()
}

// DefaultDatadogSettings is the package level instance of DatadogSettings.
var DefaultDatadogSettings DatadogSettings = getDatadogInstance()

func getDatadogInstance() *datadogSettings {
	return newDatadogSettings()
}

// DatadogAPIKey returns the "DD_API_KEY" environment variable.
func DatadogAPIKey() string {
	return DefaultDatadogSettings.DatadogAPIKey()
}

// DatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func DatadogLogEndpoint() string {
	return DefaultDatadogSettings.DatadogLogEndpoint()
}

// DatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func DatadogEnv() string {
	return DefaultDatadogSettings.DatadogEnv()
}

// DatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func DatadogService() string {
	return DefaultDatadogSettings.DatadogService()
}

// DatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func DatadogVersion() string {
	return DefaultDatadogSettings.DatadogVersion()
}

// DatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func DatadogAgentHost() string {
	return DefaultDatadogSettings.DatadogAgentHost()
}

// DatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func DatadogStatsDPort() int {
	return DefaultDatadogSettings.DatadogStatsDPort()
}

// DatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func DatadogTimeoutInMs() int {
	return DefaultDatadogSettings.DatadogTimeoutInMs()
}

// DatadogSite returns the "DD_SITE" environment variable.
func DatadogSite() string {
	return DefaultDatadogSettings.DatadogSite()
}

// DatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func DatadogLogLevel() string {
	return DefaultDatadogSettings.DatadogLogLevel()
}

// DefaultSentrySettings is the package level instance of SentrySettings.
var DefaultSentrySettings SentrySettings = getSentryInstance()

func getSentryInstance() *sentrySettings {
	return newSentrySettings()
}

// SentryDSN returns the "SENTRY_DSN" environment variable.
func SentryDSN() string {
	return DefaultSentrySettings.SentryDSN()
}

// SentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func SentryFlushTimeoutInMs() int {
	return DefaultSentrySettings.SentryFlushTimeoutInMs()
}
