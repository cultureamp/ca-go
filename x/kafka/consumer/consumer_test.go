package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	kafkatrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/segmentio/kafka.go.v0"
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
				GroupID: "some-group-id",
			},
		},
		{
			name: "new group with invalid number of consumers defaults to 1",
			config: Config{
				Brokers: []string{"some-address"},
				Topic:   "some-topic",
				GroupID: "some-group-id",
			},
			wantGenID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBackOff := NonStopExponentialBackOff
			wantNotify := func(ctx context.Context, err error, msg Message) {}
			wantBalancers := []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}}

			dialer, err := DialerSCRAM512("username", "password")
			require.NoError(t, err)

			c := NewConsumer(dialer, tt.config,
				WithExplicitCommit(),
				WithGroupBalancers(wantBalancers...),
				WithHandlerBackOffRetry(wantBackOff),
				WithNotifyError(wantNotify),
				WithReaderLogger(func(s string, i ...interface{}) { log.Println(s) }),
				WithReaderErrorLogger(func(s string, i ...interface{}) { log.Println(s) }),
				WithDataDogTracing(),
			)
			require.NotNil(t, c)
			assert.Equal(t, c.readerConfig.Dialer, dialer)
			assert.NotNil(t, c.reader)
			assert.IsType(t, &kafkatrace.Reader{}, c.reader)
			assert.NotNil(t, c.backOffConstructor)
			assert.NotNil(t, c.notifyErr)
			assert.Equal(t, wantBalancers, c.readerConfig.GroupBalancers)

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

	consumer := &Consumer{reader: reader}

	i := 0
	handler := func(ctx context.Context, msg Message) error {
		require.Equal(t, currMsg, msg.Message)
		i++
		if i == wantTimes {
			require.NoError(t, consumer.Close())
		}
		return nil
	}

	require.NoError(t, consumer.Run(ctx, handler))
}

func TestConsumer_Run_error(t *testing.T) {
	var wantErr error
	var gotAttempts int
	var didNotify bool

	wantConsumerID := "123"
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
			wantError:    fmt.Errorf("consumer %s unable to handle message: %w", wantConsumerID, wantHandlerErr),
			shouldNotify: false,
			numRetries:   0,
		},
		{
			name:         "handler error after backoff retry and notify",
			wantError:    fmt.Errorf("consumer %s unable to handle message: %w", wantConsumerID, wantHandlerErr),
			shouldNotify: true,
			numRetries:   3,
			setup: func(t *testing.T, consumer *Consumer) {
				consumer.backOffConstructor = func() backoff.BackOff {
					return &testBackoff{
						maxAttempts: 3,
					}
				}
				consumer.notifyErr = func(ctx context.Context, err error, msg Message) {
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
				consumer.backOffConstructor = func() backoff.BackOff {
					return &testBackoff{}
				}
				consumer.notifyErr = func(ctx context.Context, err error, msg Message) {
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
			wantError:        fmt.Errorf("consumer %s unable to handle message: context canceled", wantConsumerID),
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

			consumer := &Consumer{
				id:     wantConsumerID,
				reader: reader,
			}
			if tt.setup != nil {
				tt.setup(t, consumer)
			}

			handler := func(ctx context.Context, msg Message) error {
				gotAttempts++
				if wantErr != nil || gotAttempts < tt.numRetries {
					return wantHandlerErr
				}
				consumer.Close()
				return nil
			}

			gotErr := consumer.Run(ctx, handler)
			if wantErr != nil {
				assert.EqualError(t, gotErr, wantErr.Error())
			} else {
				assert.NoError(t, consumer.Close())
			}
			assert.Equal(t, tt.shouldNotify, didNotify)
		})
	}
}

func TestNewGroup(t *testing.T) {
	tests := []struct {
		name             string
		config           GroupConfig
		wantNumConsumers int
	}{
		{
			name: "new group with valid number of consumers",
			config: GroupConfig{
				Count:   10,
				Brokers: []string{"some-address"},
				Topic:   "some-topic",
				GroupID: "some-group-id",
			},
			wantNumConsumers: 10,
		},
		{
			name: "new group with invalid number of consumers defaults to 1",
			config: GroupConfig{
				Count:   0,
				Brokers: []string{"some-address"},
				Topic:   "some-topic",
				GroupID: "some-group-id",
			},
			wantNumConsumers: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantBackOff := NonStopExponentialBackOff
			wantNotify := func(ctx context.Context, err error, msg Message) {}
			wantOpts := []Option{WithHandlerBackOffRetry(wantBackOff), WithNotifyError(wantNotify)}

			g := NewGroup(&kafka.Dialer{}, tt.config, wantOpts...)
			require.NotNil(t, g)
			assert.Empty(t, g.consumers)
			wantConfig := tt.config
			wantConfig.Count = tt.wantNumConsumers
			assert.Equal(t, wantConfig, g.config)
			assert.Len(t, g.opts, len(wantOpts))
		})
	}
}

