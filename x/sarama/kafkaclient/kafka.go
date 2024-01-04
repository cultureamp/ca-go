package kafkaclient

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/IBM/sarama"
)

// EnableDebugLog provides more details on how sarama connecting to brokers.
func EnableDebugLog() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
}

// DefaultProducerConfiguration creates a Sarama configuration with the default settings
// appropriate for use by a Sarama SyncProducer in a SASL environment.
func DefaultProducerConfiguration(clientID string, username string, password string) *sarama.Config {
	conf := sarama.NewConfig()

	conf.Producer.RequiredAcks = sarama.WaitForAll

	// Return.Successes and Return.Errors specify what channels will be populated.
	// If this config is used to create a `SyncProducer`, both must be set to true
	// and you shall not read from the channels since the producer does this internally.
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true

	// A user-provided string sent with every request to the brokers for logging, debugging, and auditing purposes.
	// Defaults to "sarama", should set it to be something specific to your application.
	conf.ClientID = clientID

	// SASL must be enabled to authenticate with kafka using username and password
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = username
	conf.Net.SASL.Password = password

	// enable SHA512 and TLS as required according to our internal Kafka config
	conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &xDGSCRAMClient{HashGenerator: SHA512} }
	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	return conf
}

// NonSSLProducerConfiguration creates a Sarama configuration with the default settings without SSL config
// appropriate for use to test connection with local Kafka.
func NonSSLProducerConfiguration(clientID string) *sarama.Config {
	conf := sarama.NewConfig()

	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	conf.ClientID = clientID

	return conf
}
