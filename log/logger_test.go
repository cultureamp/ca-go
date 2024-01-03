package log

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggerLevels(t *testing.T) {
	defer os.Unsetenv("LOG_LEVEL")

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
			os.Setenv("LOG_LEVEL", tC.logLevel)

			config := newLoggerConfig()
			assert.NotNil(t, config)

			logger := newDefaultLogger(config)
			assert.NotNil(t, logger)

			level := logger.impl.GetLevel()
			assert.Equal(t, tC.expectedLevel, level, tC.desc)
		})
	}
}
