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
	lg := l.impl.Debug().Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	lg = l.addSystem(ctx, lg)
	return lg
}

// Info starts a new message with info level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Info(ctx context.Context, event string) *zerolog.Event {
	return defaultLogger.Info(ctx, event)
}

func (l *Logger) Info(ctx context.Context, event string) *zerolog.Event {
	lg := l.impl.Info().Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	lg = l.addSystem(ctx, lg)
	return lg
}

// Warn starts a new message with warn level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Warn(ctx context.Context, event string) *zerolog.Event {
	return defaultLogger.Warn(ctx, event)
}

func (l *Logger) Warn(ctx context.Context, event string) *zerolog.Event {
	lg := l.impl.Warn().Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	lg = l.addSystem(ctx, lg)
	return lg
}

// Error starts a new message with error level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Error(ctx context.Context, event string, err error) *zerolog.Event {
	return defaultLogger.Error(ctx, event, err)
}

func (l *Logger) Error(ctx context.Context, event string, err error) *zerolog.Event {
	lg := l.impl.Error().Err(err).Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	lg = l.addSystem(ctx, lg)
	return lg
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Fatal(ctx context.Context, event string, err error) *zerolog.Event {
	return defaultLogger.Fatal(ctx, event, err)
}

func (l *Logger) Fatal(ctx context.Context, event string, err error) *zerolog.Event {
	lg := l.impl.Fatal().Err(err).Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	lg = l.addSystem(ctx, lg)
	return lg
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func Panic(ctx context.Context, event string, err error) *zerolog.Event {
	return defaultLogger.Panic(ctx, event, err)
}

func (l *Logger) Panic(ctx context.Context, event string, err error) *zerolog.Event {
	lg := l.impl.Panic().Err(err).Str("event", toSnakeCase(event))
	lg = l.addRequestIDs(ctx, lg)
	lg = l.addAuthenticatedUserIDs(ctx, lg)
	lg = l.addSystem(ctx, lg)
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
	ids, ok := AuthPayloadFromContext(ctx)
	if ok {
		lg = lg.Str("account_id", ids.CustomerAccountID).Str("user_id", ids.UserID).Str("real_user_id", ids.RealUserID)
	}
	return lg
}

func (l *Logger) addSystem(ctx context.Context, lg *zerolog.Event) *zerolog.Event {
	host, err := os.Hostname()
	if err == nil {
		lg = lg.Str("host", host)
	}

	lg = lg.Str("os", runtime.GOOS)

	_, file, line, ok := runtime.Caller(3)
	if ok {
		lg = lg.Str("loc", file+":"+strconv.Itoa(line))
	}

	return lg
}
