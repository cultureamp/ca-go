package log

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// Logger that implements the CA Logging standard.
type Logger struct {
	impl zerolog.Logger
}

var defaultLogger *Logger = getInstance()

func getInstance() *Logger {
	setGlobalLogger()
	config := newLoggerConfig()
	return newDefaultLogger(config)
}

func newDefaultLogger(config *LoggerConfig) *Logger {
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

	logger := &Logger{
		zerolog.
			New(os.Stdout).
			With().
			Str("app", config.AppName).
			Str("app_version", config.AppVersion).
			Str("aws_region", config.AwsRegion).
			Str("aws_account_id", config.AwsAccountID).
			Str("farm", config.Farm).
			Str("product", config.Product).
			Timestamp().
			Logger().
			Level(lvl),
	}

	return logger
}

func setGlobalLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "details"
	zerolog.LevelFieldName = "severity"
	zerolog.DurationFieldInteger = true
}

// Properties creates an Event to be used with the *Event.Dict method.
// Call usual field methods like Str, Int etc to add fields to this
// event and give it as argument the *Event.Dict method.
func Properties() *zerolog.Event {
	return zerolog.Dict()
}

// Debug starts a new message with debug level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Debug(ctx context.Context, event string) *zerolog.Event {
	return defaultLogger.debug(ctx, event)
}

func (l *Logger) debug(ctx context.Context, event string) *zerolog.Event {
	le := l.impl.Debug().Str("event", toSnakeCase(event))
	le = l.autoAddFields(ctx, le)
	return le
}

// Info starts a new message with info level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Info(ctx context.Context, event string) *zerolog.Event {
	return defaultLogger.info(ctx, event)
}

func (l *Logger) info(ctx context.Context, event string) *zerolog.Event {
	le := l.impl.Info().Str("event", toSnakeCase(event))
	le = l.autoAddFields(ctx, le)
	return le
}

// Warn starts a new message with warn level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Warn(ctx context.Context, event string) *zerolog.Event {
	return defaultLogger.warn(ctx, event)
}

func (l *Logger) warn(ctx context.Context, event string) *zerolog.Event {
	le := l.impl.Warn().Str("event", toSnakeCase(event))
	le = l.autoAddFields(ctx, le)
	return le
}

// Error starts a new message with error level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Error(ctx context.Context, event string, err error) *zerolog.Event {
	return defaultLogger.error(ctx, event, err)
}

func (l *Logger) error(ctx context.Context, event string, err error) *zerolog.Event {
	le := l.impl.Error().Err(err).Str("event", toSnakeCase(event))
	le = l.autoAddFields(ctx, le)
	return le
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Fatal(ctx context.Context, event string, err error) *zerolog.Event {
	return defaultLogger.fatal(ctx, event, err)
}

func (l *Logger) fatal(ctx context.Context, event string, err error) *zerolog.Event {
	le := l.impl.Fatal().Err(err).Str("event", toSnakeCase(event))
	le = l.autoAddFields(ctx, le)
	return le
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Panic(ctx context.Context, event string, err error) *zerolog.Event {
	return defaultLogger.panic(ctx, event, err)
}

func (l *Logger) panic(ctx context.Context, event string, err error) *zerolog.Event {
	le := l.impl.Panic().Err(err).Str("event", toSnakeCase(event))
	le = l.autoAddFields(ctx, le)
	return le
}

func (l *Logger) autoAddFields(ctx context.Context, le *zerolog.Event) *zerolog.Event {
	le = l.addRequestIDs(ctx, le)
	le = l.addAuthenticatedUserIDs(ctx, le)
	le = l.addSystem(le)
	return le
}

func (l *Logger) addRequestIDs(ctx context.Context, le *zerolog.Event) *zerolog.Event {
	ids, ok := RequestIDsFromContext(ctx)
	if ok {
		le = le.Str("trace_id", ids.TraceID).Str("request_id", ids.RequestID).Str("correlation_id", ids.CorrelationID)
	}
	return le
}

func (l *Logger) addAuthenticatedUserIDs(ctx context.Context, le *zerolog.Event) *zerolog.Event {
	ids, ok := AuthPayloadFromContext(ctx)
	if ok {
		le = le.Str("account_id", ids.CustomerAccountID).Str("user_id", ids.UserID).Str("real_user_id", ids.RealUserID)
	}
	return le
}

func (l *Logger) addSystem(le *zerolog.Event) *zerolog.Event {
	host, err := os.Hostname()
	if err == nil {
		le = le.Str("host", host)
	}

	le = le.Str("os", runtime.GOOS)

	_, file, line, ok := runtime.Caller(3)
	if ok {
		le = le.Str("loc", file+":"+strconv.Itoa(line))
	}

	return le
}
