//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -destination=mock_reader_test.go -package consumer . Reader
package consumer

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

const (
	defaultDialerTimeout              = 10 * time.Second
	defaultBackoffRandomizationFactor = 0
	defaultBackoffMaxInterval         = 5 * time.Hour
	defaultBackoffMultiplier          = 8
	defaultBackoffMaxElapsedTime      = 0
)

// GroupConfig is a configuration object used to create a new Group. The default
// consumer count in a group is 1 unless specified otherwise.
type GroupConfig struct {
	Count   int
	Brokers []string
	Topic   string
	GroupID string

	MinBytes      int           // Default: 1MB
	MaxBytes      int           // Default: 10MB
	MaxWait       time.Duration // Default: 250ms
	QueueCapacity int           // Default: 100
}

// Group groups consumers together to concurrently consume and handle messages
// from a Kafka topic. Many groups with the same group ID are safe to use, which
// is particularly useful for groups across separate instances.
//
// It is worth noting that publishing failed messages to a dead letter queue is
// not supported and instead would need to be included in your handler implementation.
type Group struct {
	ID      string
	config  GroupConfig
	opts    []Option
	stopChs []chan struct{}
}

// NewGroup returns a new Group configured with the provided dialer and config.
func NewGroup(config GroupConfig, opts ...Option) *Group {
	if config.Count <= 0 {
		config.Count = 1
	}

	return &Group{
		ID:     fmt.Sprintf("%s-%s", strings.ToLower(config.GroupID), uuid.New().String()[:7]), // semi-random slug
		config: config,
		opts:   opts,
	}
}

// Run concurrently consumes and handles messages from the topic across all
// consumers in the group. The method call returns an error channel that is used
// to receive any consumer errors. The run process is only stopped if the context
// is canceled, the consumer has been closed, or all consumers in the group have
// errored.
func (g *Group) Run(ctx context.Context, handler Handler) <-chan error {
	var wg sync.WaitGroup
	errCh := make(chan error, g.config.Count)

	for i := range g.config.Count {
		wg.Add(1)

		// Consumers must be created and run in sequential order so that Kafka can
		// successfully re-balance the group as each is added. This unfortunately
		// prevents us from receiving a passed in list of consumers, which is
		// arguably a cleaner approach.
		cfg := Config{
			ID:            fmt.Sprintf("%s-%d", g.ID, i),
			Brokers:       g.config.Brokers,
			Topic:         g.config.Topic,
			MinBytes:      g.config.MinBytes,
			MaxBytes:      g.config.MaxBytes,
			MaxWait:       g.config.MaxWait,
			QueueCapacity: g.config.QueueCapacity,
			groupID:       g.config.GroupID,
		}
		c := NewConsumer(cfg, g.opts...)

		go func() {
			defer wg.Done()
			if err := c.Run(ctx, handler); err != nil {
				errCh <- errors.Errorf("consumer failed: %w", err)
			}
		}()
		g.stopChs = append(g.stopChs, c.stopCh)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return errCh
}

// Stop stops the group. It waits for the current message (if any) in each
// consumer to finish being handled before closing the reader streams, preventing
// each consumer from reading any more messages.
func (g *Group) Stop() {
	for _, stopCh := range g.stopChs {
		close(stopCh)
	}
}

// DialerSCRAM512 returns a Kafka dialer configured with SASL authentication
// to securely transmit the provided credentials to Kafka using SCRAM-SHA-512.
func DialerSCRAM512(username string, password string) (*kafka.Dialer, error) {
	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		return nil, err
	}

	return &kafka.Dialer{
		Timeout:       defaultDialerTimeout,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS:           &tls.Config{MinVersion: tls.VersionTLS12},
	}, nil
}

type HandlerRetryBackOffConstructor func() backoff.BackOff

// NonStopExponentialBackOff is the suggested backoff retry strategy for consumers
// groups handling messages where ordering matters e.g. a data-capture stream from
// a database. This results in endless retries to prevent Kafka from re-balancing
// the group so that each consumer does not eventually experience the same error.
//
// Retry intervals: 500ms, 4s, 32s, 4m, 34m, 4.5h, 5h (max).
//
// The max interval of 5 hours is intended to leave enough time for manual
// intervention if necessary.
func NonStopExponentialBackOff() backoff.BackOff {
	bo := backoff.NewExponentialBackOff()
	bo.RandomizationFactor = defaultBackoffRandomizationFactor
	bo.MaxInterval = defaultBackoffMaxInterval
	bo.Multiplier = defaultBackoffMultiplier
	bo.MaxElapsedTime = defaultBackoffMaxElapsedTime
	return bo
}
