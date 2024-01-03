package log

import (
	"os"
)

// LoggerConfig must have fields listed in https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard
type LoggerConfig struct {
	LogLevel     string // The logging level
	AppName      string // The name of the application the log was generated from
	AppVersion   string // The version of the application
	AwsRegion    string // the AWS region this code is running in
	AwsAccountID string // the AWS account ID this code is running in
	Product      string // performance, engagmentment, etc.
	Farm         string // The name of the farm or where the code is running
}

func newLoggerConfig() *LoggerConfig {
	appName := os.Getenv("APP")
	awsRegion := os.Getenv("AWS_REGION")
	product := os.Getenv("PRODUCT")

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "INFO"
	}

	awsAccountID, ok := os.LookupEnv("AWS_ACCOUNT_ID")
	if !ok {
		awsAccountID = "local"
	}

	farm, ok := os.LookupEnv("FARM")
	if !ok {
		farm = "local"
	}

	appVersion, ok := os.LookupEnv("APP_VERSION")
	if !ok {
		appVersion = "1.0.0"
	}

	return &LoggerConfig{
		LogLevel:     logLevel,
		AppName:      appName,
		AppVersion:   appVersion,
		AwsRegion:    awsRegion,
		AwsAccountID: awsAccountID,
		Product:      product,
		Farm:         farm,
	}
}
