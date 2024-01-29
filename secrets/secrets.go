package secrets

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// AWSSecretsManager supports the GetSecretValue method.
type AWSSecretsManager struct {
	Client AWSSecretsManagerClient
}

func NewAWSSecretsManager(region string) *AWSSecretsManager {
	client := newSecretManagerClient(region)
	return NewAWSSecretsManagerWithClient(client)
}

func NewAWSSecretsManagerWithClient(client AWSSecretsManagerClient) *AWSSecretsManager {
	return &AWSSecretsManager{
		Client: client,
	}
}

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
