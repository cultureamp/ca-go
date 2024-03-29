package log

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cultureamp/ca-go/x/request"
	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

// SetupFormatter decides the output formatter based on the environment where the app is running on.
// It uses text formatter with color if you run the app locally,
// while using json formatter if it's running on the cloud.
func (l *Logger) SetupFormatter(farm string) {
	if farm == "local" {
		l.Logger = l.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		l.Logger = l.Output(os.Stdout)
	}
}

func (l *Logger) WithCaller() *Logger {
	l.Logger = l.With().Caller().Logger()
	return l
}

func (l *Logger) WithTimestamp() *Logger {
	l.Logger = l.With().Timestamp().Logger()
	return l
}

func newDefaultLogger(config EnvConfig) *Logger {
	logger := &Logger{zerolog.New(os.Stdout).Hook(EventHook{})}
	logger.SetupFormatter(config.Farm)

	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.
			Str("app", config.AppName).
			Str("app_version", config.AppVersion).
			Str("aws_region", config.AwsRegion).
			Str("aws_account_id", config.AwsAccountID).
			Str("farm", config.Farm)
	})

	return logger
}

// NewFromCtx creates a new defaultLogger from a context, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromCtx(ctx context.Context) *Logger {
	config := EnvConfigFromContext(ctx)
	logger := newDefaultLogger(config)

	// RequestID - Set by clients externally to a random string in the X-Request-ID header (if missing the web gateway sets this to a new UUID4 string).
	// CorrelationID - A UUID4 contained in the X-Correlation-ID header (the web gateway always sets this to a new UUID4 string).
	reqIds, ok := request.RequestIDsFromContext(ctx)
	if ok {
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Str("request_id", reqIds.RequestID).
				Str("correlation_id", reqIds.CorrelationID)
		})
	}

	// customer - The aggregate id for the customer account that can be used to look up a specific customer account's details
	// user -  The aggregate id for a user that can be used to look up a specific user's details
	userIds, ok := request.AuthenticatedUserFromContext(ctx)
	if ok {
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Str("user_id", userIds.UserID).
				Str("real_user_id", userIds.RealUserID).
				Str("account_id", userIds.CustomerAccountID)
		})
	}

	return logger
}

// NewFromRequest creates a new defaultLogger from a http.Request, which should contain RequestScopedFields.
// If the context does not contain then, then this method will NOT add them in.
func NewFromRequest(r *http.Request) *Logger {
	return NewFromCtx(r.Context())
}
