package errorreport_test

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cultureamp/ca-go/x/sentry/errorreport"
)

// Set these variables at build time using linker flags (-ldflags)
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

func Example() {
	// This is an example of how to use the errorreport package in a Lambda
	// function. The following is an example `main` function.

	ctx := context.Background()

	// in a real application, use something like "github.com/kelseyhightower/envconfig"
	settings := &Settings{
		SentryDSN: os.Getenv("SENTRY_DSN"),
		Farm:      os.Getenv("FARM"),
		AppEnv:    os.Getenv("APP_ENV"),
	}

	// configure error reporting settings
	err := errorreport.Init(
		errorreport.WithDSN(settings.SentryDSN),
		errorreport.WithRelease(app, appVersion),
		errorreport.WithEnvironment(settings.AppEnv),
		errorreport.WithBuildDetails(settings.Farm, buildNumber, branch, commit),
		errorreport.WithServerlessTransport(),
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
