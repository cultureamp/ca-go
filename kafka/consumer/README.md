# ca-go/kafka/consumer

The `kafka/consumer` package implements a blocking and non-blocking Kafka consumer with Avro decoding via our Schema Registry. Clients using this library need only supply a `Receiver` function which is called for every consumed message.

To avoid using CGO, this `ca-go` implementation does not use `https://github.com/confluentinc/confluent-kafka-go` and instead provides a wrapper over `https://github.com/IBM/sarama`. While this comes with some small risk of incompatibilities or feature lag, we think this trade-off is worth the cost to avoid having to enable CGO (and having to live with all the downsides that introduces). We are happy to revisit this if need be.

The implementation follows most of the guidance and recommendations found in this Sarama guide: [Effective Kafka Consumption in Golang: A Comprehensive Guide](https://medium.com/@moabbas.ch/effective-kafka-consumption-in-golang-a-comprehensive-guide-aac54b5b79f0)

## Environment Variables

- KAFKA_BROKERS = The list of Kafka brokers to use. Calling `WithBrokers(brokers)` overwrites this default.
- KAFKA_TOPICS = The list of Kafka topics to consume. Calling `WithTopics(topics)` overwrites this default.
- SCHEMA_REGISTRY_URL = The Avro Schema Registry to use. Calling `WithSchemaRegistryURL(url)` overwrites this default.
- KAFKA_SASL_USERNAME = The sasl scram username used to authenticate with the Kafka cluster (see [Kafka Security](https://cultureamp.atlassian.net/wiki/spaces/TDS/pages/2859368483/Kafka+Security)).
- KAFKA_SASL_PASSWORD = The sasl scram password used to authenticate with the Kafka cluster (see [Kafka Security](https://cultureamp.atlassian.net/wiki/spaces/TDS/pages/2859368483/Kafka+Security)).

## Avro and the Schema Registry

Before each Kafka message is sent to the client Receiver function, the msg.Value is arvo decoded via our [Schema Registry](https://cultureamp.atlassian.net/wiki/spaces/TDS/pages/2809562876/Team+Data+Services+On-Call+doc#Endpoints)

- <https://schema-registry.kafka.usw2.prod-us.cultureamp.io/>
- <https://schema-registry.kafka.usw2.dev-us.cultureamp.io/>  (via Dev VPN)

The decoded value is added to the `consumer.ReceivedMessage` struct as `DecodedText`. Typically this will be a `json` document that you can then `json.Unmarshal()` to your domain object and then process how you wish (eg. upsert into a database table).

## Subscriber

Below is an example from `example_test.go` which creates a typical consumer.Subscriber which blocks on `ConsumeAll` so clients either need to manage this in their own go-routine or some other fashion.

```go
func ExampleNewSubscriber() {
 // create a new subscriber
 subscriber, err := consumer.NewSubscriber(
  consumer.WithBrokers([]string{"localhost:9092"}),        // if missing, will default to env var 'KAFKA_BROKERS'
  consumer.WithTopics([]string{"test-topic"}),             // if missing, will default to env var 'KAFKA_TOPICS'
  consumer.WithSchemaRegistryURL("http://localhost:8081"), // if missing, will default to env var 'SCHEMA_REGISTRY_URL'
  consumer.WithGroupID("group_id"),
  consumer.WithReturnOnClientDispathError(true),
  consumer.WithSaslUsername("test_user"),
  consumer.WithSaslPassword("test_pwd"),
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

Below is an example from `example_test.go` which creates a typical consumer.Service which does not block the client routinue and will continue to run until `Stop` is called, or an error is encounted in the Reciever function and WithReturnOnClientDispathError is `true`.

```go
func ExampleNewService() {
 // create a new service
 service, err := consumer.NewService(
  consumer.WithBrokers([]string{"localhost:9092"}),        // if missing, will default to env var 'KAFKA_BROKERS'
  consumer.WithTopics([]string{"test-topic"}),             // if missing, will default to env var 'KAFKA_TOPICS'
  consumer.WithSchemaRegistryURL("http://localhost:8081"), // if missing, will default to env var 'SCHEMA_REGISTRY_URL'
  consumer.WithGroupID("group_id"),
  consumer.WithReturnOnClientDispathError(true),
  consumer.WithSaslUsername("test_user"),
  consumer.WithSaslPassword("test_pwd"),
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

 // stop the Service
 err = service.Stop()
 if err != nil {
  panic(err)
 }

 // Output:
 //
}
```
