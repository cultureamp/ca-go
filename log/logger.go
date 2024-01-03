package log

import (
	"context"
	"os"
	"sync"
	"time"

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
	return l.impl.Debug().Str("event", ToSnakeCase(event))
}

func (l *Logger) Info(ctx context.Context, event string) *zerolog.Event {
	return l.impl.Info().Str("event", ToSnakeCase(event))
}

func (l *Logger) Warn(ctx context.Context, event string) *zerolog.Event {
	return l.impl.Warn().Str("event", ToSnakeCase(event))
}

func (l *Logger) Error(ctx context.Context, event string, err error) *zerolog.Event {
	return l.impl.Error().Err(err).Str("event", ToSnakeCase(event))
}

func (l *Logger) Fatal(ctx context.Context, event string, err error) *zerolog.Event {
	return l.impl.Fatal().Err(err).Str("event", ToSnakeCase(event))
}
