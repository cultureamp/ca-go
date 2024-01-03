package log

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/rs/zerolog"
)

type Logger struct {
	impl zerolog.Logger
}

var (
	defaultLogger *Logger
	dlOnce        sync.Once
)

func GetInstance(config EnvConfig) *Logger {
	dlOnce.Do(func() {
		setGlobalLogger(config)

		defaultLogger = newDefaultLogger(config)
	})
	return defaultLogger
}

func newDefaultLogger(config EnvConfig) *Logger {
	logger := &Logger{zerolog.New(os.Stdout)}
	logger.setupFormatter(config.Farm)

	logger.impl.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.
			Str("app", config.AppName).
			Str("app_version", config.AppVersion).
			Str("aws_region", config.AwsRegion).
			Str("aws_account_id", config.AwsAccountID).
			Str("farm", config.Farm)
	})

	return logger
}

func setGlobalLogger(config EnvConfig) {
	zerolog.TimeFieldFormat = time.RFC3339

	// Default level for this example is info, unless debug flag is present
	switch config.LogLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// ToSnakeCase returns a new string in the format word_word
func toSnakeCase(s string) string {
	var sb strings.Builder

	in := []rune(strings.TrimSpace(s))
	for i, r := range in {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(in[i-1]) {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			if unicode.IsSpace(r) {
				if !unicode.IsSpace(in[i-1]) {
					sb.WriteRune('_')
				}
			} else {
				sb.WriteRune(r)
			}
		}
	}

	return sb.String()
}

// SetupFormatter decides the output formatter based on the environment where the app is running on.
// It uses text formatter with color if you run the app locally,
// while using json formatter if it's running on the cloud.
func (l *Logger) setupFormatter(farm string) {
	if farm == "local" {
		l.impl = l.impl.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		l.impl = l.impl.Output(os.Stdout)
	}
}

func (l *Logger) Debug(ctx context.Context, event string) *zerolog.Event {
	lg := l.impl.Debug().Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	return lg
}

func (l *Logger) Info(ctx context.Context, event string) *zerolog.Event {
	lg := l.impl.Info().Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	return lg
}

func (l *Logger) Warn(ctx context.Context, event string) *zerolog.Event {
	lg := l.impl.Warn().Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	return lg
}

func (l *Logger) Error(ctx context.Context, event string, err error) *zerolog.Event {
	lg := l.impl.Error().Err(err).Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	return lg
}

func (l *Logger) Fatal(ctx context.Context, event string, err error) *zerolog.Event {
	lg := l.impl.Fatal().Err(err).Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	return lg
}

func (l *Logger) addRequestIDs(ctx context.Context, lg *zerolog.Event) *zerolog.Event {
	ids, ok := RequestIDsFromContext(ctx)
	if ok {
		lg = lg.Str("trace_id", ids.TraceID).Str("request_id", ids.RequestID).Str("correlation_id", ids.CorrelationID)
	}
	return lg
}

func (l *Logger) addAuthenticatedUserIDs(ctx context.Context, lg *zerolog.Event) *zerolog.Event {
	ids, ok := AuthUserIDsFromContext(ctx)
	if ok {
		lg = lg.Str("account_id", ids.CustomerAccountID).Str("user_id", ids.UserID).Str("real_user_id", ids.RealUserID)
	}
	return lg
}
