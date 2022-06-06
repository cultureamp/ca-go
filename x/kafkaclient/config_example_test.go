package kafkaclient

import (
	"github.com/Shopify/sarama"
	"log"
	"strings"
)

func Example() {
	clientCfg := Config{
		Brokers:  "brokers",
		Username: "username",
		Password: "password",
		Topic:    "topic",
		ClientID: "client_id",
	}

	splitBrokers := strings.Split(clientCfg.Brokers, ",")

	connConfig := GetConnConfig(clientCfg)
	producer, err := sarama.NewSyncProducer(splitBrokers, connConfig)
	if err != nil {
		log.Fatal("failed to create producer")
	}
	defer producer.Close()
}
