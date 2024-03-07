package env

import (
	senv "github.com/caarlos0/env/v10"
)

// sentrySettings that drive behavior.
type sentrySettings struct {
	SentryDSN       string `env:"SENTRY_DSN"`
	SentryFlushInMs int    `env:"SENTRY_FLUSH_TIMEOUT_IN_MS" envDefault:"100"`
}

func newSentrySettings() *sentrySettings {
	settings := sentrySettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// GetSentryDSN returns the "SENTRY_DSN" environment variable.
func (s *sentrySettings) GetSentryDSN() string {
	return s.SentryDSN
}

// GetSentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func (s *sentrySettings) GetSentryFlushTimeoutInMs() int {
	return s.SentryFlushInMs
}
