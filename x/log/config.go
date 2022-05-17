package log

import (
	"context"

	"github.com/kelseyhightower/envconfig"
)

type contextValueKey string

const envConfigKey = contextValueKey("env")

type EnvConfig struct {
	AppName    string `env:"APP"`
	AppVersion string `env:"APP_VERSION" default:"1.0.0"`
	AwsRegion  string `env:"AWS_REGION"`
	Farm       string `env:"FARM" default:"local"`
}

func EnvConfigFromContext(ctx context.Context) (EnvConfig, bool) {
	config, ok := ctx.Value(envConfigKey).(EnvConfig)
	return config, ok
}

// ContextWithDefaultEnvConfig returns a new context with the default EnvConfig embedded as a value.
func ContextWithDefaultEnvConfig(ctx context.Context) context.Context {
	var envConfig EnvConfig
	envconfig.MustProcess("", &envConfig)
	return context.WithValue(ctx, envConfigKey, envConfig)
}

// ContextWithEnvConfig returns a new context with the given EnvConfig embedded as a value.
func ContextWithEnvConfig(ctx context.Context, envConfig EnvConfig) context.Context {
	return context.WithValue(ctx, envConfigKey, envConfig)
}

// ContextHasEnvConfig returns whether the given context contains an EnvConfig value.
func ContextHasEnvConfig(ctx context.Context) bool {
	_, ok := EnvConfigFromContext(ctx)
	return ok
}
