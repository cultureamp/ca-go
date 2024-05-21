package consumer

import (
	"io"
	"log"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

// Config is a configuration object used to create a new Consumer.
type Config struct {
	id                          string           // Default: UUID
	brokers                     []string         // Kafka bootstrap brokers to connect to
	version                     string           // Kafka cluster version (Default )
	topics                      []string         // Kafka topics to be consumed
	groupId                     string           // Kafka consumer group definition
	assignor                    string           // Consumer group partition assignment strategy (range, roundrobin, sticky)
	handler                     Handler          // The client handler to receive and process messages
	oldest                      bool             // Kafka consumer consume initial offset from oldest (Default true)
	returnOnClientDispatchError bool             // If the receiver.dispatch returns error, then exit consume (Default false)
	stdLogger                   sarama.StdLogger // Consumer logging (Default nil)
	debugLogger                 sarama.StdLogger // Sarama logger (Default nil)
	saramaConfig                *sarama.Config
}

func newConfig() *Config {
	// set defaults
	conf := &Config{
		id:                          uuid.New().String(),
		stdLogger:                   log.New(io.Discard, "", log.LstdFlags),
		debugLogger:                 log.New(io.Discard, "", log.LstdFlags),
		handler:                     nil,
		oldest:                      true,
		returnOnClientDispatchError: false,
		version:                     sarama.DefaultVersion.String(),
		saramaConfig:                sarama.NewConfig(),
	}

	// ConsumerGroup <- Errors returns a read channel of errors that occurred during the consumer life-cycle.
	// By default, errors are logged and not returned over this channel.
	// If you want to implement any custom error handling, set your config's
	// Consumer.Return.Errors setting to true, and read from this channel.

	// conf.saramaConfig.Consumer.Return.Errors = true
	conf.saramaConfig.Consumer.Offsets.AutoCommit.Enable = true
	conf.saramaConfig.Consumer.IsolationLevel = sarama.ReadCommitted

	return conf
}

func (conf *Config) shouldProcess() error {
	conf.saramaConfig.ClientID = conf.id

	if conf.stdLogger != nil {
		// sarama.Logger is a global package level variable
		sarama.Logger = conf.stdLogger
	}

	if conf.debugLogger != nil {
		// sarama.DebugLogger is a global package level variable
		sarama.DebugLogger = conf.debugLogger
	}

	if len(conf.brokers) == 0 {
		return errors.Errorf("missing brokers")
	}

	if len(conf.topics) == 0 {
		return errors.Errorf("missing topics")
	}

	if conf.groupId == "" {
		return errors.Errorf("missing group")
	}

	version, err := sarama.ParseKafkaVersion(conf.version)
	if err != nil {
		return errors.Errorf("invalid kafka version '%s': %w", conf.version, err)
	}
	conf.saramaConfig.Version = version

	switch conf.assignor {
	case "sticky":
		conf.saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		conf.saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range", "":
		conf.saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		// can we default to avoid this error?
		return errors.Errorf("unrecognized consumer group partition assignor: %s", conf.assignor)
	}

	if conf.handler == nil {
		return errors.Errorf("missing message handler")
	}

	if conf.oldest {
		conf.saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	return nil
}
