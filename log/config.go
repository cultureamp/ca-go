package log

import (
	"flag"
	"os"
	"strings"
)

// LoggerConfig contains logging configuration values.
type LoggerConfig struct {
	LogLevel string // The logging level
	Quiet    bool   // Are we running inside tests and we should be quiet?

	// Mandatory fields listed in https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard
	AppName      string // The name of the application the log was generated from
	AppVersion   string // The version of the application
	AwsRegion    string // the AWS region this code is running in
	AwsAccountID string // the AWS account ID this code is running in
	Product      string // performance, engagmentment, etc.
	Farm         string // The name of the farm or where the code is running
}

// NewLoggerConfig creates a new configuration based on environment variables
// which can easily be reset before passing to NewLogger().
func NewLoggerConfig() *LoggerConfig {
	appName, ok := os.LookupEnv("APP")
	if !ok || appName == "" {
		appName = os.Getenv("APP_NAME")
	}

	awsRegion := os.Getenv("AWS_REGION")
	product := os.Getenv("PRODUCT")

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok || logLevel == "" {
		logLevel = "INFO"
	}

	awsAccountID, ok := os.LookupEnv("AWS_ACCOUNT_ID")
	if !ok || awsAccountID == "" {
		awsAccountID = "local"
	}

	farm, ok := os.LookupEnv("FARM")
	if !ok || farm == "" {
		farm = "local"
	}

	appVersion, ok := os.LookupEnv("APP_VERSION")
	if !ok || appVersion == "" {
		appVersion = "1.0.0"
	}

	var quiet bool
	testMode := os.Getenv("QUIET_MODE")
	switch testMode {
	case "ON", "TRUE", "YES":
		quiet = true
	case "OFF", "FALSE", "NO":
		quiet = false
	default:
		quiet = isTestMode()
	}

	return &LoggerConfig{
		LogLevel:     logLevel,
		AppName:      appName,
		AppVersion:   appVersion,
		AwsRegion:    awsRegion,
		AwsAccountID: awsAccountID,
		Product:      product,
		Farm:         farm,
		Quiet:        quiet,
	}
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
