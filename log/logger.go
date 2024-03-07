package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// standardLogger that implements the CA Logging standard.
type standardLogger struct {
	impl   zerolog.Logger
	config *Config
}

func NewLogger(config *Config) *standardLogger {
	var lvl zerolog.Level
	switch config.LogLevel {
	case "DEBUG":
		lvl = zerolog.DebugLevel
	case "WARN":
		lvl = zerolog.WarnLevel
	case "ERROR":
		lvl = zerolog.ErrorLevel
	case "FATAL":
		lvl = zerolog.FatalLevel
	case "PANIC":
		lvl = zerolog.PanicLevel
	default:
		lvl = zerolog.InfoLevel
	}

	// Default to Stdout, but if running in QuietMode then set the logger to silently NoOp
	var writer io.Writer = os.Stdout
	if config.Quiet {
		writer = io.Discard
	}

	if config.isLocal() {
		// NOTE: only allow ConsoleWriter to be configured if we are NOT production
		// as the ConsoleWriter is NOT performant and should just be used for local only
		if config.ConsoleWriter {
			writer = zerolog.ConsoleWriter{
				Out:        writer,
				TimeFormat: time.RFC3339,
				NoColor:    !config.ConsoleColour,
				FormatMessage: func(i interface{}) string {
					if i == nil {
						return ""
					}
					return fmt.Sprintf("event=\"%s\"", i)
				},
			}
		}
	}

	impl := zerolog.
		New(writer).
		Level(lvl).
		With().
		Str("app", config.AppName).
		Str("app_version", config.AppVersion).
		Str("aws_region", config.AwsRegion).
		Str("aws_account_id", config.AwsAccountID).
		Str("farm", config.Farm).
		Str("product", config.Product).
		Logger()

	// We have our own Timestamp hook so that we can mock in tests
	impl = impl.Hook(&timestampHook{config: config})

	return &standardLogger{
		impl:   impl,
		config: config,
	}
}

func setGlobalLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "details"
	zerolog.LevelFieldName = "severity"
	zerolog.DurationFieldInteger = true
}

// Debug starts a new message with debug level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Debug(event string) *Property {
	le := l.impl.Debug().Str("event", toSnakeCase(event))
	return &Property{le}
}

// Info starts a new message with info level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Info(event string) *Property {
	le := l.impl.Info().Str("event", toSnakeCase(event))
	return &Property{le}
}

// Warn starts a new message with warn level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Warn(event string) *Property {
	le := l.impl.Warn().Str("event", toSnakeCase(event))
	return &Property{le}
}

// Error starts a new message with error level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Error(event string, err error) *Property {
	le := l.impl.Error().
		Err(err).
		Str("event", toSnakeCase(event))
	fields := &Property{le}
	return fields.withFullStack()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Fatal(event string, err error) *Property {
	le := l.impl.Fatal().Err(err).Str("event", toSnakeCase(event))
	return &Property{le}
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Panic(event string, err error) *Property {
	le := l.impl.Panic().Err(err).Str("event", toSnakeCase(event))
	return &Property{le}
}
