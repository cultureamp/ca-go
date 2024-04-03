# ca-go/kafka/consumer

The `kafka/consumer` package provides access to consume Kafka messages from a topic. The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

The package creates a default consumer that uses the `KAFKA_BROKERS` environment variable. For ease of use, it is recommended that you use the package level method Consume. However, if you need to support other configurations then you can create a `NewConsumer()` and manage the class life-cycle yourself.


## Environment Variables

To use the package level methods `Consume` method you MUST set these:

- KAFKA_BROKERS = A comma seperate list of strings containing the kafka brokers.

## Managing Consumers Yourself

While we recommend using the package level methods for their ease of use, you may desire to create and manage consumers yourself, which you can do by calling:

- NewConsumer(config Config, opts ...Option) *Consumer


## Examples

### Package Level Consumer

```
func PackageConsumeTopic() {
	ctx := context.Background()
	ch, stop := consumer.Consume(ctx, "topic-name")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	// Read the next message from the topic
	msg, ok := <-ch
	fmt.Printf("Channel open=%t, topic=%s message_offset=%v\n", ok, msg.Topic, msg.Offset)
}

```


## Testing and Mocks


```
func MockedPackageExample() {
    defConsumers := consumer.DefaultConsumers
	consumer.DefaultConsumers = new(packageConsumerExampleMock)
	defer func() {
		consumer.DefaultConsumers = defConsumers
	}()

	ctx := context.Background()
	ch, stop := consumer.Consume(ctx, "topic-name")
	// when finished close it
	defer func() {
		err := stop()
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	ok := true
	for ok {
		select {
		// Read the messages from the topic
		case msg, opened := <-ch:
			if opened {
				fmt.Printf("Channel open=%t, topic=%s message_offset=%v\n", opened, msg.Topic, msg.Offset)
			} else {
				// time to stop, channel is closed
				ok = false
			}
		}
	}
}

type packageConsumerExampleMock struct{}

func (c *packageConsumerExampleMock) Consume(ctx context.Context, topic string) (<-chan consumer.Message, consumer.StopFunc) {
	var channel chan consumer.Message

	channel = make(chan consumer.Message, 10)
	go c.run(channel)
	return channel, func() error { return nil }
}

func (c *packageConsumerExampleMock) run(channel chan consumer.Message) {
	defer close(channel)

	// generate some messages and then exit when done
	for i := 0; i < 10; i++ {
		m := c.newMessage()
		channel <- m
	}
}

func (c *packageConsumerExampleMock) newMessage() consumer.Message {
	msg := consumer.Message{}

	msg.Topic = "topic-name"
	msg.Offset = rand.Int63()           //#nosec G404 -- This for the test runner reader, not used in production.
	msg.Partition = rand.Intn(20-1) + 1 //#nosec G404 -- This for the test runner reader, not used in production.
	msg.Value = []byte(uuid.New().String())
	return msg
}

```
