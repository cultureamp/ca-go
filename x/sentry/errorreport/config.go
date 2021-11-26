package errorreport

import (
	"github.com/getsentry/sentry-go"
)

type config struct {
	environment string
	dsn         string
	release     string
	debug       bool

	buildNumber string
	branch      string
	commit      string
	farm        string

	beforeFilter SentryBeforeFilter
	transport    sentry.Transport

	sentryOpts sentry.ClientOptions
	connected  bool
}

type Option func(c *config)

// SentryBeforeFilter is executed before a Sentry event is sent. It allows attributes
// of the event to be modified. The event can be discarded by returning nil.
type SentryBeforeFilter func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event

func WithEnvironment(env string) Option {
	return func(c *config) {
		c.environment = env
	}
}

func WithDSN(dsn string) Option {
	return func(c *config) {
		c.dsn = dsn
	}
}

func WithRelease(release string) Option {
	return func(c *config) {
		c.release = release
	}
}

func WithDebug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// WithBeforeFilter configures a function that will be called before an
// error is reported. This can be used to filter out certain errors from
// being reported.
func WithBeforeFilter(filter SentryBeforeFilter) Option {
	return func(c *config) {
		c.beforeFilter = filter
	}
}

func WithTransport(transport sentry.Transport) Option {
	return func(c *config) {
		c.transport = transport
	}
}

func WithServerlessTransport() Option {
	return WithTransport(sentry.NewHTTPSyncTransport())
}

func WithBuildDetails(farm, buildNumber, branch, commit string) Option {
	return func(c *config) {
		c.farm = farm
		c.buildNumber = buildNumber
		c.branch = branch
		c.commit = commit
	}
}
