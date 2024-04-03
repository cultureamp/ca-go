package consumer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			wantOpts := []Option{WithHandlerBackOffRetry(wantBackOff), WithKafkaDialer(&kafka.Dialer{})}

			g := NewGroup(tt.config, wantOpts...)
			require.NotNil(t, g)
			assert.Empty(t, g.stopChs)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wantConsumers := 10
	wantNumMsgs := 100

	reader := NewMockReader(gomock.NewController(t))
	reader.EXPECT().Close().AnyTimes()
	reader.EXPECT().ReadMessage(ctx).Return(randMsg(), nil).AnyTimes()

	group := NewGroup(GroupConfig{Count: wantConsumers},
		WithKafkaReader(func() Reader { return reader }),
		WithKafkaDialer(&kafka.Dialer{}),
	)

	handlerInvocations := new(safeCounter)
	stopCh := make(chan bool)

	handler := func(ctx context.Context, msg Message) error {
		require.NotEmpty(t, msg.Value)
		handlerInvocations.inc()
		if handlerInvocations.val() == wantNumMsgs {
			stopCh <- true
		}
		return nil
	}

	errCh := group.Run(ctx, handler)
	for {
		select {
		case <-stopCh:
			return
		case err := <-errCh:
			require.NoError(t, err)
		}
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
			},
		},
		{
			name:               "reader fetch error",
			withExplicitCommit: true,
			setupReader: func(m *MockReader) {
				m.EXPECT().FetchMessage(gomock.Any()).Return(kafka.Message{}, wantErr).Times(1)
			},
		},
		{
			name:               "reader commit error",
			withExplicitCommit: true,
			setupReader: func(m *MockReader) {
				m.EXPECT().FetchMessage(gomock.Any()).Return(wantMsg, nil).Times(1)
				m.EXPECT().CommitMessages(gomock.Any(), wantMsg).Return(wantErr).Times(1)
			},
		},
		{
			name:     "reader close error",
			closeErr: true,
			setupReader: func(m *MockReader) {
				m.EXPECT().ReadMessage(gomock.Any()).Return(kafka.Message{}, wantErr).Times(1)
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
				require.ErrorIs(t, err, wantErr)
				group.Stop()
				require.Contains(t, err.Error(), wantErr.Error())
			}
		})
	}
}
