package log

import "context"

var defaultLegacyLogger *LegacyLogger = getLegacyInstance()

func getLegacyInstance() *LegacyLogger {
	config := NewLoggerConfig()
	return NewLegacyLogger(config)
}

// Deprecated: GlamplifyDebug emits a new log message with debug level using the legacy glamplify logging interface.
func GlamplifyDebug(ctx context.Context, event string, fields ...Fields) {
	defaultLegacyLogger.Debug(ctx, event, fields...)
}

// Deprecated: Info emits a new log message with info level using the legacy glamplify logging interface.
func GlamplifyInfo(ctx context.Context, event string, fields ...Fields) {
	defaultLegacyLogger.Info(ctx, event, fields...)
}

// Deprecated: Warn emits a new log message with warn level using the legacy glamplify logging interface.
func GlamplifyWarn(ctx context.Context, event string, fields ...Fields) {
	defaultLegacyLogger.Warn(ctx, event, fields...)
}

// Deprecated: Error emits a new log message with error level using the legacy glamplify logging interface.
func GlamplifyError(ctx context.Context, event string, err error, fields ...Fields) {
	defaultLegacyLogger.Error(ctx, event, err, fields...)
}

// Deprecated: Fatal emits a new log message with fatal level using the legacy glamplify logging interface.
// The os.Exit(1) function is called immediately afterwards which terminates the program.
func GlamplifyFatal(ctx context.Context, event string, err error, fields ...Fields) {
	defaultLegacyLogger.Fatal(ctx, event, err, fields...)
}

// Deprecated: Panic emits a new log message with panic level using the legacy glamplify logging interface.
// The panic() function is called immediately afterwards, which stops the ordinary flow of a goroutine.
func GlamplifyPanic(ctx context.Context, event string, err error, fields ...Fields) {
	defaultLegacyLogger.Panic(ctx, event, err, fields...)
}
