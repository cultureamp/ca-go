package env

import (
	"os"
)

// Settings that drive behavior.
// Consumers of this library should inherit this by embedding this struct in their own Settings.
type Settings struct {
	// Common environment variable values used by at least 80% of apps
	App        string
	AppVersion string
	AppEnv     string
	Farm       string
	Product    string
	// Aws
	AwsProfile   string
	AwsRegion    string
	AwsAccountID string
	XrayLogging  bool
	// Logging
	LogLevel string
	// AuthZ
	AuthzClientTimeoutInMs  int
	AuthzCacheDurationInSec int
	AuthzDialerTimeoutInMs  int
	AuthzTLSTimeoutInMs     int
	// Cache
	CacheDurationInSec int
	// Datadog
	DatadogApiKey      string
	DatadogLogEndpoint string
	DatadogEnv         string
	DatadogService     string
	DatadogVersion     string
	DatadogAgentHost   string
	DatadogStatsDPort  int
	DatadogTimeoutInMs int
	DatadogSite        string
	DatadogLogLevel    string
	// Sentry - soon to be deprecated
	SentryDSN       string
	SentryFlushInMs int
}

var defaultSettings *Settings = getInstance()

func getInstance() *Settings {
	return newSettings()
}

func newSettings() *Settings {
	settings := &Settings{}

	// Global
	settings.App = os.Getenv(AppNameEnv)
	settings.AppVersion = GetString(AppVerEnv, "1.0.0")
	settings.AppEnv = GetString(AppEnvironmentEnv, "development")
	settings.Farm = GetString(AppFarmEnv, "local")
	settings.Product = os.Getenv(ProductEnv)
	// AWS
	settings.AwsProfile = GetString(AwsProfileEnv, "default")
	settings.AwsRegion = os.Getenv(AwsRegionEnv)
	settings.AwsAccountID = os.Getenv(AwsAccountIDEnv)
	settings.XrayLogging = GetBool(AwsXrayEnv, true)
	// Logging
	settings.LogLevel = GetString(LogLevelEnv, "INFO")
	// AuthZ
	settings.AuthzClientTimeoutInMs = GetInt(AuthzClientTimeoutEnv, 1000)
	settings.AuthzCacheDurationInSec = GetInt(AuthzCacheDurationEnv, 0)
	settings.AuthzDialerTimeoutInMs = GetInt(AuthzDialerTimeoutEnv, 100)
	settings.AuthzTLSTimeoutInMs = GetInt(AuthzTLSTimeoutEnv, 500)
	// Cache
	settings.CacheDurationInSec = GetInt(CacheDurationEnv, 0)
	// Datadog
	settings.DatadogApiKey = os.Getenv(DatadogAPIEnv)
	settings.DatadogLogEndpoint = os.Getenv(DatadogLogEndpointEnv)
	settings.DatadogEnv = GetString(DatadogEnvironmentEnv, settings.AppEnv)
	settings.DatadogService = GetString(DatadogServiceEnv, settings.App)
	settings.DatadogVersion = GetString(DatadogVersionEnv, settings.AppVersion)
	settings.DatadogAgentHost = os.Getenv(DatadogAgentHostEnv)
	settings.DatadogStatsDPort = GetInt(DatadogStatsdPortEnv, 8125)
	settings.DatadogTimeoutInMs = GetInt(DatadogTimeoutEnv, 500)
	settings.DatadogSite = os.Getenv(DatadogSiteEnv)
	settings.DatadogLogLevel = GetString(DatadogLogLevelEnv, settings.LogLevel)
	// Sentry - soon to be deprecated
	settings.SentryDSN = os.Getenv(SentryDsnEnv)
	settings.SentryFlushInMs = GetInt(SentryFlushTimeoutInMsEnv, 100)

	return settings
}

// AppName returns the application name from the "APP" environment variable.
func AppName() string {
	return defaultSettings.App
}

// AppVersion returns the application version from the "APP_VER" environment variable.
// Default: "1.0.0".
func AppVersion() string {
	return defaultSettings.AppVersion
}

// AppEnv returns the application environment from the "APP_ENV" environment variable.
// Examples: "development", "production".
func AppEnv() string {
	return defaultSettings.AppEnv
}

// Farm returns the farm running the application from the "FARM" environment variable.
// Examples: "local", "dolly", "production".
func Farm() string {
	return defaultSettings.Farm
}

// ProductSuite returns the product suite this application belongs to from the "PRODUCT" environment variable.
// Examples: "engagement", "performance".
func ProductSuite() string {
	return defaultSettings.Product
}

// AwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func AwsProfile() string {
	return defaultSettings.AwsProfile
}

// AwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func AwsRegion() string {
	return defaultSettings.AwsRegion
}

// AwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
func AwsAccountID() string {
	return defaultSettings.AwsAccountID
}

// IsXrayTracingEnabled returns "true" if the "XRAY_LOGGING" environment variable is turned on.
func IsXrayTracingEnabled() bool {
	return defaultSettings.XrayLogging
}

// LogLevel returns the "LOG_LEVEL" environment variable.
// Examles: "DEBUG, "INFO", "WARN", "ERROR".
func LogLevel() string {
	return defaultSettings.LogLevel
}

// DatadogApiKey returns the "DD_API_KEY" environment variable.
func DatadogApiKey() string {
	return defaultSettings.DatadogApiKey
}

// DatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func DatadogLogEndpoint() string {
	return defaultSettings.DatadogLogEndpoint
}

// DatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func DatadogEnv() string {
	return defaultSettings.DatadogEnv
}

// DatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func DatadogService() string {
	return defaultSettings.DatadogService
}

// DatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func DatadogVersion() string {
	return defaultSettings.DatadogVersion
}

// DatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func DatadogAgentHost() string {
	return defaultSettings.DatadogAgentHost
}

// DatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func DatadogStatsDPort() int {
	return defaultSettings.DatadogStatsDPort
}

// DatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func DatadogTimeoutInMs() int {
	return defaultSettings.DatadogTimeoutInMs
}

// DatadogSite returns the "DD_SITE" environment variable.
func DatadogSite() string {
	return defaultSettings.DatadogSite
}

// DatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func DatadogLogLevel() string {
	return defaultSettings.DatadogLogLevel
}

// SentryDSN returns the "SENTRY_DSN" environment variable.
func SentryDSN() string {
	return defaultSettings.SentryDSN
}

// SentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func SentryFlushTimeoutInMs() int {
	return defaultSettings.SentryFlushInMs
}

// IsProduction returns true if "APP_ENV" == "production".
func IsProduction() bool {
	return defaultSettings.isProduction()
}

func (s *Settings) isProduction() bool {
	return s.AppEnv == "production"
}

// IsRunningInAWS returns true if "APP_ENV" == "production".
func IsRunningInAWS() bool {
	return defaultSettings.isRunningInAWS()
}

func (s *Settings) isRunningInAWS() bool {
	return !s.isRunningLocal()
}

// IsRunningLocal returns true if FARM" == "local".
func IsRunningLocal() bool {
	return defaultSettings.Farm == "local"
}

func (s *Settings) isRunningLocal() bool {
	return s.Farm == "local"
}
