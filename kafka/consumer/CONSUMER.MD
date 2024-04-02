# ca-go/kafka/consumer

The `kafka/consumer` package provides access to consume Kafka messages from a topic. The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

The package creates a default consumer that uses the `KAFKA_BROKERS` environment variable. For ease of use, it is recommended that you use the package level method Consume. However, if you need to support other configurations then you can create a `NewConsumer()` and manage the class life-cycle yourself.


## Environment Variables

To use the package level methods `Consume` method you MUST set these:

- KAFKA_BROKERS = A comma seperate list of strings containing the kafka brokers.

## Managing Consumers Yourself

While we recommend using the package level methods for their ease of use, you may desire to create and manage consumers yourself, which you can do by calling:

- NewConsumer(dialer *kafka.Dialer, config Config, opts ...Option) *Consumer


## Examples

### Package Level Consumer

```
ch, stop := Consume(ctx, "topic-name")
// Note: call "stop()" when finished to close channel
// when finished close it
defer func() {
    err := stop()
    if err != nil {
        fmt.Printf("Error: %v", err)
    }
}()

for running {
	select {
	case msg, ok := <-ch:
		// if ok then do something with the msg

	case <-time.After(time.Duration(1) * time.Second):
		// timed out before receiving the msg then do something else

    case <-ctx.Done():
        // context deadline reached, time to stop
        running = false
	}
}
```

## Testing and Mocks
