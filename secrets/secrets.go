package secrets

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWSSecretsManager supports wraps the Secrets interface.
type AWSSecretsManager struct {
	smClient *secretsmanager.Client
}

// NewAWSSecretsManager creates a new AWS Secret Manager for a given region.
func NewAWSSecretsManager(ctx context.Context, region string) (*AWSSecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	smc := secretsmanager.NewFromConfig(cfg)
	return &AWSSecretsManager{smClient: smc}, nil
}

// NewAWSSecretsManagerWithClient creates a new AWS Secret Manager with a custom client
// that supports the AWSSecretsManagerClient interface.
func NewAWSSecretsManagerWithClient(client *secretsmanager.Client) *AWSSecretsManager {
	return &AWSSecretsManager{
		smClient: client,
	}
}

// Get retrieves the secret from AWS SecretsManager.
func (sm *AWSSecretsManager) Get(ctx context.Context, secretKey string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretKey),
	}

	result, err := sm.smClient.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}
