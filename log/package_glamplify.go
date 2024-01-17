package log

var defaultLegacyLogger *LegacyLogger = getLegacyInstance()

func getLegacyInstance() *LegacyLogger {
	config := NewLoggerConfig()
	return NewLegacyLogger(config)
}

// Deprecated: LogDebug emits a new log message with debug level using the legacy glamplify logging interface.
func LogDebug(event string, fields ...Fields) {
	defaultLegacyLogger.Debug(event, fields...)
}

// Deprecated: Info emits a new log message with info level using the legacy glamplify logging interface.
func LogInfo(event string, fields ...Fields) {
	defaultLegacyLogger.Info(event, fields...)
}

// Deprecated: Warn emits a new log message with warn level using the legacy glamplify logging interface.
func LogWarn(event string, fields ...Fields) {
	defaultLegacyLogger.Warn(event, fields...)
}

// Deprecated: Error emits a new log message with error level using the legacy glamplify logging interface.
func LogError(event string, err error, fields ...Fields) {
	defaultLegacyLogger.Error(event, err, fields...)
}

// Deprecated: Fatal emits a new log message with fatal level using the legacy glamplify logging interface.
// The os.Exit(1) function is called immediately afterwards which terminates the program.
func LogFatal(event string, err error, fields ...Fields) {
	defaultLegacyLogger.Fatal(event, err, fields...)
}

// Deprecated: Panic emits a new log message with panic level using the legacy glamplify logging interface.
// The panic() function is called immediately afterwards, which stops the ordinary flow of a goroutine.
func LogPanic(event string, err error, fields ...Fields) {
	defaultLegacyLogger.Panic(event, err, fields...)
}
