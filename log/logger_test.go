package log

import (
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerLevels(t *testing.T) {
	testCases := []struct {
		desc          string
		logLevel      string
		expectedLevel zerolog.Level
	}{
		{
			desc:          "debug level test 1",
			logLevel:      "DEBUG",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			desc:          "debug level test 2",
			logLevel:      "Debug",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			desc:          "debug level test 3",
			logLevel:      "debug",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			desc:          "info level test 1",
			logLevel:      "INFO",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			desc:          "info level test 2",
			logLevel:      "Info",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			desc:          "info level test 3",
			logLevel:      "info",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			desc:          "warn level test 1",
			logLevel:      "WARN",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			desc:          "warn level test 2",
			logLevel:      "Warn",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			desc:          "warn level test 3",
			logLevel:      "warn",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			desc:          "error level test 1",
			logLevel:      "ERROR",
			expectedLevel: zerolog.ErrorLevel,
		},
		{
			desc:          "error level test 2",
			logLevel:      "Error",
			expectedLevel: zerolog.ErrorLevel,
		},
		{
			desc:          "error level test 3",
			logLevel:      "error",
			expectedLevel: zerolog.ErrorLevel,
		},
		{
			desc:          "fatal level test 1",
			logLevel:      "FATAL",
			expectedLevel: zerolog.FatalLevel,
		},
		{
			desc:          "fatal level test 2",
			logLevel:      "Fatal",
			expectedLevel: zerolog.FatalLevel,
		},
		{
			desc:          "fatal level test 3",
			logLevel:      "fatal",
			expectedLevel: zerolog.FatalLevel,
		},
		{
			desc:          "panic level test 1",
			logLevel:      "PANIC",
			expectedLevel: zerolog.PanicLevel,
		},
		{
			desc:          "panic level test 2",
			logLevel:      "Panic",
			expectedLevel: zerolog.PanicLevel,
		},
		{
			desc:          "panic level test 3",
			logLevel:      "panic",
			expectedLevel: zerolog.PanicLevel,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			config := NewLoggerConfig()
			assert.NotNil(t, config)

			config.LogLevel = tC.logLevel
			config.Quiet = true
			logger := NewLogger(config)
			assert.NotNil(t, logger)

			level := logger.impl.GetLevel()
			assert.Equal(t, tC.expectedLevel, level, tC.desc)
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	testCases := []struct {
		desc          string
		logLevel      string
		expectedLevel zerolog.Level
	}{
		{
			desc:          "debug level test",
			logLevel:      "DEBUG",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			desc:          "info level test",
			logLevel:      "INFO",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			desc:          "warn level test",
			logLevel:      "WARN",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			desc:          "error level test",
			logLevel:      "ERROR",
			expectedLevel: zerolog.ErrorLevel,
		},
		{
			desc:          "fatal level test",
			logLevel:      "FATAL",
			expectedLevel: zerolog.FatalLevel,
		},
		{
			desc:          "panic level test",
			logLevel:      "PANIC",
			expectedLevel: zerolog.PanicLevel,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			config := NewLoggerConfig()
			config.AppName = "logger-test"
			config.AwsRegion = "def"
			config.Product = "cago"

			assert.NotNil(t, config)

			config.LogLevel = tC.logLevel
			config.Quiet = true
			logger := NewLogger(config)
			assert.NotNil(t, logger)

			assert.True(t, logger.Enabled(tC.logLevel))
			switch tC.logLevel {
			case "DEBUG":
				logger.Debug("debug_event").Send()
			case "INFO":
				logger.Info("info_event").Send()
			case "WARN":
				logger.Warn("warn_event").Send()
			case "ERROR":
				logger.Error("error_event", errors.New("test error")).Send()
			}
		})
	}
}
