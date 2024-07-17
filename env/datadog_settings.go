package env

import (
	senv "github.com/caarlos0/env/v11"
)

// DatadogSettings implements Datadog settings.
// This is an interface so that clients can mock out this behaviour in tests.
type DatadogSettings interface {
	DatadogAPIKey() string
	DatadogLogEndpoint() string
	DatadogEnv() string
	DatadogService() string
	DatadogVersion() string
	DatadogAgentHost() string
	DatadogStatsDPort() int
	DatadogTimeoutInMs() int
	DatadogSite() string
	DatadogLogLevel() string
}

// datadogSettings that drive behavior.
type datadogSettings struct {
	// These have to be public so that "github.com/caarlos0/env/v10" can populate them
	DDApiKeyEnv      string `env:"DD_API_KEY"`
	DDLogEndpointEnv string `env:"DD_LOG_ENDPOINT"`
	DDTypeEnv        string `env:"DD_ENV"            envDefault:"development"`
	DDServiceEnv     string `env:"DD_SERVICE"`
	DDVersionEnv     string `env:"DD_VERSION"        envDefault:"1.0.0"`
	DDAgentHostEnv   string `env:"DD_AGENT_HOST"`
	DDStatsDPortEnv  int    `env:"DD_DOGSTATSD_PORT" envDefault:"8125"`
	DDTimeoutInMsEnv int    `env:"DD_TIMEOUT"        envDefault:"500"`
	DDSiteEnv        string `env:"DD_SITE"`
	DDLogLevelEnv    string `env:"DD_LOG_LEVEL"      envDefault:"INFO"`
}

func newDatadogSettings() *datadogSettings {
	settings := datadogSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// DatadogAPIKey returns the "DD_API_KEY" environment variable.
func (s *datadogSettings) DatadogAPIKey() string {
	return s.DDApiKeyEnv
}

// DatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func (s *datadogSettings) DatadogLogEndpoint() string {
	return s.DDLogEndpointEnv
}

// DatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func (s *datadogSettings) DatadogEnv() string {
	return s.DDTypeEnv
}

// DatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func (s *datadogSettings) DatadogService() string {
	return s.DDServiceEnv
}

// DatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func (s *datadogSettings) DatadogVersion() string {
	return s.DDVersionEnv
}

// DatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func (s *datadogSettings) DatadogAgentHost() string {
	return s.DDAgentHostEnv
}

// DatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func (s *datadogSettings) DatadogStatsDPort() int {
	return s.DDStatsDPortEnv
}

// DatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func (s *datadogSettings) DatadogTimeoutInMs() int {
	return s.DDTimeoutInMsEnv
}

// DatadogSite returns the "DD_SITE" environment variable.
func (s *datadogSettings) DatadogSite() string {
	return s.DDSiteEnv
}

// DatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func (s *datadogSettings) DatadogLogLevel() string {
	return s.DDLogLevelEnv
}
