package secrets

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSSecretsManager(t *testing.T) {
	ctx := context.Background()
	secrets, err := NewAWSSecretsManager(ctx, "us-west-2")
	assert.Nil(t, err)
	assert.NotNil(t, secrets)

	// with own client
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	smc := secretsmanager.NewFromConfig(cfg)
	secrets = NewAWSSecretsManagerWithClient(smc)
	assert.NotNil(t, secrets)
}
