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
	brokers       string
	topic         string
	groupID       string
	consumerCount int
)

func main() {
	parseFlags()

	groupCfg := consumer.GroupConfig{
		Count:   consumerCount,
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
		GroupID: groupID,
	}
	consumerGroup := consumer.NewGroup(kafka.DefaultDialer, groupCfg)
	defer consumerGroup.Close()

	log.Printf("consumer group %s started for topic %s\n", consumerGroup.ID, topic)
	errCh := consumerGroup.Run(context.Background(), handle)
	for err := range errCh {
		if err != nil {
			log.Println(err)
		}
	}
}

func handle(ctx context.Context, msg kafka.Message, metadata consumer.Metadata) error {
	log.Printf("message at consumer: %s topic:%v partition:%v offset:%v	%s = %s\n",
		metadata.ConsumerID, msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value),
	)
	return nil
}

func parseFlags() {
	flag.StringVar(&brokers, "brokers", "", "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&groupID, "group-id", "", "Kafka consumer group identifier")
	flag.StringVar(&topic, "topic", "", "Kafka topic to be consumed")
	flag.IntVar(&consumerCount, "count", 1, "Kafka consumer group number of consumers")
	flag.Parse()

	if brokers == "" {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}
	if topic == "" {
		panic("no topics given to be consumed, please set the -topic flag")
	}
	if groupID == "" {
		panic("no Kafka consumer group defined, please set the -group flag")
	}
	if consumerCount == 0 {
		panic("no Kafka consumer group count defined, please set the -count flag")
	}
}
