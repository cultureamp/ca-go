package consumer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
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
			wantNotify := func(ctx context.Context, err error, msg kafka.Message, md Metadata) {}

			dialer, err := DialerSCRAM512("username", "password")
			require.NoError(t, err)

			c := NewConsumer(dialer, tt.config,
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

			if tt.wantGenID {
				assert.NotEmpty(t, c.id)
			} else {
				assert.Equal(t, tt.config.ID, c.id)
			}
		})
	}
}

func TestConsumer_Run(t *testing.T) {
	wantMsgs := dummyMessages(100, 20)

	consumer := &Consumer{
		id:     "1",
		reader: newMockReader(wantMsgs),
	}

	var gotMsgs []kafka.Message
	handler := func(ctx context.Context, msg kafka.Message, md Metadata) error {
		gotMsgs = append(gotMsgs, msg)
		return nil
	}

	require.NoError(t, consumer.Run(context.Background(), handler))
	require.Equal(t, wantMsgs, gotMsgs)
	require.NoError(t, consumer.Close())
}

func TestConsumer_Run_error(t *testing.T) {
	var wantErr error
	var gotAttempts int
	var didNotify bool

	wantConsumerID := "123"
	wantHandlerErr := errors.New("some downstream error")
	wantMsg := kafka.Message{Value: []byte(uuid.New().String())}

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
				consumer.notifyErr = func(ctx context.Context, err error, msg kafka.Message, md Metadata) {
					assert.Equal(t, gotAttempts, md.Attempt)
					assert.Equal(t, wantHandlerErr, err)
					assert.Equal(t, wantMsg, msg)
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
				consumer.notifyErr = func(ctx context.Context, err error, msg kafka.Message, md Metadata) {
					assert.Equal(t, gotAttempts, md.Attempt)
					if err != nil {
						assert.Equal(t, wantHandlerErr, err)
					}
					assert.Equal(t, wantMsg, msg)
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

			consumer := &Consumer{
				id:     wantConsumerID,
				reader: newMockReader([]kafka.Message{wantMsg}),
			}
			if tt.setup != nil {
				tt.setup(t, consumer)
			}

			handler := func(ctx context.Context, msg kafka.Message, md Metadata) error {
				gotAttempts++
				if wantErr != nil || gotAttempts < tt.numRetries {
					return wantHandlerErr
				}
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
			wantNotify := func(ctx context.Context, err error, msg kafka.Message, md Metadata) {}
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
	wantMsgs := dummyMessages(1000, 20)
	wantConsumers := 10

	readerFn := func() Reader {
		return newMockReader(wantMsgs)
	}
	group := &Group{
		config: GroupConfig{
			Count: wantConsumers,
		},
		opts: []Option{WithKafkaReader(readerFn)},
	}

	// Concurrent safe counter since handler runs from multiple consumer go routines.
	var count syncCounter

	handler := func(ctx context.Context, msg kafka.Message, md Metadata) error {
		count.Lock()
		defer count.Unlock()
		count.val++
		return nil
	}

	errCh := group.Run(context.Background(), handler)
	for err := range errCh {
		require.NoError(t, err)
	}

	wantCount := len(wantMsgs) * wantConsumers
	require.Equal(t, wantCount, count.val)
	require.NoError(t, group.Close())
}

func TestGroup_Run_readerError(t *testing.T) {
	wantFetchErr := errors.New("error fetching message")
	wantCommitErr := errors.New("error committing offset")
	wantCloseErr := errors.New("error closing reader")
	wantConsumers := 1

	handler := func(ctx context.Context, msg kafka.Message, md Metadata) error {
		return nil
	}

	tests := []struct {
		name          string
		reader        mockReader
		wantFetchErr  error
		wantCommitErr error
		wantCloseErr  error
	}{
		{
			name:         "reader fetch error",
			reader:       mockReader{fetchErr: wantFetchErr},
			wantFetchErr: wantFetchErr,
		},
		{
			name:          "reader commit error",
			reader:        mockReader{commitErr: wantCommitErr},
			wantCommitErr: wantCommitErr,
		},
		{
			name:         "reader close error",
			reader:       mockReader{closeErr: wantCloseErr},
			wantCloseErr: wantCloseErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &Group{
				config: GroupConfig{
					Count: wantConsumers,
				},
				opts: []Option{WithKafkaReader(func() Reader {
					return &tt.reader
				})},
			}

			errCh := group.Run(context.Background(), handler)
			for err := range errCh {
				if tt.wantFetchErr != nil {
					require.Contains(t, err.Error(), tt.wantFetchErr.Error())
				}
				if tt.wantCommitErr != nil {
					require.Contains(t, err.Error(), tt.wantCommitErr.Error())
				}
			}
			err := group.Close()
			if tt.wantCloseErr != nil {
				require.Contains(t, err.Error(), tt.wantCloseErr.Error())
			} else {
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

type mockReader struct {
	msgs      []kafka.Message
	fetchErr  error
	commitErr error
	closeErr  error
}

func newMockReader(wantMsgs []kafka.Message) *mockReader {
	return &mockReader{
		msgs: wantMsgs,
	}
}

func (m *mockReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	if m.fetchErr != nil {
		return kafka.Message{}, m.fetchErr
	}
	if len(m.msgs) == 0 {
		return kafka.Message{}, io.EOF
	}

	msg := m.msgs[0]
	m.msgs = m.msgs[1:]
	return msg, nil
}

func (m *mockReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	if m.commitErr != nil {
		return m.commitErr
	}
	return nil
}
func (m *mockReader) Close() error {
	if m.closeErr != nil {
		return m.closeErr
	}
	return nil
}

func dummyMessages(n int, partitions int) []kafka.Message {
	rand.Seed(time.Now().UnixNano())
	var wantMsgs []kafka.Message

	for i := 0; i < n; i++ {
		partition := 1
		if partitions > 1 {
			partition = rand.Intn(partitions-1) + 1 //nolint:gosec
		}
		wantMsgs = append(wantMsgs, kafka.Message{
			Topic:     "some-topic",
			Partition: partition,
			Value:     []byte(uuid.New().String()),
		})
	}

	return wantMsgs
}

type syncCounter struct {
	sync.Mutex
	val int
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
