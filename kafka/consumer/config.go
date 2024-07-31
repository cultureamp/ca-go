package consumer

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

const (
	mBytes             = 1024 * 1024
	defaultChannelSize = 256
	defaultFetchSize   = 10 * mBytes
	defaultMaxWaitTime = 500 * time.Millisecond
)

// Config is a configuration object used to create a new Consumer.
type Config struct {
	id                          string           // Default: UUID
	brokers                     []string         // Kafka bootstrap brokers to connect to
	version                     string           // Kafka cluster version (Default )
	topics                      []string         // Kafka topics to be consumed
	groupID                     string           // Kafka consumer group definition
	assignor                    string           // Consumer group partition assignment strategy (range, roundrobin, sticky)
	schemaRegistryURL           string           // The client avro registry URL
	oldest                      bool             // Kafka consumer consume initial offset from oldest (Default true)
	returnOnClientDispatchError bool             // If the receiver.dispatch returns error, then exit consume (Default false)
	tlsEnabled                  bool             // Enable TLS (Default false)
	tlsCertPem                  []byte           // TLS certificate pem
	tlsCertKey                  []byte           // TLS certificate key
	tlsCaPem                    []byte           // TLS CA pem
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
		oldest:                      true,
		returnOnClientDispatchError: false,
		tlsEnabled:                  false,
		version:                     sarama.DefaultVersion.String(),
		saramaConfig:                sarama.NewConfig(),
	}

	// set env var defaults
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers != "" {
		conf.brokers = strings.Split(brokers, ",")
	}

	topics := os.Getenv("KAFKA_TOPICS")
	if topics != "" {
		conf.topics = strings.Split(topics, ",")
	}

	schemaRegistryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	if schemaRegistryURL != "" {
		conf.schemaRegistryURL = schemaRegistryURL
	}

	certPem := os.Getenv("KAFKA_CERT_PEM")
	if certPem != "" {
		conf.tlsCertPem = []byte(certPem)
	}

	certKey := os.Getenv("KAFKA_CERT_KEY")
	if certKey != "" {
		conf.tlsCertKey = []byte(certKey)
	}

	caPem := os.Getenv("KAFKA_CA_PEM")
	if caPem != "" {
		conf.tlsCaPem = []byte(caPem)
	}

	conf.saramaConfig.ChannelBufferSize = 256
	conf.saramaConfig.Consumer.Fetch.Default = defaultFetchSize
	conf.saramaConfig.Consumer.IsolationLevel = sarama.ReadCommitted
	conf.saramaConfig.Consumer.MaxWaitTime = defaultMaxWaitTime
	conf.saramaConfig.Consumer.Offsets.AutoCommit.Enable = true

	// ConsumerGroup <- Errors returns a read channel of errors that occurred during the consumer life-cycle.
	// By default, errors are logged and not returned over this channel.
	// If you want to implement any custom error handling, set your config's
	// Consumer.Return.Errors setting to true, and read from this channel.
	// conf.saramaConfig.Consumer.Return.Errors = true

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

	if conf.groupID == "" {
		return errors.Errorf("missing group")
	}

	err := conf.shouldProcessAssignor()
	if err != nil {
		return err
	}

	version, err := sarama.ParseKafkaVersion(conf.version)
	if err != nil {
		return errors.Errorf("invalid kafka version '%s': %w", conf.version, err)
	}
	conf.saramaConfig.Version = version

	if conf.schemaRegistryURL == "" {
		return errors.Errorf("missing schema registry URL")
	}

	if conf.oldest {
		conf.saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	err = conf.shouldProcessTLS()
	if err != nil {
		return err
	}

	return nil
}

func (conf *Config) shouldProcessAssignor() error {
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

	return nil
}

func (conf *Config) shouldProcessTLS() error {
	if conf.tlsEnabled {
		cert, err := tls.X509KeyPair(conf.tlsCertPem, conf.tlsCertKey)
		if err != nil {
			return errors.Errorf("error loading client certificate: %w", err)
		}

		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM(conf.tlsCaPem)
		if !ok {
			return errors.Errorf("failed to parse root certificate")
		}

		tlsConf := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
			MinVersion:   tls.VersionTLS12,
		}

		conf.saramaConfig.Net.TLS.Enable = true
		conf.saramaConfig.Net.TLS.Config = tlsConf
	}

	return nil
}
