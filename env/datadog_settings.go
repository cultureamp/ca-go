package env

import (
	senv "github.com/caarlos0/env/v10"
)

// datadogSettings that drive behavior.
type datadogSettings struct {
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
}

func newDatadogSettings() *datadogSettings {
	settings := datadogSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// GetDatadogApiKey returns the "DD_API_KEY" environment variable.
func (s *datadogSettings) GetDatadogApiKey() string {
	return s.DatadogApiKey
}

// GetDatadogLogEndpoint returns the "DD_LOG_ENDPOINT" environment variable.
func (s *datadogSettings) GetDatadogLogEndpoint() string {
	return s.DatadogLogEndpoint
}

// GetDatadogEnv returns the "DD_ENV" environment variable.
// Default: AppEnv().
func (s *datadogSettings) GetDatadogEnv() string {
	return s.DatadogEnv
}

// GetDatadogService returns the "DD_SERVICE" environment variable.
// Default: App().
func (s *datadogSettings) GetDatadogService() string {
	return s.DatadogService
}

// GetDatadogVersion returns the "DD_VERSION" environment variable.
// Default: AppVersion().
func (s *datadogSettings) GetDatadogVersion() string {
	return s.DatadogVersion
}

// GetDatadogAgentHost returns the "DD_AGENT_HOST" environment variable.
func (s *datadogSettings) GetDatadogAgentHost() string {
	return s.DatadogAgentHost
}

// GetDatadogStatsDPort returns the "DD_DOGSTATSD_PORT" environment variable.
// Default: 8125.
func (s *datadogSettings) GetDatadogStatsDPort() int {
	return s.DatadogStatsDPort
}

// GetDatadogTimeoutInMs returns the "DD_TIMEOUT" environment variable.
// Default: 500.
func (s *datadogSettings) GetDatadogTimeoutInMs() int {
	return s.DatadogTimeoutInMs
}

// GetDatadogSite returns the "DD_SITE" environment variable.
func (s *datadogSettings) GetDatadogSite() string {
	return s.DatadogSite
}

// GetDatadogLogLevel returns the "DD_LOG_LEVEL" environment variable.
// Default: LogLevel().
func (s *datadogSettings) GetDatadogLogLevel() string {
	return s.DatadogLogLevel
}
