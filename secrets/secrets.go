package secrets

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
)

// AWSSecretsManager supports wraps the AWSSecretsManagerClient interface.
type AWSSecretsManager struct {
	Client AWSSecretsManagerClient
}

// NewAWSSecretsManager creates a new AWS Secret Manager for a given region.
func NewAWSSecretsManager(ctx context.Context, region string) (*AWSSecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := newSecretManagerClient(cfg)
	return NewAWSSecretsManagerWithClient(client), nil
}

// NewAWSSecretsManagerWithClient creates a new AWS Secret Manager with a custom client
// that supports the AWSSecretsManagerClient interface.
func NewAWSSecretsManagerWithClient(client AWSSecretsManagerClient) *AWSSecretsManager {
	return &AWSSecretsManager{
		Client: client,
	}
}

// Get retrieves the secret from AWS SecretsManager.
func (s *AWSSecretsManager) Get(ctx context.Context, secretKey string) (string, error) {
	secret, err := s.Client.GetSecretValue(ctx, secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve '%s': %w", secretKey, err)
	}
	if secret == "" {
		return "", fmt.Errorf("retrieved secret '%s' is empty", secretKey)
	}
	return secret, nil
}
