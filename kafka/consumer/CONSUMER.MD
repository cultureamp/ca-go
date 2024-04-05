# ca-go/kafka/consumer

The `kafka/consumer` package provides access to consume Kafka messages from a topic. The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

For now, the recommendation is to create and manage consumers yourself, which you can do by calling:
- NewConsumer(config Config, opts ...Option) *Consumer

# Consumer

A single consumer is the most basic use case when consuming messages from a
topic.

```
c := consumer.NewConsumer(dialer, config)
err := c.Run(ctx, handler)
```

# Consumer Group

Multiple consumers can be created as Group to consume messages at a higher rate.

```
g := consumer.Group(dialer, config)
errCh := g.Run(ctx, handler)

	for err := range errCh {
	    if err != nil {
		       panic(err)
	    }
	}
```

The consumer count and for the group is specified in the config parameter, which
determines the number of goroutines to spawn. Each goroutine has their own
Consumer and is set up so that individual consumer errors are reported back
via the error channel returned from Group.Run.

It is important to note that Kafka is ultimately responsible for managing group
members. As a result, a consumer group can easily be spread across multiple
instances by simply using the same group ID for each Group. Kafka will then
take care of re-balancing the group if members are added/removed.

# Examples

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/cultureamp/ca-go/kafka/consumer"
)

var (
	brokers string
	topic   string
)

func main() {
	parseFlags()

	cfg := consumer.Config{
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
	}
	c := consumer.NewConsumer(cfg)

	log.Printf("consumer started for topic %s\n", topic)
	if err := c.Run(context.Background(), handle); err != nil {
		panic(err)
	}
}

func handle(ctx context.Context, msg consumer.Message) error {
	log.Printf("message at consumer: %s topic:%v partition:%v offset:%v	%s = %s\n",
		msg.Metadata.ConsumerID, msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value),
	)
	return nil
}

func parseFlags() {
	flag.StringVar(&brokers, "brokers", "", "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&topic, "topic", "", "Kafka topic to be consumed")
	flag.Parse()

	if brokers == "" {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}
	if topic == "" {
		panic("no topics given to be consumed, please set the -topic flag")
	}
}
```
