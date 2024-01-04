package env

const (

	// *** AuthZ Environment Variables ***.
	AuthzClientTimeoutEnv = "AUTHZ_CLIENT_TIMEOUT_IN_MS"
	AuthzCacheDurationEnv = "AUTHZ_CACHE_DURATION_IN_SEC"
	AuthzDialerTimeoutEnv = "AUTHZ_DIALER_TIMEOUT_IN_MS"
	AuthzTLSTimeoutEnv    = "AUTHZ_TLS_TIMEOUT_IN_MS"

	// *** AWS Environment Variables ***.
	AwsProfileEnv   = "AWS_PROFILE"
	AwsRegionEnv    = "AWS_REGION"
	AwsAccountIDEnv = "AWS_ACCOUNT_ID"
	AwsXrayEnv      = "XRAY_LOGGING"

	// *** Cache Environment Variables ***.
	CacheDurationEnv = "CACHE_DURATION_IN_SEC"

	// *** Datadog Environment Variables ***.
	DatadogAPIEnvVar   = "DD_API_KEY"
	DatadogLogEndpoint = "DD_LOG_ENDPOINT"
	DatadogEnv         = "DD_ENV"
	DatadogService     = "DD_SERVICE"
	DatadogVersion     = "DD_VERSION"
	DatadogAgentHost   = "DD_AGENT_HOST"
	DatadogStatsdPort  = "DD_DOGSTATSD_PORT"
	DatadogTimeout     = "DD_TIMEOUT"
	DatadogSite        = "DD_SITE"
	DatadogLogLevel    = "DD_LOG_LEVEL"

	// *** Log Environment Variables ***.
	LogLevel      = "LOG_LEVEL"
	LogOmitEmpty  = "LOG_OMITEMPTY"
	LogUseColours = "LOG_COLOURS"

	// *** Sentry Environment Variables ***.
	SentryDsnEnv              = "SENTRY_DSN"
	SentryFlushTimeoutInMsEnv = "SENTRY_FLUSH_TIMEOUT_IN_MS"

	// *** Global Environment Variables ***.
	AppNameEnv       = "APP"
	AppVerEnv        = "APP_VERSION"
	AppEnv           = "APP_ENV"
	AppFarmEnv       = "FARM"
	AppFarmLegacyEnv = "APP_ENV"
	ProductEnv       = "PRODUCT"
)
