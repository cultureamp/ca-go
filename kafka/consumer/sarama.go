package consumer

import "github.com/IBM/sarama"

type kafkaClient interface {
	newConsumerGroup(brokers []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error)
	commitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage)
}

type saramaClient struct{}

func newSaramaClient() *saramaClient {
	return &saramaClient{}
}

func (sc *saramaClient) newConsumerGroup(brokers []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(brokers, groupId, config)
}

func (sc *saramaClient) commitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	session.MarkMessage(msg, "")
}
