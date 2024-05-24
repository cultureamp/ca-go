package consumer

import "github.com/IBM/sarama"

type client interface {
	NewConsumerGroup(brokers []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error)
	CommitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage)
}

type saramaClient struct{}

func newSaramaClient() *saramaClient {
	return &saramaClient{}
}

func (sc *saramaClient) NewConsumerGroup(brokers []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(brokers, groupId, config)
}

func (sc *saramaClient) CommitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	session.MarkMessage(msg, "")
}
