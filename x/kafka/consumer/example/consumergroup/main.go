package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	errs := consumerGroup.Run(context.Background(), handle)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err, ok := <-errs:
			if !ok { // err channel closed
				return
			}
			log.Println(err)
		case <-sigterm:
			signal.Stop(sigterm)
			return
		}
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
