package errorreport

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/getsentry/sentry-go"
)

var sentryConfig *config

const (
	sentryTracingSubheading = "Culture Amp - Tracing"
)

func Configure(opts ...Option) error {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	var missingMandatoryConfigs []string

	if cfg.environment == "" {
		missingMandatoryConfigs = append(missingMandatoryConfigs, "environment")
	}

	if cfg.dsn == "" {
		missingMandatoryConfigs = append(missingMandatoryConfigs, "DSN")
	}

	if cfg.release == "" {
		missingMandatoryConfigs = append(missingMandatoryConfigs, "release")
	}

	if len(missingMandatoryConfigs) > 0 {
		return fmt.Errorf("mandatory fields missing: %s", strings.Join(missingMandatoryConfigs, ", "))
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

	cfg.sentryOpts = sentryOpts

	sentryConfig = cfg

	return nil
}

func Connect() error {
	if sentryConfig == nil {
		return errors.New("attempt to connect an unconfigured client")
	}

	if sentryConfig.connected {
		// Don't attempt to connect if not necessary.
		return nil
	}

	err := sentry.Init(sentryConfig.sentryOpts)
	if err != nil {
		return fmt.Errorf("initialise sentry: %w", err)
	}

	sentryConfig.connected = true

	// Add build information to the scope for all error reports.
	// This can't be done before we initialsie the Sentry client.
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("build_number", sentryConfig.buildNumber)
		scope.SetTag("branch", sentryConfig.branch)
		scope.SetTag("commit", sentryConfig.commit)
		scope.SetTag("farm", sentryConfig.farm)
	})

	return nil
}

func ReportError(ctx context.Context, err error) {
	hub := sentry.CurrentHub()

	hub.WithScope(func(scope *sentry.Scope) {
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
		hub.CaptureException(err)
	})
}
