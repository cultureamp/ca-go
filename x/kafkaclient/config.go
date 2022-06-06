package kafkaclient

import (
	"crypto/tls"
	"github.com/Shopify/sarama"
)

type Config struct {
	Brokers  string
	Username string
	Password string
	Topic    string
	ClientID string
}

func GetConnConfig(cfg Config) *sarama.Config {
	conf := sarama.NewConfig()

	conf.Producer.RequiredAcks = sarama.WaitForAll
	// Return.Successes and Return.Errors specify what channels will be populated.
	// If this config is used to create a `SyncProducer`, both must be set to true
	// and you shall not read from the channels since the producer does this internally.
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	// A user-provided string sent with every request to the brokers for logging, debugging, and auditing purposes.
	// Defaults to "sarama", should set it to be something specific to your application.
	conf.ClientID = cfg.ClientID

	// SASL must be enabled to authenticate with kafka using username and password
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = cfg.Username
	conf.Net.SASL.Password = cfg.Password

	// enable SHA512 and TLS as required according to our internal Kafka config
	conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = createTLSConfiguration()

	// verbose debugging (uncomment this line to help debugging)
	//sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	return conf
}

func createTLSConfiguration() (t *tls.Config) {
	return &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}
}
