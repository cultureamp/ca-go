package env

import (
	senv "github.com/caarlos0/env/v11"
)

// SentrySettings implements Sentry settings.
// This is an interface so that clients can mock out this behaviour in tests.
type SentrySettings interface {
	SentryDSN() string
	SentryFlushTimeoutInMs() int
}

// sentrySettings that drive behavior.
type sentrySettings struct {
	// These have to be public so that "github.com/caarlos0/env/v10" can populate them
	SentryDSNEnv       string `env:"SENTRY_DSN"`
	SentryFlushInMsEnv int    `env:"SENTRY_FLUSH_TIMEOUT_IN_MS" envDefault:"100"`
}

func newSentrySettings() *sentrySettings {
	settings := sentrySettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// SentryDSN returns the "SENTRY_DSN" environment variable.
func (s *sentrySettings) SentryDSN() string {
	return s.SentryDSNEnv
}

// SentryFlushTimeoutInMs returns the "SENTRY_FLUSH_TIMEOUT_IN_MS" environment variable.
// Default: 100.
func (s *sentrySettings) SentryFlushTimeoutInMs() int {
	return s.SentryFlushInMsEnv
}
