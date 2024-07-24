# ca-go/kafka/consumer

The `kafka/consumer` package implements a blocking and non-blocking Kafka consumer with Avro decoding via our Schema Registry. Clients using this library need only supply a `Receiver` function which is called for every consumed message.

To avoid using CGO, the `ca-go` implementation does not use `https://github.com/confluentinc/confluent-kafka-go` and instead provides a wrapper over the `https://github.com/IBM/sarama` implementation. While this comes with some small risk of impatibilities or feature lag, we think this trade-off is worth the cost to avoid having to enable CGO (and having to live with all the downsides that introduces).

## Environment Variables

- KAFKA_BROKERS = The list of Kafka brokers to use. Calling `WithBrokers(brokers)` overwrites this default.
- KAFKA_TOPICS = The list of Kafka topics to consume. Calling `WithTopics(topics)` overwrites this default.
- SCHEMA_REGISTRY_URL = The Avro Schema Registry to use. Calling `WithSchemaRegistryURL(url)` overwrites this default.

## Avro and the Schema Registry

- <https://schema-registry.kafka.usw2.prod-us.cultureamp.io/>
- <https://schema-registry.kafka.usw2.dev-us.cultureamp.io/>  (via Dev VPN)

Documentation is here:  <https://cultureamp.atlassian.net/wiki/spaces/TDS/pages/2809562876/Team+Data+Services+On-Call+doc>

## Subscriber

```go
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
    // if WithReturnOnClientDispathError(true)
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
```

## Service

```go
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
    // if WithReturnOnClientDispathError(true)
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
```
