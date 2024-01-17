package log

var defaultLogger *Logger = getInstance()

func getInstance() *Logger {
	setGlobalLogger()
	config := NewLoggerConfig()
	return NewLogger(config)
}

// Debug starts a new message with debug level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Debug(event string) *LoggerField {
	return defaultLogger.Debug(event)
}

// Info starts a new message with info level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Info(event string) *LoggerField {
	return defaultLogger.Info(event)
}

// Warn starts a new message with warn level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Warn(event string) *LoggerField {
	return defaultLogger.Warn(event)
}

// Error starts a new message with error level.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Error(event string, err error) *LoggerField {
	return defaultLogger.Error(event, err)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method, which terminates the program immediately.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Fatal(event string, err error) *LoggerField {
	return defaultLogger.Fatal(event, err)
}

// Panic starts a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
//
// You must call Details, Detailsf or Send on the returned event in order to send the event to the output.
func Panic(event string, err error) *LoggerField {
	return defaultLogger.Panic(event, err)
}