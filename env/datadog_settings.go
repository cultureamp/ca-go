package env

import (
	senv "github.com/caarlos0/env/v10"
)

// DatadogSettings implements Datadog settings.
// This is an interface so that clients can mock out this behaviour in tests.
type DatadogSettings interface {
	DatadogApiKey() string
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
	DSDatadogApiKey      string `env:"DD_API_KEY"`
	DSDatadogLogEndpoint string `env:"DD_LOG_ENDPOINT"`
	DSDatadogEnv         string `env:"DD_ENV"            envDefault:"development"`
	DSDatadogService     string `env:"DD_SERVICE"`
	DSDatadogVersion     string `env:"DD_VERSION"        envDefault:"1.0.0"`
	DSDatadogAgentHost   string `env:"DD_AGENT_HOST"`
	DSDatadogStatsDPort  int    `env:"DD_DOGSTATSD_PORT" envDefault:"8125"`
	DSDatadogTimeoutInMs int    `env:"DD_TIMEOUT"        envDefault:"500"`
	DSDatadogSite        string `env:"DD_SITE"`
	DSDatadogLogLevel    string `env:"DD_LOG_LEVEL"      envDefault:"INFO"`
}

func newDatadogSettings() *datadogSettings {
	settings := datadogSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// DatadogApiKey returns the "DD_API_KEY" environment variable.
func (s *datadogSettings) DatadogApiKey() string {
	return s.DSDatadogApiKey
}

// DatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func (s *datadogSettings) DatadogLogEndpoint() string {
	return s.DSDatadogLogEndpoint
}

// DatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func (s *datadogSettings) DatadogEnv() string {
	return s.DSDatadogEnv
}

// DatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func (s *datadogSettings) DatadogService() string {
	return s.DSDatadogService
}

// DatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func (s *datadogSettings) DatadogVersion() string {
	return s.DSDatadogVersion
}

// DatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func (s *datadogSettings) DatadogAgentHost() string {
	return s.DSDatadogAgentHost
}

// DatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func (s *datadogSettings) DatadogStatsDPort() int {
	return s.DSDatadogStatsDPort
}

// DatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func (s *datadogSettings) DatadogTimeoutInMs() int {
	return s.DSDatadogTimeoutInMs
}

// DatadogSite returns the "DD_SITE" environment variable.
func (s *datadogSettings) DatadogSite() string {
	return s.DSDatadogSite
}

// DatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func (s *datadogSettings) DatadogLogLevel() string {
	return s.DSDatadogLogLevel
}
