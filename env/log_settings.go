package env

import (
	senv "github.com/caarlos0/env/v11"
)

// LogSettings implements Logging settings.
// This is an interface so that clients can mock out this behaviour in tests.
type LogSettings interface {
	LogLevel() string
}

// logSettings that drive behavior.
type logSettings struct {
	// These have to be public so that "github.com/caarlos0/env/v10" can populate them
	LogLevelEnv string `env:"LOG_LEVEL" envDefault:"INFO"`
}

func newLogSettings() *logSettings {
	settings := logSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// LogLevel returns the "LOG_LEVEL" environment variable.
// Examples: "DEBUG, "INFO", "WARN", "ERROR".
func (s *logSettings) LogLevel() string {
	return s.LogLevelEnv
}
