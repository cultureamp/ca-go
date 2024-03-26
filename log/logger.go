package log

import (
	"github.com/rs/zerolog"
)

// standardLogger that implements the CA Logging standard.
type standardLogger struct {
	impl   zerolog.Logger
	config *Config
}

func NewLogger(config *Config) *standardLogger {
	lvl := config.severityToLevel()
	writer := config.getWriter()

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

// Debug starts a new message with debug level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Debug(event string) *Property {
	le := l.impl.Debug().Str("event", toSnakeCase(event))
	return newLoggerProperty(le)
}

// Info starts a new message with info level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Info(event string) *Property {
	le := l.impl.Info().Str("event", toSnakeCase(event))
	return newLoggerProperty(le)
}

// Warn starts a new message with warn level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Warn(event string) *Property {
	le := l.impl.Warn().Str("event", toSnakeCase(event))
	return newLoggerProperty(le)
}

// Error starts a new message with error level.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Error(event string, err error) *Property {
	le := l.impl.Error()
	le.Dict("error", zerolog.Dict().
		Stack().
		Err(err),
	).Str("event", toSnakeCase(event))
	return newLoggerProperty(le).WithSystemTracing()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Fatal(event string, err error) *Property {
	le := l.impl.Fatal()
	le.Dict("error", zerolog.Dict().
		Stack().
		Err(err),
	).Str("event", toSnakeCase(event))
	return newLoggerProperty(le).WithSystemTracing()
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Msg or Send on the returned event in order to send the event to the output.
func (l *standardLogger) Panic(event string, err error) *Property {
	le := l.impl.Panic()
	le.Dict("error", zerolog.Dict().
		Stack().
		Err(err),
	).Str("event", toSnakeCase(event))
	return newLoggerProperty(le).WithSystemTracing()
}
