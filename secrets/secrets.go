package secrets

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// AWSSecretsManager supports wraps the AWSSecretsManagerClient interface.
type AWSSecretsManager struct {
	Client AWSSecretsManagerClient
}

// NewAWSSecretsManager creates a new AWS Secret Manager for a given region.
func NewAWSSecretsManager(region string) *AWSSecretsManager {
	client := newSecretManagerClient(region)
	return NewAWSSecretsManagerWithClient(client)
}

// NewAWSSecretsManagerWithClient creates a new AWS Secret Manager with a custom client
// that supports the AWSSecretsManagerClient interface.
func NewAWSSecretsManagerWithClient(client AWSSecretsManagerClient) *AWSSecretsManager {
	return &AWSSecretsManager{
		Client: client,
	}
}

// Get retrives the secret from AWS SecretsManager.
func (s *AWSSecretsManager) Get(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := s.Client.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve '%s': %w", secretName, err)
	}
	if result == nil || result.SecretString == nil {
		return "", fmt.Errorf("retrieved secret '%s' is empty", secretName)
	}
	return *result.SecretString, nil
}
