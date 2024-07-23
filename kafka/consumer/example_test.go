package consumer_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cultureamp/ca-go/kafka/consumer"
)

func ExampleNewSubscriber() {
	// create a new subscriber
	subscriber, err := consumer.NewSubscriber(
		consumer.WithBrokers([]string{"localhost:9092"}),        // if missing, will default to env var 'KAFKA_BROKERS'
		consumer.WithTopics([]string{"test-topic"}),             // if missing, will default to env var 'KAFKA_TOPICS'
		consumer.WithSchemaRegistryURL("http://localhost:8081"), // if missing, will default to env var 'SCHEMA_REGISTRY_URL'
		consumer.WithGroupID("group_id"),
		consumer.WithHandler(func(ctx context.Context, msg *consumer.ReceivedMessage) error {
			// check topic, timestamp, etc. if need be

			// do something with the message, typically unmarshal the json to your domain object
			var myDomainObject interface{}
			err := json.Unmarshal([]byte(msg.DecodedText), &myDomainObject)
			if err != nil {
				// recover, or return this error, which will stop the subscriber from consuming any more messages
				return err
			}

			// save the domain object to a database, etc.
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	// Set an optional ctx deadline
	ctx := context.Background()
	deadline := time.Now().Add(1 * time.Second)
	ctx, cancelCtx := context.WithDeadline(ctx, deadline)
	defer cancelCtx()

	// consume all the messages.
	// Note: this blocks until ctx is done, cancelled, deadline reached or an error occurs.
	err = subscriber.ConsumeAll(ctx)

	// Output:
	//
}

func ExampleNewService() {
	// create a new service
	service, err := consumer.NewService(
		consumer.WithBrokers([]string{"localhost:9092"}),        // if missing, will default to env var 'KAFKA_BROKERS'
		consumer.WithTopics([]string{"test-topic"}),             // if missing, will default to env var 'KAFKA_TOPICS'
		consumer.WithSchemaRegistryURL("http://localhost:8081"), // if missing, will default to env var 'SCHEMA_REGISTRY_URL'
		consumer.WithGroupID("group_id"),
		consumer.WithHandler(func(ctx context.Context, msg *consumer.ReceivedMessage) error {
			// check topic, timestamp, etc. if need be

			// do something with the message, typically unmarshal the json to your domain object
			var myDomainObject interface{}
			err := json.Unmarshal([]byte(msg.DecodedText), &myDomainObject)
			if err != nil {
				// recover, or return this error, which will stop the subscriber from consuming any more messages
				return err
			}

			// save the domain object to a database, etc.
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	// The service runs in its own go-routine, so is non blocking
	ctx := context.Background()
	service.Start(ctx)

	// do other work here
	time.Sleep(1 * time.Second)

	// stop the Service
	err = service.Stop()
	if err != nil {
		panic(err)
	}

	// Output:
	//
}
