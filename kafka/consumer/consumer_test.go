package consumer

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConsumer(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantGenID bool
	}{
		{
			name: "new group with valid number of consumers",
			config: Config{
				ID:      "some-id",
				Brokers: []string{"some-address"},
				Topic:   "some-topic",
				groupID: "some-group-id",
			},
		},
		{
			name: "new group with invalid number of consumers defaults to 1",
			config: Config{
				Brokers: []string{"some-address"},
				Topic:   "some-topic",
				groupID: "some-group-id",
			},
			wantGenID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBackOff := NonStopExponentialBackOff
			wantBalancers := []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}}

			dialer, err := DialerSCRAM512("username", "password")
			require.NoError(t, err)

			c := NewConsumer(tt.config,
				WithExplicitCommit(),
				WithGroupBalancers(wantBalancers...),
				WithHandlerBackOffRetry(wantBackOff),
				WithDataDogTracing(),
				WithKafkaReader(func() Reader {
					return &MockReader{}
				}),
				WithKafkaDialer(dialer),
			)
			require.NotNil(t, c)
			assert.Equal(t, c.conf.Dialer, dialer)
			assert.NotNil(t, c.reader)
			assert.Implements(t, (*Reader)(nil), c.reader)
			assert.NotNil(t, c.clientHandler.BackOffConstructor)
			assert.NotNil(t, c.clientHandler.BackOffConstructor)
			assert.Equal(t, wantBalancers, c.conf.GroupBalancers)

			if tt.wantGenID {
				assert.NotEmpty(t, c.id)
			} else {
				assert.Equal(t, tt.config.ID, c.id)
			}
		})
	}
}

func TestConsumer_Run(t *testing.T) {
	var currMsg kafka.Message
	ctx := context.Background()
	wantTimes := 50

	reader := NewMockReader(gomock.NewController(t))
	reader.EXPECT().Close().Return(nil).Times(1)
	reader.EXPECT().ReadMessage(ctx).DoAndReturn(func(_ context.Context) (kafka.Message, error) {
		currMsg = randMsg()
		return currMsg, nil
	}).Times(wantTimes)

	consumer := NewConsumer(Config{},
		WithKafkaReader(func() Reader { return reader }),
		WithKafkaDialer(&kafka.Dialer{}),
	)

	i := 0
	handler := func(ctx context.Context, msg Message) error {
		require.Equal(t, currMsg, msg.Message)
		i++
		if i == wantTimes {
			require.NoError(t, consumer.Stop())
		}
		return nil
	}

	require.NoError(t, consumer.Run(ctx, handler))
}

func TestConsumer_Run_error(t *testing.T) {
	var wantErr error
	var gotAttempts int
	var didNotify bool

	wantHandlerErr := errors.New("some downstream error")

	tests := []struct {
		name             string
		wantError        error
		shouldNotify     bool
		numRetries       int
		contextCancelled bool
		setup            func(t *testing.T, consumer *Consumer)
	}{
		{
			name:         "consumer unable to handle message",
			wantError:    fmt.Errorf("consumer error: unable to handle message: %w", wantHandlerErr),
			shouldNotify: false,
			numRetries:   0,
		},
		{
			name:         "handler error after backoff retry and notify",
			wantError:    fmt.Errorf("consumer error: unable to handle message: %w", wantHandlerErr),
			shouldNotify: true,
			numRetries:   3,
			setup: func(t *testing.T, consumer *Consumer) {
				consumer.clientHandler.BackOffConstructor = func() backoff.BackOff {
					return &testBackoff{
						maxAttempts: 3,
					}
				}
				consumer.clientHandler.clientNotify = func(ctx context.Context, err error, msg Message) {
					assert.Equal(t, gotAttempts, msg.Metadata.Attempt)
					assert.Equal(t, wantHandlerErr, err)
					didNotify = true
				}
			},
		},
		{
			name:         "handler success after error backoff retry and notify",
			numRetries:   3,
			wantError:    nil,
			shouldNotify: true,
			setup: func(t *testing.T, consumer *Consumer) {
				consumer.clientHandler.BackOffConstructor = func() backoff.BackOff {
					return &testBackoff{}
				}
				consumer.clientHandler.clientNotify = func(ctx context.Context, err error, msg Message) {
					assert.Equal(t, gotAttempts, msg.Metadata.Attempt)
					if err != nil {
						assert.Equal(t, wantHandlerErr, err)
					}
					didNotify = true
				}
			},
		},
		{
			name:             "consumer context done error",
			wantError:        errors.New("consumer error: unable to handle message: context canceled"),
			contextCancelled: true,
			shouldNotify:     false,
			numRetries:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.contextCancelled {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			didNotify = false
			gotAttempts = 0
			wantErr = tt.wantError

			reader := NewMockReader(gomock.NewController(t))
			reader.EXPECT().Close().Return(nil).AnyTimes()
			reader.EXPECT().ReadMessage(ctx).Return(randMsg(), nil).AnyTimes()

			consumer := NewConsumer(Config{},
				WithKafkaReader(func() Reader { return reader }),
				WithKafkaDialer(&kafka.Dialer{}),
			)

			if tt.setup != nil {
				tt.setup(t, consumer)
			}

			handler := func(ctx context.Context, msg Message) error {
				gotAttempts++
				if wantErr != nil || gotAttempts < tt.numRetries {
					return wantHandlerErr
				}
				require.NoError(t, consumer.Stop())
				return nil
			}

			gotErr := consumer.Run(ctx, handler)
			if wantErr != nil {
				require.EqualError(t, gotErr, wantErr.Error())
			}
			assert.Equal(t, tt.shouldNotify, didNotify)
		})
	}
}

func TestNonStopExponentialBackOff(t *testing.T) {
	bo := NonStopExponentialBackOff()
	assert.Equal(t, 500*time.Millisecond, bo.NextBackOff())
	assert.Equal(t, 4*time.Second, bo.NextBackOff())
	assert.Equal(t, 32*time.Second, bo.NextBackOff())
	assert.Equal(t, (4*time.Minute)+(16*time.Second), bo.NextBackOff())
	assert.Equal(t, (34*time.Minute)+(8*time.Second), bo.NextBackOff())
	assert.Equal(t, (4*time.Hour)+(33*time.Minute)+(4*time.Second), bo.NextBackOff())
	assert.Equal(t, 5*time.Hour, bo.NextBackOff())
}

func randMsg() kafka.Message {
	return kafka.Message{
		Topic:     "some-topic",
		Partition: rand.Intn(20-1) + 1, //nolint:gosec
		Value:     []byte(uuid.New().String()),
	}
}

type testBackoff struct {
	maxAttempts int
	gotAttempts int
}

func (b *testBackoff) Reset() {}

func (b *testBackoff) NextBackOff() time.Duration {
	b.gotAttempts++
	if b.gotAttempts == b.maxAttempts {
		return backoff.Stop
	}
	return 0
}

type safeCounter struct {
	sync.RWMutex
	v int
}

func (m *safeCounter) inc() {
	m.Lock()
	defer m.Unlock()
	m.v++
}

func (m *safeCounter) val() int {
	m.RLock()
	defer m.RUnlock()
	return m.v
}
