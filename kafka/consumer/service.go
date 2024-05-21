package consumer

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

type Service struct {
	consumer KafkaConsumer
	logger   sarama.StdLogger

	runnningMutex sync.Mutex
	running       bool
}

func NewService(opts ...Option) (*Service, error) {
	c, err := NewConsumer(opts...)
	if err != nil {
		return nil, err
	}

	s := &Service{
		consumer: c,
		logger:   c.conf.stdLogger,
		running:  false,
	}
	return s, nil
}

func (s *Service) Start(ctx context.Context) {
	s.runnningMutex.Lock()
	defer s.runnningMutex.Unlock()

	// if running, do nothing
	if s.running {
		s.logf("service: already running!")
		return
	}

	// blocking call, so run in a go-routine
	s.logf("service: starting...")
	go s.run(ctx)
	s.running = true
}

func (s *Service) run(ctx context.Context) {
	// blocking call until context Done, client dispatch error, or Kafka rebalance
	err := s.consumer.Consume(ctx)
	if err != nil {
		s.logf("service: error consuming topic: '%s'", err.Error())
	}
}

func (s *Service) Stop() error {
	s.runnningMutex.Lock()
	defer s.runnningMutex.Unlock()

	// if already stopped, do nothing
	if !s.running {
		s.logf("service: already stopped!")
		return nil
	}

	s.logf("service: stopping...")
	err := s.consumer.Stop()
	if err != nil {
		s.logf("service: error stopping consumer: '%s'", err.Error())
	}

	s.running = false
	return err
}

func (s *Service) logf(format string, v ...interface{}) {
	if s.logger == nil {
		return
	}

	s.logger.Printf(format, v...)
}
