package log

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLegacyLoggerLevels(t *testing.T) {
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
			logger := NewLegacyLogger(config)
			assert.NotNil(t, logger)

			level := logger.impl.impl.GetLevel()
			assert.Equal(t, tC.expectedLevel, level, tC.desc)
		})
	}
}

func TestLegacyLoggerMethods(t *testing.T) {
	ctx := context.Background()

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
			config.Quiet = false
			logger := NewLegacyLogger(config)
			assert.NotNil(t, logger)

			switch tC.logLevel {
			case "DEBUG":
				logger.Debug(ctx, "debug_event")
			case "INFO":
				logger.Info(ctx, "info_event")
			case "WARN":
				logger.Warn(ctx, "warn_event")
			case "ERROR":
				logger.Error(ctx, "error_event", errors.New("test error"))
			}
		})
	}

	// Output:
	// {"severity":"debug","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"debug_event","time":"2024-01-17T16:10:54+11:00"}
	// {"severity":"info","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"info_event","time":"2024-01-17T16:10:54+11:00"}
	// {"severity":"warn","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","event":"warn_event","time":"2024-01-17T16:10:54+11:00"}
	// {"severity":"error","app":"","app_version":"1.0.0","aws_region":"","aws_account_id":"local","farm":"local","product":"","error":"test error","event":"error_event","time":"2024-01-17T16:10:54+11:00"}
}
