package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAwsSettings(t *testing.T) {
	settings := newAWSSettings()
	assert.NotNil(t, settings)
}

func TestAwsSettings(t *testing.T) {
	t.Setenv(AwsProfileEnv, "dev")
	t.Setenv(AwsRegionEnv, "us-west-1")
	t.Setenv(AwsAccountIDEnv, "123456789")
	t.Setenv(AwsXrayEnv, "true")

	settings := newAWSSettings()
	assert.Equal(t, "dev", settings.AwsProfile)
	assert.Equal(t, "us-west-1", settings.AwsRegion)
	assert.Equal(t, "123456789", settings.AwsAccountID)
	assert.Equal(t, true, settings.XrayLogging)
}
