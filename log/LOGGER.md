# ca-go/log

The `log` package implements the [Logging Standard](https://cultureamp.atlassian.net/wiki/spaces/TV/pages/3114598406/Logging+Standard). The design of this package is to provide a simple structured logging system that can be used in a variety of situations without requiring high cognitive load.

There are no new loggers to create or pass around, instead there is a singleton logger created in the package that you can call directly.

The `log` package wraps [zerolog](https://github.com/rs/zerolog) and therefore requires that you end all logging statements with a `Details("_your_message_here")` to actually emit the log.

## Environment Variables

You MUST set these:
- APP = The application name (eg. "employee-tasks-service")
- AWS_REGION = The AWS region this code is running in (eg. "us-west-1")
- PRODUCT = The product suite the service belongs to (eg. "engagement")

You can OPTIONALLY set these:
- LOG_LEVEL = One of DEBUG, INFO, WARN, ERROR, defaults to "INFO"
- AWS_ACCOUNT_ID = The AWS account Id this code is running in, defaults to  "development"
- FARM = The name of the farm or where the code is running, defaults to "local" (eg. "production", "dolly") 
- APP_VERSION = The version of the application, defaults to "1.0.0"

## Use in Unit Tests

By default the logger will not emit any output when running inside a test. You can override this behaviour by setting the `QUIET_MODE` environment variable to "false".

When running localling you can also set the `CONSOLE_WRITER` to "true" to change from json to key-value colour coded output. Note: Never run with the `CONSOLE_WRITER` set to "true" in production.

## Managing Loggers Yourself

While we recommend using the package level methods for their ease of use, you may desire to create and manage loggers yourself, which you can do by calling:

```
config := NewLoggerConfig()
// optionally override default properties on the config
return NewLogger(config)
```

## Log Examples
```
package cagoexample

import (
	"github.com/cultureamp/ca-go/log"
)

func basic_example() {
	var ipv4 net.IP

	then := time.Now()
	u := uuid.New()
	duration := time.Since(then)

	props := SubDoc().
		Str("str", "value").
		Int("int", 1).
		Bool("bool", true).
		Duration("dur", duration).
		IPAddr("ipaddr", ipv4).
		UUID("uuid", u)

	Debug("debug_with_all_field_types").
		WithSystemTracing().
		Properties(props).
		Details("logging should contain all types")

	Debug("debug_with_all_field_types").
		WithSystemTracing().
		Properties(props).
		Detailsf("logging should contain all types: %s", "ok")

	Debug("debug_with_all_field_types").
		WithSystemTracing().
		Properties(props).
		Send()
}

func http_request_example() {
	// create a dummy request and add it to the context
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Add(log.TraceHeader, "trace_123_id")
	req.Header.Add(log.RequestHeader, "request_456_id")
	req.Header.Add(log.CorrelationHeader, "correlation_789_id")

	Debug("debug_with_request_and_system_tracing").
		WithSystemTracing().
		WithRequestTracing(req).
		Properties(log.SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain both")
}

func jwtauth_payload_example() {
	// create a jwt payload
	auth := &log.AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	log.Info("info_with_auth_and_system_tracing").
		WithSystemTracing().
		WithAuthenticatedUserTracing(auth).
		Properties(log.SubDoc().
			Str("resource", "resource_id").
			Int("test-number", 1),
		).Details("logging should contain both")
}
```

## Legacy Glamplify Loggers

Included are both package and NewLegacyLogger methods that support the glamplfy `log.Fields{}` interface. Feel free to use this when migrating an existing project off `glamplify` to `ca-go`, but these are NOT recommended for use for new projects.


## Legacy Glamplify Examples
```
package cagoexample

import (
	"github.com/cultureamp/ca-go/log"
)

func glamplify_example() {
	ctx := context.Background()

	now := time.Now()
	f := log.Fields{
		"key1":    "string value",
		"key2":    1,
		"now":     now.Format(time.RFC3339),
		"later":   time.Now(),
		"details": "detailed message",
	}
	log.GlamplifyDebug(ctx, "log_fields", f)
	log.GlamplifyInfo(ctx, "log_fields", f)
	log.GlamplifyWarn(ctx, "log_fields", f)
	log.GlamplifyError(ctx, "log_fields", errors.New("test error"), f)

	// log.GlamplifyFatal calls os.exit() so this is hard to demonstrate!

	defer recoverFromPanic()
	log.GlamplifyPanic(ctx, "panic_error", errors.New("test error"), f)
}

func http_request_example() {
	ctx := context.Background()

	f := log.Fields{
		"key1":    "string value",
		"key2":    1,
		"now":     now.Format(time.RFC3339),
		"later":   time.Now(),
		"details": "detailed message",
	}

	// create a dummy request and add it to the context
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Add(log.TraceHeader, "trace_123_id")
	req.Header.Add(log.RequestHeader, "request_456_id")
	req.Header.Add(log.CorrelationHeader, "correlation_789_id")

	// Note: Glamplify logger automatically adds WithSystemTracing() to log.Fields{}
	log.GlamplifyDebug(ctx, "log_fields", f.WithRequestTracing(req))
}

func jwtauth_payload_example() {
	ctx := context.Background()

	f := log.Fields{
		"key1":    "string value",
		"key2":    1,
		"now":     now.Format(time.RFC3339),
		"later":   time.Now(),
		"details": "detailed message",
	}

	// create a jwt payload
	auth := &log.AuthPayload{
		CustomerAccountID: "account_123_id",
		RealUserID:        "real_456_id",
		UserID:            "user_789_id",
	}

	// Note: Glamplify logger automatically adds WithSystemTracing() to log.Fields{}
	log.GlamplifyInfo(ctx, "log_fields", f.WithAuthenticatedUserTracing(req))
}

func recoverFromLogPanic() {
	if saved := recover(); saved != nil {
		// convert to an error if it's not one already
		err, ok := saved.(error)
		if !ok {
			err = errors.New(fmt.Sprint(saved))
		}

		log.GlamplifyError("recovered_from_panic", err).Send()
	}
}
```
