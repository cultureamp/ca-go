package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvConfigFromContext(t *testing.T) {
	t.Setenv("APP", "")
	t.Setenv("APP_VERSION", "0.0.0")
	t.Setenv("AWS_REGION", "")
	t.Setenv("AWS_ACCOUNT_ID", "")
	t.Setenv("FARM", "local")

	tests := []struct {
		name        string
		env         *EnvConfig
		expectedCfg EnvConfig
	}{
		{
			name: "should overwrite default env config values if provided in context",
			env: &EnvConfig{
				AppName:    "test-app",
				AppVersion: "1.0.0",
				AwsRegion:  "us-east-1",
				Farm:       "test-farm",
			},
			expectedCfg: EnvConfig{
				AppName:    "test-app",
				AppVersion: "1.0.0",
				AwsRegion:  "us-east-1",
				Farm:       "test-farm",
			},
		},
		{
			name: "should have default env config values if not provided in context",
			expectedCfg: EnvConfig{
				AppName:      "",
				AppVersion:   "0.0.0",
				AwsRegion:    "",
				AwsAccountID: "",
				Farm:         "local",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.env != nil {
				ctx = context.WithValue(context.Background(), envConfigKey, *tt.env)
			}
			cfg := EnvConfigFromContext(ctx)
			assert.Equal(t, tt.expectedCfg, cfg)
		})
	}
}

func TestContextWithEnvConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  EnvConfig
	}{
		{
			name: "should have env config values in context",
			cfg: EnvConfig{
				AppName: "test-app",
				Farm:    "test-farm",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextWithEnvConfig(context.Background(), tt.cfg)
			assert.Equal(t, tt.cfg, ctx.Value(envConfigKey))
		})
	}
}
