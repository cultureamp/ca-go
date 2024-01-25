package secrets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAWSSecretsClient(t *testing.T) {
	client := NewAWSSecrets("us-west-2")
	assert.NotNil(t, client)
}
