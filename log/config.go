package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
)

const (
	AppVerDefault       = "1.0.0"
	AwsAccountIDDefault = "development"
	AppFarmDefault      = "local"
	LogLevelDefault     = "INFO"
)

type timeNowFunc = func() time.Time

// Config contains logging configuration values.
type Config struct {
	// Mandatory fields listed in https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard
	AppName    string // The name of the application the log was generated from
	AppVersion string // The version of the application

	AwsRegion    string // the AWS region this code is running in
	AwsAccountID string // the AWS account ID this code is running in
	Product      string // performance, engagmentment, etc.
	Farm         string // The name of the farm or where the code is running

	LogLevel      string // The logging level
	Quiet         bool   // Are we running inside tests and we should be quiet?
	ConsoleWriter bool   // If Farm=local use key-value colour console logging
	ConsoleColour bool   // If ConsoleWriter=true then enable/disable colour

	TimeNow timeNowFunc // Defaults to "time.Now" but useful to set in tests
}

// NewLoggerConfig creates a new configuration based on environment variables
// which can easily be reset before passing to NewLogger().
func NewLoggerConfig() *Config {
	appName := os.Getenv(AppNameEnv)
	if appName == "" {
		appName = os.Getenv(AppNameLeagcyEnv)
	}

	appVersion := os.Getenv(AppVerEnv)
	if appVersion == "" {
		appVersion = AppVerDefault
	}

	awsRegion := os.Getenv(AwsRegionEnv)
	product := os.Getenv(ProductEnv)

	awsAccountID := os.Getenv(AwsAccountIDEnv)
	if awsAccountID == "" {
		awsAccountID = AwsAccountIDDefault
	}

	farm := os.Getenv(AppFarmEnv)
	if farm == "" {
		farm = os.Getenv(AppFarmLegacyEnv)
		if farm == "" {
			farm = AppFarmDefault
		}
	}

	logLevel := os.Getenv(LogLevelEnv)
	if logLevel == "" {
		logLevel = LogLevelDefault
	}

	quiet := getEnvBool(LogQuietModeEnv, false)
	consoleWriter := getEnvBool(LogConsoleWriterEnv, false)
	consoleColour := getEnvBool(LogConsoleWriterEnv, false)

	return &Config{
		LogLevel:      logLevel,
		AppName:       appName,
		AppVersion:    appVersion,
		AwsRegion:     awsRegion,
		AwsAccountID:  awsAccountID,
		Product:       product,
		Farm:          farm,
		Quiet:         quiet,
		ConsoleWriter: consoleWriter,
		ConsoleColour: consoleColour,
		TimeNow:       time.Now,
	}
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

func (c *Config) shouldProcess() error {
	if c.AppName == "" {
		return errors.Errorf("config.AppName is empty - missing APP or APP_NAME environment variable?")
	}

	if c.AwsRegion == "" {
		return errors.Errorf("config.AwsRegion is empty - missing AWS_REGION environment variable?")
	}

	if c.Product == "" {
		return errors.Errorf("config.Product is empty - missing PRODUCT environment variable?")
	}

	return nil
}

func (c *Config) mustProcess() {
	err := c.shouldProcess()
	if err != nil {
		// panics if mandatory env vars are not set
		panic(err)
	}
}

// GetBool gets the environment variable for 'key' if present, otherwise returns 'fallback'.
func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return b
}
