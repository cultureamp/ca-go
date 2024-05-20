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

func (m *mockKafkaClient) NewConsumerGroup(brokers []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	args := m.Called(brokers, groupId, config)
	gc := args.Get(0).(sarama.ConsumerGroup)
	return gc, args.Error(1)
}

func (m *mockKafkaClient) CommitMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	m.Called(session, msg)
}

// mockSaramaConsumerGroup must support the sarama.ConsumerGroup interface.
type mockSaramaConsumerGroup struct {
	sesson sarama.ConsumerGroupSession
	claim  sarama.ConsumerGroupClaim
	mock.Mock
}

func newMockConsumerGroup(sesson sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) *mockSaramaConsumerGroup {
	return &mockSaramaConsumerGroup{
		sesson: sesson,
		claim:  claim,
	}
}

func (m *mockSaramaConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := m.Called(ctx, topics, handler)

	err := args.Error(0)
	// success case
	if err == nil {
		// we need to call the handler with a session & claim
		err = handler.ConsumeClaim(m.sesson, m.claim)
	}

	return err
}

func (m *mockSaramaConsumerGroup) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockSaramaConsumerGroup) Errors() <-chan error {
	args := m.Called()
	ce := args.Get(0).(<-chan error)
	return ce
}

func (m *mockSaramaConsumerGroup) Pause(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *mockSaramaConsumerGroup) Resume(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *mockSaramaConsumerGroup) PauseAll() {
	m.Called()
}

func (m *mockSaramaConsumerGroup) ResumeAll() {
	m.Called()
}

// mockSaramaConsumerGroupHandler must support the sarama.ConsumerGroupHandler interface.
type mockSaramaConsumerGroupHandler struct {
	mock.Mock
}

func newMockReceiver() *mockSaramaConsumerGroupHandler {
	return &mockSaramaConsumerGroupHandler{}
}

func (m *mockSaramaConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *mockSaramaConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *mockSaramaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	args := m.Called(session, claim)
	return args.Error(0)
}

// mockSaramaConsumerGroupSession must support the sarama.ConsumerGroupSession interface.
type mockSaramaConsumerGroupSession struct {
	mock.Mock
}

func newMockConsumerGroupSession() *mockSaramaConsumerGroupSession {
	return &mockSaramaConsumerGroupSession{}
}

func (m *mockSaramaConsumerGroupSession) Claims() map[string][]int32 {
	args := m.Called()
	claims := args.Get(0).(map[string][]int32)
	return claims
}

func (m *mockSaramaConsumerGroupSession) MemberID() string {
	args := m.Called()
	memberId := args.Get(0).(string)
	return memberId
}

func (m *mockSaramaConsumerGroupSession) GenerationID() int32 {
	args := m.Called()
	genId := args.Get(0).(int32)
	return genId
}

func (m *mockSaramaConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *mockSaramaConsumerGroupSession) Commit() {
	m.Called()
}

func (m *mockSaramaConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *mockSaramaConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}

func (m *mockSaramaConsumerGroupSession) Context() context.Context {
	args := m.Called()
	ctx := args.Get(0).(context.Context)
	return ctx
}

// mockSaramaConsumerGroupClaim must support the sarama.ConsumerGroupClaim interface.
type mockSaramaConsumerGroupClaim struct {
	mock.Mock
}

func newMockConsumerGroupClaim() *mockSaramaConsumerGroupClaim {
	return &mockSaramaConsumerGroupClaim{}
}

func (m *mockSaramaConsumerGroupClaim) Topic() string {
	args := m.Called()
	topic := args.Get(0).(string)
	return topic
}

func (m *mockSaramaConsumerGroupClaim) Partition() int32 {
	args := m.Called()
	partition := args.Get(0).(int32)
	return partition
}

func (m *mockSaramaConsumerGroupClaim) InitialOffset() int64 {
	args := m.Called()
	iOffiset := args.Get(0).(int64)
	return iOffiset
}

func (m *mockSaramaConsumerGroupClaim) HighWaterMarkOffset() int64 {
	args := m.Called()
	hmOffset := args.Get(0).(int64)
	return hmOffset
}

func (m *mockSaramaConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	args := m.Called()
	messages := args.Get(0).(<-chan *sarama.ConsumerMessage)
	return messages
}
