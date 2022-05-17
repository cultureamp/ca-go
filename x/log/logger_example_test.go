package log_test

import (
	"context"
	"fmt"

	"github.com/cultureamp/ca-go/x/log"
	"github.com/cultureamp/ca-go/x/request"
)

func Example() {
	// This is an example of how to use the log package in a Lambda handler function.
	// The following is an example `main` function.

	ctx := context.Background()
	ctx = log.ContextWithEnvConfig(ctx, log.EnvConfig{
		AppName:    "my-app",
		AppVersion: "1.0.0",
		AwsRegion:  "",
		Farm:       "test",
		//Farm:       "local",
	})
	ctx = request.ContextWithRequestIDs(ctx, request.RequestIDs{
		RequestID:     "id1",
		CorrelationID: "id2",
	})
	logger := log.NewFromCtx(ctx)

	_ = createHandler(ctx)
	logger.Debug("initialise handler")
	// Output
	//{
	//  "AppName": "my-app",
	//  "AppVersion": "1.0.0",
	//  "AwsRegion": "",
	//  "CorrelationID": "id2",
	//  "Farm": "test",
	//  "RequestID": "id1",
	//  "file": "path/to/file.go:38",
	//  "func": "package.function",
	//  "level": "debug",
	//  "msg": "initialise handler",
	//  "time": "2022-05-17T11:44:08+10:00"
	//}
}

func createHandler(ctx context.Context) func() {
	return func() {
		logger := log.NewFromCtx(ctx)

		err := process()
		if err != nil {
			logger.WithError(err).Error("process failed")
		}

		logger.WithFields(map[string]any{
			"key1": "val1",
			"key2": "val2",
		}).Debug("process finished")
	}
}

func process() error {
	// processing something
	return fmt.Errorf("something went wrong")
}
