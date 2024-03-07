package env

import (
	senv "github.com/caarlos0/env/v10"
)

// standardSettings that drive behavior.
// Consumers of this library should inherit this by embedding this struct in their own standardSettings.
type standardSettings struct {
	// Common environment variable values used by at least 80% of apps
	App        string `env:"APP"`
	AppVersion string `env:"APP_VERSION" envDefault:"1.0.0"`
	AppEnv     string `env:"APP_ENV"     envDefault:"development"`
	Farm       string `env:"FARM"        envDefault:"local"`
	Product    string `env:"PRODUCT"`

	// Aws
	AwsProfile   string `env:"AWS_PROFILE"    envDefault:"default"`
	AwsRegion    string `env:"AWS_REGION"`
	AwsAccountID string `env:"AWS_ACCOUNT_ID"`
	XrayLogging  bool   `env:"XRAY_LOGGING"   envDefault:"true"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`

	// AuthZ
	AuthzClientTimeoutInMs  int `env:"AUTHZ_CLIENT_TIMEOUT_IN_MS"  envDefault:"1000"`
	AuthzCacheDurationInSec int `env:"AUTHZ_CACHE_DURATION_IN_SEC" envDefault:"0"`
	AuthzDialerTimeoutInMs  int `env:"AUTHZ_DIALER_TIMEOUT_IN_MS"  envDefault:"100"`
	AuthzTLSTimeoutInMs     int `env:"AUTHZ_TLS_TIMEOUT_IN_MS"     envDefault:"500"`

	// Cache
	CacheDurationInSec int `env:"CACHE_DURATION_IN_SEC" envDefault:"0"`

	// Datadog
	DatadogApiKey      string `env:"DD_API_KEY"`
	DatadogLogEndpoint string `env:"DD_LOG_ENDPOINT"`
	DatadogEnv         string `env:"DD_ENV"            envDefault:"development"`
	DatadogService     string `env:"DD_SERVICE"`
	DatadogVersion     string `env:"DD_VERSION"        envDefault:"1.0.0"`
	DatadogAgentHost   string `env:"DD_AGENT_HOST"`
	DatadogStatsDPort  int    `env:"DD_DOGSTATSD_PORT" envDefault:"8125"`
	DatadogTimeoutInMs int    `env:"DD_TIMEOUT"        envDefault:"500"`
	DatadogSite        string `env:"DD_SITE"`
	DatadogLogLevel    string `env:"DD_LOG_LEVEL"      envDefault:"INFO"`

	// Sentry - soon to be deprecated
	SentryDSN       string `env:"SENTRY_DSN"`
	SentryFlushInMs int    `env:"SENTRY_FLUSH_TIMEOUT_IN_MS" envDefault:"100"`
}

func newSettings() *standardSettings {
	settings := standardSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// GetAppName returns the application name from the "APP" environment variable.
func (s *standardSettings) GetAppName() string {
	return s.App
}

// GetAppVersion returns the application version from the "APP_VER" environment variable.
// Default: "1.0.0".
func (s *standardSettings) GetAppVersion() string {
	return s.AppVersion
}

// GetAppEnv returns the application environment from the "APP_ENV" environment variable.
// Examples: "development", "production".
func (s *standardSettings) GetAppEnv() string {
	return s.AppEnv
}

// GetFarm returns the farm running the application from the "FARM" environment variable.
// Examples: "local", "dolly", "production".
func (s *standardSettings) GetFarm() string {
	return s.Farm
}

// GetProductSuite returns the product suite this application belongs to from the "PRODUCT" environment variable.
// Examples: "engagement", "performance".
func (s *standardSettings) GetProductSuite() string {
	return s.Product
}

// GetAwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func (s *standardSettings) GetAwsProfile() string {
	return s.AwsProfile
}

// GetAwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func (s *standardSettings) GetAwsRegion() string {
	return s.AwsRegion
}

// GetAwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
func (s *standardSettings) GetAwsAccountID() string {
	return s.AwsAccountID
}

// IsXrayTracingEnabled returns "true" if the "XRAY_LOGGING" environment variable is turned on.
func (s *standardSettings) IsXrayTracingEnabled() bool {
	return s.XrayLogging
}

// GetLogLevel returns the "LOG_LEVEL" environment variable.
// Examples: "DEBUG, "INFO", "WARN", "ERROR".
func (s *standardSettings) GetLogLevel() string {
	return s.LogLevel
}

// GetDatadogApiKey returns the "DD_API_KEY" environment variable.
func (s *standardSettings) GetDatadogApiKey() string {
	return s.DatadogApiKey
}

// GetDatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func (s *standardSettings) GetDatadogLogEndpoint() string {
	return s.DatadogLogEndpoint
}

// GetDatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func (s *standardSettings) GetDatadogEnv() string {
	return s.DatadogEnv
}

// GetDatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func (s *standardSettings) GetDatadogService() string {
	return s.DatadogService
}

// GetDatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func (s *standardSettings) GetDatadogVersion() string {
	return s.DatadogVersion
}

// GetDatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func (s *standardSettings) GetDatadogAgentHost() string {
	return s.DatadogAgentHost
}

// GetDatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func (s *standardSettings) GetDatadogStatsDPort() int {
	return s.DatadogStatsDPort
}

// GetDatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func (s *standardSettings) GetDatadogTimeoutInMs() int {
	return s.DatadogTimeoutInMs
}

// GetDatadogSite returns the "DD_SITE" environment variable.
func (s *standardSettings) GetDatadogSite() string {
	return s.DatadogSite
}

// GetDatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func (s *standardSettings) GetDatadogLogLevel() string {
	return s.DatadogLogLevel
}

// GetSentryDSN returns the "SENTRY_DSN" environment variable.
func (s *standardSettings) GetSentryDSN() string {
	return s.SentryDSN
}

// GetSentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func (s *standardSettings) GetSentryFlushTimeoutInMs() int {
	return s.SentryFlushInMs
}

// IsProduction returns true if "APP_ENV" == "production".
func (s *standardSettings) IsProduction() bool {
	return s.AppEnv == "production"
}

// IsRunningInAWS returns true if "APP_ENV" != "local".
func (s *standardSettings) IsRunningInAWS() bool {
	return !s.IsRunningLocal()
}

// IsRunningLocal returns true if FARM" == "local".
func (s *standardSettings) IsRunningLocal() bool {
	return s.Farm == "local"
}
