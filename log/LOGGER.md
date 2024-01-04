# ca-go/log

The `log` package implements the [Logging Standard](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard). The design of this package is to provide a simple structured logging system that can be used in a variety of situations without requiring high cognitive load.

There are no new loggers to create or pass around, instead there is a singleton logger created in the package that you can call directly.

The `log` package wraps [zerolog](https://github.com/rs/zerolog) and therefore requires that you end all logging statements with a `Msg("_your_message_here")` to actually emit the log.

## Environment Variables

You MUST set these:
- APP = The application name (eg. "employee-tasks-service")
- AWS_REGION = The AWS region this code is running in (eg. "us-west-1")
- PRODUCT = The product suite the service belongs to (eg. "engagement")

You can OPTIONALLY set these:
- LOG_LEVEL = One of DEBUG, INFO, WARN, ERROR, defaults to "INFO"
- AWS_ACCOUNT_ID = The AWS account Id this code is running in, defaults to  "local"
- FARM = The name of the farm or where the code is running, defaults to "local" (eg. "production", "dolly") 
- APP_VERSION = The version of the application, defaults to "1.0.0"


## Examples
```
package cago

import (
	"context"

	"github.com/cultureamp/ca-go/log"
	"github.com/cultureamp/ca-go/jwt"
)


func basic_example() {
	ctx := context.Background()

	log.Debug(ctx, "something_just_happened")

    log.Info(ctx, "something_else_happened").
		Str("resource", "resource_id").
        
		Int("test-number", 2).
		Msg("detailed information go here")

    log.Error(ctx, "user_added", err).
		Str("resource", "resource_id").
		Int("test-number", 4).
		Dict("properties", log.Properties().
			Str("bar", "baz").
			Int("n", 1),
		).Msg("further details can be added here")
}

func glamplify_example(t *testing.T) {
	ctx := context.Background()

	type fields map[string]interface{}

	now := time.Now()
	f := &fields{
		"key1":  "string value",
		"key2":  1,
		"now":   now.Format(time.RFC3339),
		"later": time.Now(),
	}

	log.Info(ctx, "log_fields").
		Interface("properties", f).
		Msg("detailed information explain")
}

func http_request_example(t *testing.T) {
	ctx := context.Background()

	// create a dummy request and add it to the context
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Add(log.TraceHeader, "trace_123_id")
	req.Header.Add(log.RequestHeader, "request_456_id")
	req.Header.Add(log.CorrelationHeader, "correlation_789_id")

	ctx = log.ContextWithRequest(ctx, req)

    // later when you log with that context, the http request headers 
    // will automatically be added to the log (as per the standard)
	log.Info(ctx, "info_event").
		Str("resource", "resource_id").
		Int("test-number", 1).
		Msg("logging should contain request headers")
}

func jwtauth_payload_example(t *testing.T) {
	ctx := context.Background()

	// copy to the jwt auth payload and add it to the context
	auth := log.AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}
	ctx = log.ContextWithAuthPayload(ctx, auth)

    // later when you log with that context, the account, user and 
    // real_user ids will automatically be added to the log
    // (as per the standard)
	log.Info(ctx, "info_event").
		Str("resource", "resource_id").
		Int("test-number", 1).
		Msg("logging should contain auth payload")
}
```
