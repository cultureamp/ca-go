package consumer

import "github.com/IBM/sarama"

type kafkaClient interface {
	NewConsumerGroup(brokers []string, groupID string, config *sarama.Config) (sarama.ConsumerGroup, error)
	CommitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage)
	Commit(session sarama.ConsumerGroupSession)
}

type saramaClient struct{}

func newSaramaClient() *saramaClient {
	return &saramaClient{}
}

func (sc *saramaClient) NewConsumerGroup(brokers []string, groupID string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(brokers, groupID, config)
}

func (sc *saramaClient) CommitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	session.MarkMessage(msg, "")
}

func (sc *saramaClient) Commit(session sarama.ConsumerGroupSession) {
	session.Commit()
}
