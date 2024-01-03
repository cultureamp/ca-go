package log

import (
	"context"

	"github.com/kelseyhightower/envconfig"
)

type contextValueKey string

const envConfigKey = contextValueKey("env")

// EnvConfig must have fields listed in https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard
type EnvConfig struct {
	AppName      string `envconfig:"APP"`                      // The name of the application the log was generated from
	AppVersion   string `default:"0.0.0" split_words:"true"`   // The version of the application
	AwsRegion    string `split_words:"true"`                   // the AWS region this code is running in
	AwsAccountID string `split_words:"true"`                   // the AWS account ID this code is running in
	Farm         string `default:"local" envconfig:"FARM"`     // The name of the farm or where the code is running
	LogLevel     string `default:"INFO" envconfig:"LOG_LEVEL"` // The logging level
}

// EnvConfigFromContext returns the EnvConfig value embedded in the given context. Return a default EnvConfig if not exists
func EnvConfigFromContext(ctx context.Context) EnvConfig {
	var config EnvConfig
	config, ok := ctx.Value(envConfigKey).(EnvConfig)
	if !ok {
		envconfig.MustProcess("", &config)
	}
	return config
}

// ContextWithEnvConfig returns a new context with the given EnvConfig embedded as a value.
func ContextWithEnvConfig(ctx context.Context, envConfig EnvConfig) context.Context {
	ctx = context.WithValue(ctx, envConfigKey, envConfig)
	return ctx
}
