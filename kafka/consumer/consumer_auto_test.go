package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestAutoConsumer(t *testing.T) {
	const queueSize = 5

	testRunnerKafkaReader := func() Reader {
		return newTestRunnerReader("topic-name")
	}
	autoBackOff := NonStopExponentialBackOff
	autoBalancers := []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}}

	defConsumers := DefaultConsumers
	DefaultConsumers = newPackageAutoConsumerGlue(
		WithGroupBalancers(autoBalancers...),
		WithHandlerBackOffRetry(autoBackOff),
		WithKafkaReader(testRunnerKafkaReader),
		WithQueueCapacity(queueSize-1),
	)
	defer func() {
		DefaultConsumers = defConsumers
	}()

	ctx := context.Background()
	ch, stop := Consume(ctx, "topic-name")

	// wait for the channel to fill up with 10 messages (our queue capacity)
	time.Sleep(time.Millisecond * 500)
	// stop the consumer
	stop()

	// we should only read queueSize messages before the channel signals it has been closed
	i := 0
	for range ch {
		i++
	}
	assert.Equal(t, queueSize, i)
}

type packageAutoMock struct {
	brokers []string
	opts    []Option
}

func newPackageAutoConsumerGlue(opts ...Option) *packageAutoMock {
	return &packageAutoMock{
		brokers: []string{"test1", "test2"},
		opts:    opts,
	}
}

func (c *packageAutoMock) Consume(ctx context.Context, topic string) (<-chan Message, StopFunc) {
	ac := newAutoConsumer(topic, c.brokers, c.opts...)
	go ac.run(ctx)
	return ac.channel, func() error { return ac.stop() }
}
