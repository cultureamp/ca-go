package consumer

import (
	"context"
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

	conf.groupId = "group_id"
	err = conf.shouldProcess()
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "missing message handler")

	conf.handler = func(context.Context, *Message) error { return nil }
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
