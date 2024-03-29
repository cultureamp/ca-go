package log_test

import (
	"context"

	"github.com/cultureamp/ca-go/x/log"
	"github.com/cultureamp/ca-go/x/request"
)

// This example illustrates how to use NewFromCtx to initialise a logger with values from context.
func Example() {
	ctx := context.Background()

	ctx = log.ContextWithEnvConfig(ctx, log.EnvConfig{
		AppName:    "test-app",
		AppVersion: "1.0.0",
		AwsRegion:  "us-east-1",
		Farm:       "test-farm",
	})

	ctx = request.ContextWithRequestIDs(ctx, request.RequestIDs{
		RequestID:     "id1",
		CorrelationID: "id2",
	})

	logger := log.NewFromCtx(ctx)
	logger.Debug().
		Str("test-str", "str").
		Int("test-number", 1).
		Msg("initialise handler")

	// Output:
	// {"level":"debug","app":"test-app","app_version":"1.0.0","aws_region":"us-east-1","aws_account_id":"","farm":"test-farm","request_id":"id1","correlation_id":"id2","test-str":"str","test-number":1,"event":"initialise_handler","message":"initialise handler"}
}
