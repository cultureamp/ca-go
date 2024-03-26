package log

import (
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(event string) *Property
	Info(event string) *Property
	Warn(event string) *Property
	Error(event string, err error) *Property
	Fatal(event string, err error) *Property
	Panic(event string, err error) *Property
}

var DefaultLogger Logger = getInstance()

func getInstance() *standardLogger {
	setGlobalLogger()
	config := NewLoggerConfig()
	return NewLogger(config)
}

func setGlobalLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "details"
	zerolog.ErrorFieldName = "error_message"
	zerolog.LevelFieldName = "severity"
	zerolog.DurationFieldInteger = true
	zerolog.ErrorStackMarshaler = logStackTracer
}

// Debug starts a new message with debug level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Debug(event string) *Property {
	return DefaultLogger.Debug(event)
}

// Info starts a new message with info level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Info(event string) *Property {
	return DefaultLogger.Info(event)
}

// Warn starts a new message with warn level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Warn(event string) *Property {
	return DefaultLogger.Warn(event)
}

// Error starts a new message with error level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Error(event string, err error) *Property {
	return DefaultLogger.Error(event, err)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Fatal(event string, err error) *Property {
	return DefaultLogger.Fatal(event, err)
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Panic(event string, err error) *Property {
	return DefaultLogger.Panic(event, err)
}
