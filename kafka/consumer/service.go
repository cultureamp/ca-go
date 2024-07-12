package consumer

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

type Service struct {
	subscriber    *Subscriber
	logger        sarama.StdLogger
	runnningMutex sync.Mutex
	running       bool
}

func NewService(opts ...Option) (*Service, error) {
	c, err := NewSubscriber(opts...)
	if err != nil {
		return nil, err
	}

	s := &Service{
		subscriber: c,
		logger:     c.conf.stdLogger,
		running:    false,
	}
	return s, nil
}

func (s *Service) Start(ctx context.Context) {
	s.runnningMutex.Lock()
	defer s.runnningMutex.Unlock()

	// if running, do nothing
	if s.running {
		s.logger.Printf("service: already running!")
		return
	}

	// blocking call, so run in a go-routine
	s.logger.Printf("service: starting...")
	go s.run(ctx)
	s.running = true
}

func (s *Service) run(ctx context.Context) {
	// blocking call until context Done, client dispatch error, or Kafka rebalance
	err := s.subscriber.Subscribe(ctx)
	if err != nil {
		s.logger.Printf("service: error consuming topic: '%s'", err.Error())
	}
}

func (s *Service) Stop() error {
	s.runnningMutex.Lock()
	defer s.runnningMutex.Unlock()

	// if already stopped, do nothing
	if !s.running {
		s.logger.Printf("service: already stopped!")
		return nil
	}

	s.logger.Printf("service: stopping...")
	err := s.subscriber.Stop()
	if err != nil {
		s.logger.Printf("service: error stopping consumer: '%s'", err.Error())
	}

	s.running = false
	return err
}