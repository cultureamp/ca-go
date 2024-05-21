package consumer

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

type Service struct {
	consumer KafkaConsumer

	mu      sync.Mutex
	running bool
}

func NewService(consumer KafkaConsumer) *Service {
	return &Service{
		consumer: consumer,
		running:  false,
	}
}

func (s *Service) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if running, do nothing
	if s.running {
		sarama.Logger.Printf("service: already running!")
		return
	}

	// blocking call, so run in a go-routine
	sarama.Logger.Printf("service: starting...")
	go s.run(ctx)
	s.running = true
}

func (s *Service) run(ctx context.Context) {
	// blocking call until context Done or Kafka rebalance
	err := s.consumer.Consume(ctx)
	if err != nil {
		sarama.Logger.Printf("service: error consuming topic: '%s'", err.Error())
	}
}

func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if already stopped, do nothing
	if !s.running {
		sarama.Logger.Printf("service: already stopped!")
		return nil
	}

	sarama.Logger.Printf("service: stopping...")
	err := s.consumer.Stop()
	if err != nil {
		sarama.Logger.Printf("service: error stopping consumer: '%s'", err.Error())
	}

	s.running = false
	return err
}
