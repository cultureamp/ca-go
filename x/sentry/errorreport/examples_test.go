package errorreport_test

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cultureamp/ca-go/x/log"
	"github.com/cultureamp/ca-go/x/sentry/errorreport"
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
	err := errorreport.Init(
		errorreport.WithDSN(settings.SentryDSN),
		errorreport.WithRelease(app, appVersion),
		errorreport.WithEnvironment(settings.AppEnv),
		errorreport.WithBuildDetails(settings.Farm, buildNumber, branch, commit),
		errorreport.WithServerlessTransport(),

		// optionally add a tag to every error report
		errorreport.WithTag("animal", "gopher"),

		// or add multiple tags at once to be added to every error report
		errorreport.WithTags(map[string]string{
			"genus":   "phoenicoparrus",
			"species": "jamesi",
		}),

		// optionally customise error title with the root cause message
		errorreport.WithBeforeFilter(errorreport.RootCauseAsTitle),
	)
	if err != nil {
		// FIX: write error to log
		os.Exit(1)
	}

	// wrap the lambda handler function with error reporting
	handler := errorreport.LambdaMiddleware(Handler)

	// start the lambda function
	lambda.StartWithContext(ctx, handler)
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
	defer errorreport.Decorate(map[string]string{
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
	logger := log.NewFromCtx(ctx)

	settings := getSettings()

	// configure error reporting settings
	err := errorreport.Init(
		errorreport.WithDSN(settings.SentryDSN),
		errorreport.WithRelease(app, appVersion),
		errorreport.WithEnvironment(settings.AppEnv),
		errorreport.WithBuildDetails(settings.Farm, buildNumber, branch, commit),
		errorreport.WithTag("application_name", app),
		// optionally add a tag to every error report
		errorreport.WithTag("animal", "gopher"),

		// or add multiple tags at once to be added to every error report
		errorreport.WithTags(map[string]string{
			"genus":   "phoenicoparrus",
			"species": "jamesi",
		}),
		errorreport.WithBeforeFilter(errorreport.RootCauseAsTitle),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("sentry init")
	}

	// handle core business logic here
	handleBusinessLogic(ctx)

	// capture panic and report to sentry before the program exits
	defer func() {
		if err := recover(); err != nil {
			errorreport.GracefullyShutdown(err, time.Second*5)
		}
	}()
}

func handleBusinessLogic(ctx context.Context) {
	if _, err := doSomething(); err != nil {
		// report a error to sentry
		errorreport.ReportError(ctx, err)
	}
}

func doSomething() (string, error) {
	return "", nil
}
