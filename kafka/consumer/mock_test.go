package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/mock"
)

// mockKafkaClient must support the kafkaClient interface.
type mockKafkaClient struct {
	mock.Mock
}

func newMockKafkaClient() *mockKafkaClient {
	return &mockKafkaClient{}
}

func (m *mockKafkaClient) newConsumerGroup(brokers []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	args := m.Called(brokers, groupId, config)
	gc := args.Get(0).(sarama.ConsumerGroup)
	return gc, args.Error(1)
}

func (m *mockKafkaClient) commitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	m.Called(session, msg)
}

// mockConsumerGroup must support the sarama.ConsumerGroup interface.
type mockConsumerGroup struct {
	sesson sarama.ConsumerGroupSession
	claim  sarama.ConsumerGroupClaim
	mock.Mock
}

func newMockConsumerGroup(sesson sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) *mockConsumerGroup {
	return &mockConsumerGroup{
		sesson: sesson,
		claim:  claim,
	}
}

func (m *mockConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := m.Called(ctx, topics, handler)

	err := args.Error(0)
	// success case
	if err == nil {
		// we need to call the handler with a session & claim
		err = handler.ConsumeClaim(m.sesson, m.claim)
	}

	return err
}

func (m *mockConsumerGroup) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockConsumerGroup) Errors() <-chan error {
	args := m.Called()
	ce := args.Get(0).(<-chan error)
	return ce
}

func (m *mockConsumerGroup) Pause(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *mockConsumerGroup) Resume(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *mockConsumerGroup) PauseAll() {
	m.Called()
}

func (m *mockConsumerGroup) ResumeAll() {
	m.Called()
}

// mockReceiver must support the sarama.ConsumerGroupHandler
type mockReceiver struct {
	mock.Mock
}

func newMockReceiver() *mockReceiver {
	return &mockReceiver{}
}

func (m *mockReceiver) Setup(session sarama.ConsumerGroupSession) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *mockReceiver) Cleanup(session sarama.ConsumerGroupSession) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *mockReceiver) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	args := m.Called(session, claim)
	return args.Error(0)
}

type mockConsumerGroupSession struct {
	mock.Mock
}

func newMockConsumerGroupSession() *mockConsumerGroupSession {
	return &mockConsumerGroupSession{}
}

func (m *mockConsumerGroupSession) Claims() map[string][]int32 {
	args := m.Called()
	claims := args.Get(0).(map[string][]int32)
	return claims
}

func (m *mockConsumerGroupSession) MemberID() string {
	args := m.Called()
	memberId := args.Get(0).(string)
	return memberId
}

func (m *mockConsumerGroupSession) GenerationID() int32 {
	args := m.Called()
	genId := args.Get(0).(int32)
	return genId
}

func (m *mockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *mockConsumerGroupSession) Commit() {
	m.Called()
}

func (m *mockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *mockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}

func (m *mockConsumerGroupSession) Context() context.Context {
	args := m.Called()
	ctx := args.Get(0).(context.Context)
	return ctx
}

type mockConsumerGroupClaim struct {
	mock.Mock
}

func newMockConsumerGroupClaim() *mockConsumerGroupClaim {
	return &mockConsumerGroupClaim{}
}

func (m *mockConsumerGroupClaim) Topic() string {
	args := m.Called()
	topic := args.Get(0).(string)
	return topic
}

func (m *mockConsumerGroupClaim) Partition() int32 {
	args := m.Called()
	partition := args.Get(0).(int32)
	return partition
}

func (m *mockConsumerGroupClaim) InitialOffset() int64 {
	args := m.Called()
	iOffiset := args.Get(0).(int64)
	return iOffiset
}

func (m *mockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	args := m.Called()
	hmOffset := args.Get(0).(int64)
	return hmOffset
}

func (m *mockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	args := m.Called()
	messages := args.Get(0).(<-chan *sarama.ConsumerMessage)
	return messages
}
