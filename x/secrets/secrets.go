package secrets

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretClient interface {
	Get(secretName string) (string, error)
}

type AWSSecrets struct {
	Client interface {
		GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
	}
}

func NewAWSSecretsClient(region string) (*AWSSecrets, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	client := secretsmanager.New(sess)
	return &AWSSecrets{
		Client: client,
	}, nil
}

func (s *AWSSecrets) Get(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := s.Client.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s: %w", secretName, err)
	}
	if result == nil || result.SecretString == nil {
		return "", fmt.Errorf("retrieved secret %s is empty", secretName)
	}
	return *result.SecretString, nil
}
