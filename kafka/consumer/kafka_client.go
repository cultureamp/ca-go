package consumer

import "github.com/IBM/sarama"

type kafkaClient interface {
	NewConsumerGroup(brokers []string, groupID string, config *sarama.Config) (sarama.ConsumerGroup, error)
	NewConsumer(brokers []string, config *sarama.Config) (sarama.Consumer, error)
	MarkMessageConsumed(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage, metadata string)
	Commit(session sarama.ConsumerGroupSession)
}

type saramaClient struct{}

func newSaramaClient() *saramaClient {
	return &saramaClient{}
}

func (sc *saramaClient) NewConsumerGroup(brokers []string, groupID string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(brokers, groupID, config)
}

func (sc *saramaClient) NewConsumer(brokers []string, config *sarama.Config) (sarama.Consumer, error) {
	return sarama.NewConsumer(brokers, config)
}

func (sc *saramaClient) MarkMessageConsumed(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage, metadata string) {
	session.MarkMessage(msg, metadata)
}

func (sc *saramaClient) Commit(session sarama.ConsumerGroupSession) {
	session.Commit()
}
