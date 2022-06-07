package kafkaclient_test

import (
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/cultureamp/ca-go/x/sarama/kafkaclient"
	"github.com/stretchr/testify/assert"
)

//nolint:ifshort
func Example() {
	// the default behaviour can be overridden in local environments
	localEnv := false

	config := kafkaclient.DefaultProducerConfiguration("application_name", "username", "passwordvaluefromsecret")

	if localEnv {
		config.Net.TLS.Enable = false
		config.Net.SASL.Enable = false
	}

	fmt.Printf("ClientID=%s TLS:enabled=%t SASL:enabled=%t\n", config.ClientID, config.Net.TLS.Enable, config.Net.SASL.Enable)

	// Output: ClientID=application_name TLS:enabled=true SASL:enabled=true
}

func ExampleDefaultProducerConfiguration_production() {
	config := kafkaclient.DefaultProducerConfiguration("application_name", "username", "passwordvaluefromsecret")

	fmt.Printf("ClientID=%s TLS:enabled=%t SASL:enabled=%t\n", config.ClientID, config.Net.TLS.Enable, config.Net.SASL.Enable)

	// Output: ClientID=application_name TLS:enabled=true SASL:enabled=true
}

func ExampleDefaultProducerConfiguration_testingoverrides() {
	config := kafkaclient.DefaultProducerConfiguration("application_name", "username", "passwordvaluefromsecret")

	// for testing, disable TLS and SASL when connecting to a local Kafka instance
	config.Net.TLS.Enable = false
	config.Net.SASL.Enable = false

	fmt.Printf("ClientID=%s TLS:enabled=%t SASL:enabled=%t\n", config.ClientID, config.Net.TLS.Enable, config.Net.SASL.Enable)

	// Output: ClientID=application_name TLS:enabled=false SASL:enabled=false
}

func TestDefaultProducerConfiguration(t *testing.T) {
	config := kafkaclient.DefaultProducerConfiguration("test_client_tag", "user", "password")

	assert.Equal(t, "test_client_tag", config.ClientID)
	assert.Equal(t, sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA512), config.Net.SASL.Mechanism)
	assert.Equal(t, "user", config.Net.SASL.User)
	assert.Equal(t, "password", config.Net.SASL.Password)
	assert.True(t, config.Net.TLS.Enable)
	assert.True(t, config.Net.SASL.Enable)
}
