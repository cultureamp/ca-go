//nolint:testableexamples
package sentry_test

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cultureamp/ca-go/sentry"
)

var (
	app         string
	appVersion  string
	buildNumber string
	branch      string
	commit      string
)

type Settings struct {
	SentryDSN string
	Farm      string
	AppEnv    string
}

func getSettings() *Settings {
	return &Settings{
		SentryDSN: os.Getenv("SENTRY_DSN"),
		Farm:      os.Getenv("FARM"),
		AppEnv:    os.Getenv("APP_ENV"),
	}
}

func Example_lambda() {
	// This is an example of how to use the errorreport package in a Lambda
	// function. The following is an example `main` function.

	ctx := context.Background()

	// in a real application, use something like "github.com/kelseyhightower/envconfig"
	settings := getSettings()

	// configure error reporting settings
	err := sentry.Init(
		sentry.WithDSN(settings.SentryDSN),
		sentry.WithRelease(app, appVersion),
		sentry.WithEnvironment(settings.AppEnv),
		sentry.WithBuildDetails(settings.Farm, buildNumber, branch, commit),
		sentry.WithServerlessTransport(),

		// optionally add a tag to every error report
		sentry.WithTag("animal", "gopher"),

		// or add multiple tags at once to be added to every error report
		sentry.WithTags(map[string]string{
			"genus":   "phoenicoparrus",
			"species": "jamesi",
		}),

		// optionally customise error title with the root cause message
		sentry.WithBeforeFilter(sentry.RootCauseAsTitle),
	)
	if err != nil {
		log.Panic("sentry_init",err).Send()
	}

	// wrap the lambda handler function with error reporting
	handler := sentry.LambdaMiddleware(Handler)

	// start the lambda function
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}

// Handler is the lambda handler function with the logic to be executed. In this
// case, it's a Kinesis event handler, but this could be a handler for any
// Lambda event.
func Handler(ctx context.Context, event events.KinesisEvent) error {
	for _, record := range event.Records {
		if err := processRecord(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func processRecord(ctx context.Context, record events.KinesisEventRecord) error {
	// Decorate will add these details to any error report that is sent to
	// Sentry in the context of this method. (Note the use of defer.)
	defer sentry.Decorate(map[string]string{
		"event_id":        record.EventID,
		"partition_key":   record.Kinesis.PartitionKey,
		"sequence_number": record.Kinesis.SequenceNumber,
	})()

	// do something with the record
	return nil
}

func Example_fargate() {
	// This is an example of how to use the errorreport package in a Main
	// function. The following is an example `main` function.

	ctx := context.Background()
	settings := getSettings()

	// configure error reporting settings
	err := sentry.Init(
		sentry.WithDSN(settings.SentryDSN),
		sentry.WithRelease(app, appVersion),
		sentry.WithEnvironment(settings.AppEnv),
		sentry.WithBuildDetails(settings.Farm, buildNumber, branch, commit),
		sentry.WithTag("application_name", app),
		// optionally add a tag to every error report
		sentry.WithTag("animal", "gopher"),

		// or add multiple tags at once to be added to every error report
		sentry.WithTags(map[string]string{
			"genus":   "phoenicoparrus",
			"species": "jamesi",
		}),
		sentry.WithBeforeFilter(sentry.RootCauseAsTitle),
	)
	if err != nil {
		// log error
	}

	// handle core business logic here
	handleBusinessLogic(ctx)

	// capture panic and report to sentry before the program exits
	defer func() {
		if err := recover(); err != nil {
			sentry.GracefullyShutdown(err, time.Second*5)
		}
	}()
}

func handleBusinessLogic(ctx context.Context) {
	if _, err := doSomething(); err != nil {
		// report a error to sentry
		sentry.ReportError(ctx, err)
	}
}

func doSomething() (string, error) {
	return "", nil
}
