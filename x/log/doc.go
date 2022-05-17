package log

// Package log enables you to log events with different levels, customise fields,
// as well as embedding env configs from context.
//
// See the executable examples below and on individual functions for more
// details on usage.
//
// You configure and initialise logger using NewFromCtx():
//    logger := log.NewFromCtx(ctx)
//
// You could embed environment configs and request configs into context by using the request lib in ca-go
// then the configs would be logged as customised fields
//    import "github.com/cultureamp/ca-go/x/request"
//    ctx := context.Background()
//    ctx = log.ContextWithEnvConfig(ctx, log.EnvConfig{
//     AppName:    "my-app",
//     AppVersion: "1.0.0",
//     AwsRegion:  "",
//     Farm:       "local",
//    })
//    ctx = request.ContextWithRequestIDs(ctx, request.RequestIDs{
//     RequestID:     "id1",
//     CorrelationID: "id2",
//    })
//
// Log format would be plain text if Farm == "local", and using json formatt for any other cases.
// This is designed to make the log more human readble when you run your app locally.
//
// Once you created a logger entry, all methods provided by a logrus entry are available
//    logger.WithFields(map[string]any{
//      "key1": "val1",
//      "key2": "val2",
//    }).Debug("something happened")
//
// Errors can be reported using WithError():
//    logger.WithError(err).Error("something went wrong")
//
