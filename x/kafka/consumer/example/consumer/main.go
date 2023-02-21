package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"

	"github.com/cultureamp/ca-go/x/kafka/consumer"
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
	c := consumer.NewConsumer(kafka.DefaultDialer, cfg)

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
