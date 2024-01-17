package log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger that implements the CA Logging standard.
type Logger struct {
	impl zerolog.Logger
}

func NewLogger(config *LoggerConfig) *Logger {
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

	// Default to Stdout, but if running in TestMode then set the logger to silently NoOp
	var writer io.Writer = os.Stdout
	if config.Quiet {
		writer = io.Discard
	}

	impl := zerolog.
		New(writer).
		With().
		Str("app", config.AppName).
		Str("app_version", config.AppVersion).
		Str("aws_region", config.AwsRegion).
		Str("aws_account_id", config.AwsAccountID).
		Str("farm", config.Farm).
		Str("product", config.Product).
		Timestamp().
		Logger().
		Level(lvl)

	return &Logger{impl: impl}
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
func (l *Logger) Debug(event string) *Property {
	le := l.impl.Debug().Str("event", toSnakeCase(event))
	return &Property{le}
}

// Info starts a new message with info level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *Logger) Info(event string) *Property {
	le := l.impl.Info().Str("event", toSnakeCase(event))
	return &Property{le}
}

// Warn starts a new message with warn level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *Logger) Warn(event string) *Property {
	le := l.impl.Warn().Str("event", toSnakeCase(event))
	return &Property{le}
}

// Error starts a new message with error level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *Logger) Error(event string, err error) *Property {
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
func (l *Logger) Fatal(event string, err error) *Property {
	le := l.impl.Fatal().Err(err).Str("event", toSnakeCase(event))
	return &Property{le}
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *Logger) Panic(event string, err error) *Property {
	le := l.impl.Panic().Err(err).Str("event", toSnakeCase(event))
	return &Property{le}
}
