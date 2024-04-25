package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/mock"
)

type mockKafkaClient struct {
	mock.Mock
}

func newMockKafkaClient() *mockKafkaClient {
	return &mockKafkaClient{}
}

func (m *mockKafkaClient) newConsumerGroup(brokers []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	return nil, errors.Errorf("missing mock implementation")
}

func (m *mockKafkaClient) commitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	// todo
}

type mockConsumer struct {
	mock.Mock
}

func newMockConsumer() *mockConsumer {
	return &mockConsumer{}
}

func (m *mockConsumer) Consume(ctx context.Context) error {
	return nil
}
