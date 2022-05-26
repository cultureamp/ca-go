package errorreport

import (
	"fmt"

	"github.com/getsentry/sentry-go"
)

type config struct {
	environment string
	dsn         string
	release     string
	debug       bool

	tags map[string]string

	beforeFilter SentryBeforeFilter
	transport    sentry.Transport
}

func (c *config) tag(name string, value string) {
	if len(name) > 0 && len(value) > 0 {
		c.tags[name] = value
	}
}

// Option is a function type that can be provided to Configure to modify the
// behaviour of Sentry.
type Option func(c *config)

// SentryBeforeFilter is executed before a Sentry event is sent. It allows attributes
// of the event to be modified. The event can be discarded by returning nil.
type SentryBeforeFilter func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event

// WithEnvironment configures Sentry for the given environment, e.g. production-us.
// This is the name of the AWS account to which the application is deployed, and should be
// supplied to the application from the infrastructure via an environment variable. Environment
// names are defined in the Culture Amp CDK Constructs.
// This is a mandatory option.
func WithEnvironment(env string) Option {
	return func(c *config) {
		c.environment = env
	}
}

// WithDSN configures Sentry with the given DSN.
// This is a mandatory option.
func WithDSN(dsn string) Option {
	return func(c *config) {
		c.dsn = dsn
	}
}

// WithRelease formats the Sentry release with the given app name and version.
// This is a mandatory option.
func WithRelease(appName, appVersion string) Option {
	return func(c *config) {
		c.release = fmt.Sprintf("%s@%s", appName, appVersion)
	}
}

// WithDebug configures Sentry to log debug information.
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

// WithTransport configures an alternate transport for sending reports to
// Sentry.
func WithTransport(transport sentry.Transport) Option {
	return func(c *config) {
		c.transport = transport
	}
}

// WithServerlessTransport configures Sentry with the correct transport
// for serverless applications
func WithServerlessTransport() Option {
	return WithTransport(sentry.NewHTTPSyncTransport())
}

// WithBuildDetails configures Sentry to send build details along with
// error reports.
func WithBuildDetails(farm, buildNumber, branch, commit string) Option {
	return func(c *config) {
		c.tag("farm", farm)
		c.tag("build_number", buildNumber)
		c.tag("branch", branch)
		c.tag("commit", commit)
	}
}

// WithTag adds the specified name/value pair as a tag on every error report.
// Tags are used for grouping and searching error reports.
func WithTag(name string, value string) Option {
	return func(c *config) {
		c.tag(name, value)
	}
}

// WithTag adds multiple name/value pairs as tags on every error report. Tags
// are used for grouping and searching error reports.
func WithTags(tags map[string]string) Option {
	return func(c *config) {
		if tags == nil {
			return
		}

		for name, value := range tags {
			c.tag(name, value)
		}
	}
}
