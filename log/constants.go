package log

const (

	// *** Global Environment Variables ***.
	AppNameEnv       = "APP"
	AppNameLeagcyEnv = "APP_NAME"
	AppVerEnv        = "APP_VERSION"
	AppFarmEnv       = "FARM"
	AppFarmLegacyEnv = "APP_ENV"
	ProductEnv       = "PRODUCT"

	// *** AWS Environment Variables ***.
	AwsProfileEnv   = "AWS_PROFILE"
	AwsRegionEnv    = "AWS_REGION"
	AwsAccountIDEnv = "AWS_ACCOUNT_ID"

	// *** Log Environment Variables ***.
	LogLevelEnv         = "LOG_LEVEL"
	LogQuietModeEnv     = "QUIET_MODE"
	LogConsoleWriterEnv = "CONSOLE_WRITER"
	LogConsoleColourEnv = "CONSOLE_COLOUR"
)
