package log

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// Logger interface used for mock testing.
type Logger interface {
	Debug(event string) *Property
	Info(event string) *Property
	Warn(event string) *Property
	Error(event string, err error) *Property
	Fatal(event string, err error) *Property
	Panic(event string, err error) *Property

	Child(options ...LoggerOption) Logger
}

// DefaultLogger is the package level default implementation used by all package level methods.
// Package level methods are provided for ease of use.
// For testing you can replace the DefaultLogger with your own mock:
//
// DefaultLogger = newmockLogger().
var DefaultLogger Logger = nil

// Debug starts a new message with debug level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Debug(event string) *Property {
	mustHaveDefaultLogger()

	return DefaultLogger.Debug(event)
}

// Info starts a new message with info level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Info(event string) *Property {
	mustHaveDefaultLogger()

	return DefaultLogger.Info(event)
}

// Warn starts a new message with warn level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Warn(event string) *Property {
	mustHaveDefaultLogger()

	return DefaultLogger.Warn(event)
}

// Error starts a new message with error level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Error(event string, err error) *Property {
	mustHaveDefaultLogger()

	return DefaultLogger.Error(event, err)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Fatal(event string, err error) *Property {
	mustHaveDefaultLogger()

	return DefaultLogger.Fatal(event, err)
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Panic(event string, err error) *Property {
	mustHaveDefaultLogger()

	return DefaultLogger.Panic(event, err)
}

// DefaultOptions adds global properties to the DefaultLogger.
// Note: This creates a new DefaultLogger.
func DefaultOptions(options ...LoggerOption) {
	mustHaveDefaultLogger()

	// Update the DefaultLogger - not thread safe, but should be ok
	DefaultLogger = DefaultLogger.Child(options...)
}

// FromContext returns the Logger associated with the ctx. If not logger
// is associated, then a new logger is created and added to the context.
func FromContext(ctx context.Context) (context.Context, Logger, error) {
	if l, ok := ctx.Value(ctxLoggerKey{}).(Logger); ok {
		return ctx, l, nil
	}

	config, err := NewLoggerConfig()
	if err != nil {
		return ctx, nil, err
	}
	l := NewLogger(config)
	ctx = l.WithContext(ctx)
	return ctx, l, nil
}

func mustHaveDefaultLogger() {
	if DefaultLogger == nil {
		setGlobalLogger()
		config, _ := NewLoggerConfig()
		DefaultLogger = NewLogger(config)
	}
}

func setGlobalLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "details"
	zerolog.LevelFieldName = "severity"
	zerolog.DurationFieldInteger = true
	zerolog.ErrorStackMarshaler = logStackTracer
}
