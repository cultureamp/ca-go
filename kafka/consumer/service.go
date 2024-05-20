package consumer

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

type Service struct {
	consumer KafkaConsumer
	cleanup  Cleanup
	mu       sync.Mutex
}

func NewService(consumer KafkaConsumer) *Service {
	return &Service{
		consumer: consumer,
	}
}

func (s *Service) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if running, do nothing

	// blocking call, so run in a go-routine
	go s.run(ctx)
}

func (s *Service) run(ctx context.Context) {
	// blocking call until context Done or Kafka rebalance
	cleanup, err := s.consumer.Consume(ctx)
	if err != nil {
		sarama.Logger.Printf("service: error consuming topic: '%s'", err.Error())
	}
	s.cleanup = cleanup
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if already stopped, do nothing
	if s.cleanup == nil {
		return nil
	}

	err := s.cleanup()
	if err != nil {
		sarama.Logger.Printf("service: error stopping consumer: '%s'", err.Error())
	}

	return err
}
