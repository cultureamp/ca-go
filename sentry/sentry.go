package sentry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cultureamp/ca-go/request"
	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
)

const (
	sentryTracingSubheading = "Culture Amp - Tracing"
)

// Init initialises the Sentry client with the given options. It returns
// an error if mandatory options are not supplied.
func Init(opts ...Option) error {
	cfg := &config{
		tags: make(map[string]string),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	var missingMandatoryConfigs []string

	if cfg.environment == "" {
		missingMandatoryConfigs = append(missingMandatoryConfigs, "environment")
	}

	if cfg.dsn == "" && cfg.environment != "local" {
		missingMandatoryConfigs = append(missingMandatoryConfigs, "DSN")
	}

	if cfg.release == "" {
		missingMandatoryConfigs = append(missingMandatoryConfigs, "release")
	}

	if len(missingMandatoryConfigs) > 0 {
		return errors.Errorf("mandatory fields missing: %s", strings.Join(missingMandatoryConfigs, ", "))
	}

	sentryOpts := sentry.ClientOptions{
		Environment: cfg.environment,
		Dsn:         cfg.dsn,
		Release:     cfg.release,
		Debug:       cfg.debug,
	}

	if cfg.beforeFilter != nil {
		sentryOpts.BeforeSend = cfg.beforeFilter
	}

	if cfg.transport != nil {
		sentryOpts.Transport = cfg.transport
	}

	if err := sentry.Init(sentryOpts); err != nil {
		return errors.Errorf("initialise sentry: %w", err)
	}

	// Add build information to the scope for all error reports.
	// This can't be done before we initialise the Sentry client.
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		if cfg.tags != nil {
			scope.SetTags(cfg.tags)
		}
	})

	return nil
}

// ReportError reports an error to Sentry. It will attempt to
// extract request IDs and the authenticated user from the
// context.
func ReportError(ctx context.Context, err error) {
	scope := sentry.CurrentHub().PushScope()
	defer sentry.PopScope()

	addRequestFieldsToScope(ctx, scope)
	sentry.CaptureException(err)
}

// Decorate creates a new Sentry scope and adds the supplied tags. This allows
// for any errors that are generated in this scope to have additional
// information added to them. For example, this is useful when processing
// multiple event records in a batch. The return value is a function that should
// be passed to `defer` so the created scope is automatically popped.
func Decorate(tags map[string]string) func() {
	scope := sentry.CurrentHub().PushScope()
	scope.SetTags(tags)
	return sentry.PopScope
}

// GracefullyShutdown flushes the Sentry client buffered events with blocking for at most the given timeout.
func GracefullyShutdown(err interface{}, timeout time.Duration) {
	var newError error
	switch content := err.(type) {
	case string:
		newError = errors.New(content)
	case error:
		newError = content
	default:
		message := fmt.Sprintf("unknown panic: %T", content)
		newError = errors.New(message)
	}
	sentry.CurrentHub().Recover(newError)
	sentry.Flush(timeout)
}

func addRequestFieldsToScope(ctx context.Context, scope *sentry.Scope) {
	if authenticatedUser, ok := request.AuthenticatedUserFromContext(ctx); ok {
		scope.SetUser(sentry.User{
			ID: authenticatedUser.UserID,
		})

		scope.SetTag("customer", authenticatedUser.CustomerAccountID)
		scope.SetTag("user.real", authenticatedUser.RealUserID)
	}

	if requestIDs, ok := request.RequestIDsFromContext(ctx); ok {
		scope.SetTag("RequestID", requestIDs.RequestID)

		// add as a context as well for display below the stack trace
		scope.SetContext(sentryTracingSubheading, map[string]interface{}{
			"RequestID":     requestIDs.RequestID,
			"CorrelationID": requestIDs.CorrelationID,
		})
	}
}
