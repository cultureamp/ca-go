package log

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

const (
	AppVerDefault       = "1.0.0"
	AwsAccountIDDefault = "development"
	AppFarmDefault      = "local"
	LogLevelDefault     = "INFO"
)

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

	quiet := getEnvBool(LogQuietModeEnv, isTestMode())
	consoleWriter := getEnvBool(LogConsoleWriterEnv, isTestMode())

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
	}
}

func (c *Config) isLocal() bool {
	return c.Farm == AppFarmDefault && c.AwsAccountID == AwsAccountIDDefault
}

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
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
