package env

type CommonSettings interface {
	GetAppName() string
	GetAppVersion() string
	GetAppEnv() string
	GetFarm() string
	GetProductSuite() string
	IsProduction() bool
	IsRunningInAWS() bool
	IsRunningLocal() bool
}

type AWSSettings interface {
	GetAwsProfile() string
	GetAwsRegion() string
	GetAwsAccountID() string
	IsXrayTracingEnabled() bool
}

type LogSettings interface {
	GetLogLevel() string
}

type DatadogSettings interface {
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
}

type SentrySettings interface {
	GetSentryDSN() string
	GetSentryFlushTimeoutInMs() int
}

var DefaultCommonSettings CommonSettings = getCommonInstance()

func getCommonInstance() *commonSettings {
	return newCommonSettings()
}

// GetAppName returns the application name from the "APP" environment variable.
func AppName() string {
	return DefaultCommonSettings.GetAppName()
}

// GetAppVersion returns the application version from the "APP_VER" environment variable.
// Default: "1.0.0".
func AppVersion() string {
	return DefaultCommonSettings.GetAppVersion()
}

// GetAppEnv returns the application environment from the "APP_ENV" environment variable.
// Examples: "development", "production".
func AppEnv() string {
	return DefaultCommonSettings.GetAppEnv()
}

// GetFarm returns the farm running the application from the "FARM" environment variable.
// Examples: "local", "dolly", "production".
func Farm() string {
	return DefaultCommonSettings.GetFarm()
}

// GetProductSuite returns the product suite this application belongs to from the "PRODUCT" environment variable.
// Examples: "engagement", "performance".
func ProductSuite() string {
	return DefaultCommonSettings.GetProductSuite()
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

var DefaultAWSSettings AWSSettings = getAWSInstance()

func getAWSInstance() *awsSettings {
	return newAWSSettings()
}

// GetAwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func AwsProfile() string {
	return DefaultAWSSettings.GetAwsProfile()
}

// GetAwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func AwsRegion() string {
	return DefaultAWSSettings.GetAwsRegion()
}

// GetAwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
func AwsAccountID() string {
	return DefaultAWSSettings.GetAwsAccountID()
}

// IsXrayTracingEnabled returns "true" if the "XRAY_LOGGING" environment variable is turned on.
func IsXrayTracingEnabled() bool {
	return DefaultAWSSettings.IsXrayTracingEnabled()
}

var DefaultLogSettings LogSettings = getLogInstance()

func getLogInstance() *logSettings {
	return newLogSettings()
}

// GetLogLevel returns the "LOG_LEVEL" environment variable.
// Examples: "DEBUG, "INFO", "WARN", "ERROR".
func LogLevel() string {
	return DefaultLogSettings.GetLogLevel()
}

var DefaultDatadogSettings DatadogSettings = getDatadogInstance()

func getDatadogInstance() *datadogSettings {
	return newDatadogSettings()
}

// GetDatadogApiKey returns the "DD_API_KEY" environment variable.
func DatadogApiKey() string {
	return DefaultDatadogSettings.GetDatadogApiKey()
}

// GetDatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func DatadogLogEndpoint() string {
	return DefaultDatadogSettings.GetDatadogLogEndpoint()
}

// GetDatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func DatadogEnv() string {
	return DefaultDatadogSettings.GetDatadogEnv()
}

// GetDatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func DatadogService() string {
	return DefaultDatadogSettings.GetDatadogService()
}

// GetDatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func DatadogVersion() string {
	return DefaultDatadogSettings.GetDatadogVersion()
}

// GetDatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func DatadogAgentHost() string {
	return DefaultDatadogSettings.GetDatadogAgentHost()
}

// GetDatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func DatadogStatsDPort() int {
	return DefaultDatadogSettings.GetDatadogStatsDPort()
}

// GetDatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func DatadogTimeoutInMs() int {
	return DefaultDatadogSettings.GetDatadogTimeoutInMs()
}

// GetDatadogSite returns the "DD_SITE" environment variable.
func DatadogSite() string {
	return DefaultDatadogSettings.GetDatadogSite()
}

// GetDatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func DatadogLogLevel() string {
	return DefaultDatadogSettings.GetDatadogLogLevel()
}

var DefaultSentrySettings SentrySettings = getSentryInstance()

func getSentryInstance() *sentrySettings {
	return newSentrySettings()
}

// GetSentryDSN returns the "SENTRY_DSN" environment variable.
func SentryDSN() string {
	return DefaultSentrySettings.GetSentryDSN()
}

// GetSentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func SentryFlushTimeoutInMs() int {
	return DefaultSentrySettings.GetSentryFlushTimeoutInMs()
}
