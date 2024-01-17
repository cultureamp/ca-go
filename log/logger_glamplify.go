package log

// Deprecated: Fields is a drop in replacement for glamplify log.Fields in logging statements.
type Fields map[string]interface{}

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
func (l *LegacyLogger) Debug(event string, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Debug(event).
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Info emits a new log message with info level.
func (l *LegacyLogger) Info(event string, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Info(event).
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Warn emits a new log message with warn level.
func (l *LegacyLogger) Warn(event string, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Warn(event).
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Error emits a new log message with error level.
func (l *LegacyLogger) Error(event string, err error, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Error(event, err).
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Fatal emits a new log message with fatal level.
// The os.Exit(1) function is called immediately afterwards which terminates the program.
func (l *LegacyLogger) Fatal(event string, err error, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Fatal(event, err).
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Panic emits a new log message with panic level.
// The panic() function is called immediately afterwards, which stops the ordinary flow of a goroutine.
func (l *LegacyLogger) Panic(event string, err error, fields ...Fields) {
	props := Fields{}
	props = props.Merge(fields...)

	l.impl.Panic(event, err).
		LegacyFields("properties", props).
		Send()
}

// Deprecated: Merge combines multiple legacy log.Fields together.
func (fields Fields) Merge(other ...Fields) Fields {
	merged := Fields{}

	for k, v := range fields {
		merged[k] = v
	}

	for _, f := range other {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}
