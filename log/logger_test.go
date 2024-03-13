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
			assert.NotNil(t, config)

			config.LogLevel = tC.logLevel
			config.Quiet = true
			logger := NewLogger(config)
			assert.NotNil(t, logger)

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
