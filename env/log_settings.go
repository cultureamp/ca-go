package env

import (
	senv "github.com/caarlos0/env/v10"
)

// logSettings that drive behavior.
type logSettings struct {
	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
}

func newLogSettings() *logSettings {
	settings := logSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// GetLogLevel returns the "LOG_LEVEL" environment variable.
// Examples: "DEBUG, "INFO", "WARN", "ERROR".
func (s *logSettings) GetLogLevel() string {
	return s.LogLevel
}