// TestGroup_Run tests that consumers can concurrently handle messages received
// from a Reader. On the other hand, it does not test if the allocation of
// partitions to consumers are mutually exclusive, since that is the responsibility
// of Kafka itself.
func TestGroup_Run(t *testing.T) {
	ctx := context.Background()
	wantConsumers := 10
	wantNumMsgs := 100

	// Concurrent safe msg tracker since handler runs from multiple consumer go routines.
	var msgTracker messageTracker

	reader := NewMockReader(gomock.NewController(t))
	reader.EXPECT().Close().AnyTimes()
	reader.EXPECT().ReadMessage(ctx).DoAndReturn(func(_ context.Context) (kafka.Message, error) {
		msgTracker.Lock()
		defer msgTracker.Unlock()
		msg := randMsg()
		msgTracker.wantMessages = append(msgTracker.wantMessages, msg)
		return msg, nil
	}).AnyTimes()

	group := &Group{
		config: GroupConfig{Count: wantConsumers},
		opts: []Option{WithKafkaReader(func() Reader {
			return reader
		})},
	}

	handler := func(ctx context.Context, msg Message) error {
		msgTracker.Lock()
		defer msgTracker.Unlock()
		require.Contains(t, msgTracker.wantMessages, msg.Message)
		if len(msgTracker.wantMessages) == wantNumMsgs {
			require.NoError(t, group.Close())
		}
		return nil
	}

	errCh := group.Run(ctx, handler)
	for err := range errCh {
		require.NoError(t, err)
	}
}

func TestGroup_Run_readerError(t *testing.T) {
	wantErr := errors.New("some error")
	wantConsumers := 1
	wantMsg := randMsg()

	tests := []struct {
		name               string
		withExplicitCommit bool
		setupReader        func(m *MockReader)
		closeErr           bool
	}{
		{
			name: "reader read error",
			setupReader: func(m *MockReader) {
				m.EXPECT().ReadMessage(gomock.Any()).Return(kafka.Message{}, wantErr).Times(1)
				m.EXPECT().Close().Return(nil).Times(1)
			},
		},
		{
			name:               "reader fetch error",
			withExplicitCommit: true,
			setupReader: func(m *MockReader) {
				m.EXPECT().FetchMessage(gomock.Any()).Return(kafka.Message{}, wantErr).Times(1)
				m.EXPECT().Close().Return(nil).Times(1)
			},
		},
		{
			name:               "reader commit error",
			withExplicitCommit: true,
			setupReader: func(m *MockReader) {
				m.EXPECT().FetchMessage(gomock.Any()).Return(wantMsg, nil).Times(1)
				m.EXPECT().CommitMessages(gomock.Any(), wantMsg).Return(wantErr).Times(1)
				m.EXPECT().Close().Return(nil).Times(1)
			},
		},
		{
			name:     "reader close error",
			closeErr: true,
			setupReader: func(m *MockReader) {
				m.EXPECT().ReadMessage(gomock.Any()).Return(kafka.Message{}, wantErr).Times(1)
				m.EXPECT().Close().Return(wantErr).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewMockReader(gomock.NewController(t))
			tt.setupReader(reader)

			group := &Group{
				config: GroupConfig{
					Count: wantConsumers,
				},
				opts: []Option{WithKafkaReader(func() Reader {
					return reader
				})},
			}

			if tt.withExplicitCommit {
				group.opts = append(group.opts, WithExplicitCommit())
			}
			errCh := group.Run(context.Background(), func(ctx context.Context, msg Message) error {
				return nil
			})
			for err := range errCh {
				require.True(t, errors.Is(err, wantErr))
				err = group.Close()
				if tt.closeErr {
					require.Contains(t, err.Error(), wantErr.Error())
					return
				}
				require.NoError(t, err)
			}
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

type messageTracker struct {
	sync.Mutex
	wantMessages []kafka.Message
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