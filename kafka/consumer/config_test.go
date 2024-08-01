package consumer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigShouldProcess(t *testing.T) {
	conf := newConfig()
	assert.NotNil(t, conf)

	err := conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing brokers")

	conf.brokers = []string{"localhost:9092"}
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing topic")

	conf.topics = []string{"test-topic"}
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing group")

	conf.groupID = "group_id"
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing schema registry URL")

	conf.schemaRegistryURL = "http://localhost:8081"
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing sasl username")

	conf.saramaConfig.Net.SASL.User = "test_user"
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing sasl password")

	conf.saramaConfig.Net.SASL.Password = "test_pwd"
	err = conf.shouldProcess()
	assert.Nil(t, err)

	// valid options are: sticky, range, roundrobin
	conf.assignor = "abc"
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unrecognized consumer group partition assignor")

	conf.assignor = "sticky"
	err = conf.shouldProcess()
	assert.Nil(t, err)

	conf.version = "a wacky number"
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid version")

	conf.version = "1.0.0"
	err = conf.shouldProcess()
	assert.Nil(t, err)
}

func TestConfigShouldProcessWithEnvVars(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", "localhost:9092,localhost:9093")
	t.Setenv("KAFKA_TOPICS", "test-topic1, test_topic2")
	t.Setenv("SCHEMA_REGISTRY_URL", "http://test.schema.registry.com:8081")
	conf := newConfig()
	assert.NotNil(t, conf)

	conf.saramaConfig.Net.SASL.Enable = false
	conf.groupID = "group_id"

	err := conf.shouldProcess()
	assert.Nil(t, err)
}
