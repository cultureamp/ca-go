package log

import (
	"context"
	"errors"
)

func Example() {
	ctx := context.Background()

	logger := GetInstance(EnvConfig{
		AppName:    "test-app",
		AppVersion: "1.0.0",
		AwsRegion:  "us-east-1",
		Farm:       "test-farm",
	})

	logger.Debug(ctx, "event_details").
		Str("test-str", "str").
		Int("test-number", 1).
		Msg("initialise handler")

	logger = GetInstance(EnvConfig{
		AppName:    "test-app",
		AppVersion: "1.0.0",
		AwsRegion:  "us-east-1",
		Farm:       "test-farm",
	})

	err := errors.New("exception")
	logger.Error(ctx, "event_details", err).
		Str("test-str", "str").
		Int("test-number", 1).
		Msg("initialise handler")

	// Output:
	// {"level":"error","app":"test-app","app_version":"1.0.0","aws_region":"us-east-1","aws_account_id":"","farm":"test-farm","error":"exception","event":"event_details","test-str":"str","test-number":1,"message":"initialise handler"}
}
