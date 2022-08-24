package log

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvConfigFromContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectedCfg EnvConfig
	}{
		{
			name: "should overwrite default env config values if provided in context",
			ctx: context.WithValue(context.Background(), envConfigKey, EnvConfig{
				AppName:    "test-app",
				AppVersion: "1.0.0",
				AwsRegion:  "us-east-1",
				Farm:       "test-farm",
			}),
			expectedCfg: EnvConfig{
				AppName:    "test-app",
				AppVersion: "1.0.0",
				AwsRegion:  "us-east-1",
				Farm:       "test-farm",
			},
		},
		{
			name: "should have default env config values if not provided in context",
			ctx:  context.Background(),
			expectedCfg: EnvConfig{
				AppName:      "",
				AppVersion:   "0.0.0",
				AwsRegion:    "",
				AwsAccountId: "",
				Farm:         "local",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := EnvConfigFromContext(tt.ctx)
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
