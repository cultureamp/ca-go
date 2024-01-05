package env

const (

	// *** Global Environment Variables ***.
	AppNameEnv        = "APP"
	AppVerEnv         = "APP_VERSION"
	AppEnvironmentEnv = "APP_ENV"
	AppFarmEnv        = "FARM"
	AppFarmLegacyEnv  = "APP_ENV"
	ProductEnv        = "PRODUCT"

	// *** AWS Environment Variables ***.
	AwsProfileEnv   = "AWS_PROFILE"
	AwsRegionEnv    = "AWS_REGION"
	AwsAccountIDEnv = "AWS_ACCOUNT_ID"
	AwsXrayEnv      = "XRAY_LOGGING"

	// *** Log Environment Variables ***.
	LogLevelEnv = "LOG_LEVEL"

	// *** AuthZ Environment Variables ***.
	AuthzClientTimeoutEnv = "AUTHZ_CLIENT_TIMEOUT_IN_MS"
	AuthzCacheDurationEnv = "AUTHZ_CACHE_DURATION_IN_SEC"
	AuthzDialerTimeoutEnv = "AUTHZ_DIALER_TIMEOUT_IN_MS"
	AuthzTLSTimeoutEnv    = "AUTHZ_TLS_TIMEOUT_IN_MS"

	// *** Cache Environment Variables ***.
	CacheDurationEnv = "CACHE_DURATION_IN_SEC"

	// *** Datadog Environment Variables ***.
	DatadogAPIEnv         = "DD_API_KEY"
	DatadogLogEndpointEnv = "DD_LOG_ENDPOINT"
	DatadogEnvironmentEnv = "DD_ENV"
	DatadogServiceEnv     = "DD_SERVICE"
	DatadogVersionEnv     = "DD_VERSION"
	DatadogAgentHostEnv   = "DD_AGENT_HOST"
	DatadogStatsdPortEnv  = "DD_DOGSTATSD_PORT"
	DatadogTimeoutEnv     = "DD_TIMEOUT"
	DatadogSiteEnv        = "DD_SITE"
	DatadogLogLevelEnv    = "DD_LOG_LEVEL"

	// *** Sentry Environment Variables ***.
	SentryDsnEnv              = "SENTRY_DSN"
	SentryFlushTimeoutInMsEnv = "SENTRY_FLUSH_TIMEOUT_IN_MS"
)
