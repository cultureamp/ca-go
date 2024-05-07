package log

import (
	"fmt"
	"io"
	"os"
	"time"

	senv "github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

const (
	AppVerDefault       = "1.0.0"
	AwsAccountIDDefault = "development"
	AppFarmDefault      = "local"
	LogLevelDefault     = "INFO"
	MissingDefault      = "<missing>"
)

type timeNowFunc = func() time.Time

// Config contains logging configuration values.
type Config struct {
	// Mandatory fields listed in https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard
	AppName    string `env:"APP"         envDefault:"<unknown>"` // The name of the application the log was generated from
	AppVersion string `env:"APP_VERSION" envDefault:"1.0.0"`     // The version of the application

	AwsRegion    string `env:"AWS_REGION"     envDefault:"dev"`         // the AWS region this code is running in
	AwsAccountID string `env:"AWS_ACCOUNT_ID" envDefault:"development"` // the AWS account ID this code is running in
	Product      string `env:"PRODUCT"        envDefault:"<unknown>"`   // performance, engagmentment, etc.
	Farm         string `env:"FARM"           envDefault:"local"`       // The name of the farm or where the code is running

	LogLevel      string `env:"LOG_LEVEL"      envDefault:"INFO"`  // The logging level
	Quiet         bool   `env:"QUIET_MODE"     envDefault:"false"` // Are we running inside tests and we should be quiet?
	ConsoleWriter bool   `env:"CONSOLE_WRITER" envDefault:"false"` // If ConsoleWriter=true then key-value pair output
	ConsoleColour bool   `env:"CONSOLE_COLOUR" envDefault:"false"` // If ConsoleWriter=true then enable/disable colour

	TimeNow timeNowFunc // Defaults to "time.Now" but useful to set in tests
}

// NewLoggerConfig creates a new configuration based on environment variables
// which can easily be reset before passing to NewLogger().
func NewLoggerConfig() (*Config, error) {
	c := Config{
		TimeNow: time.Now,
	}
	err := senv.Parse(&c)
	return &c, err
}

func (c *Config) isLocal() bool {
	return c.Farm == AppFarmDefault && c.AwsAccountID == AwsAccountIDDefault
}

func (c *Config) Level() zerolog.Level {
	return c.ToLevel(c.LogLevel)
}

func (c *Config) ToLevel(logLevel string) zerolog.Level {
	var lvl zerolog.Level
	switch logLevel {
	case "DEBUG", "Debug", "debug":
		lvl = zerolog.DebugLevel
	case "WARN", "Warn", "warn":
		lvl = zerolog.WarnLevel
	case "ERROR", "Error", "error":
		lvl = zerolog.ErrorLevel
	case "FATAL", "Fatal", "fatal":
		lvl = zerolog.FatalLevel
	case "PANIC", "Panic", "panic":
		lvl = zerolog.PanicLevel
	default:
		lvl = zerolog.InfoLevel
	}

	return lvl
}

func (c *Config) getWriter() io.Writer {
	// Default to Stdout, but if running in QuietMode then set the logger to silently NoOp
	var writer io.Writer = os.Stdout
	if c.Quiet {
		writer = io.Discard
	}

	// NOTE: only allow ConsoleWriter to be configured if we are NOT production
	// as the ConsoleWriter is NOT performant and should just be used for local only
	if c.isLocal() && c.ConsoleWriter {
		writer = zerolog.ConsoleWriter{
			Out:             writer,
			TimeFormat:      time.RFC3339,
			NoColor:         !c.ConsoleColour,
			FormatMessage:   c.formatMessage,
			FormatTimestamp: c.formatTimestamp,
		}
	}

	return writer
}

func (c *Config) formatMessage(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("event=\"%s\"", i)
}

func (c *Config) formatTimestamp(i interface{}) string {
	if i == nil {
		return "nil"
	}
	timeString, ok := i.(string)
	if !ok {
		return "nil"
	}
	return timeString
}
