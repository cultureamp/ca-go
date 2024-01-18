package log

import "context"

// Deprecated: LegacyLogger supports the legacy glamplify logging interface.
type LegacyLogger struct {
	impl *Logger
}

// Deprecated: NewLegacyLogger is deprecated but included here for easy migration from glamplify.
func NewLegacyLogger(config *Config) *LegacyLogger {
	impl := NewLogger(config)
	return &LegacyLogger{impl: impl}
}

// Deprecated: Debug emits a new log message with debug level.
func (l *LegacyLogger) Debug(_ context.Context, event string, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Debug(event).
		WithSystemTracing().
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Info emits a new log message with info level.
func (l *LegacyLogger) Info(_ context.Context, event string, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Info(event).
		WithSystemTracing().
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Warn emits a new log message with warn level.
func (l *LegacyLogger) Warn(_ context.Context, event string, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Warn(event).
		WithSystemTracing().
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Error emits a new log message with error level.
func (l *LegacyLogger) Error(_ context.Context, event string, err error, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Error(event, err).
		WithSystemTracing().
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Fatal emits a new log message with fatal level.
// The os.Exit(1) function is called immediately afterwards which terminates the program.
func (l *LegacyLogger) Fatal(_ context.Context, event string, err error, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Fatal(event, err).
		WithSystemTracing().
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Panic emits a new log message with panic level.
// The panic() function is called immediately afterwards, which stops the ordinary flow of a goroutine.
func (l *LegacyLogger) Panic(_ context.Context, event string, err error, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Panic(event, err).
		WithSystemTracing().
		LegacyFields("properties", props).
		Send()
}

