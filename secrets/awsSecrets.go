package secrets

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWSSecrets supports the Get method.
type AWSSecrets struct {
	smClient *secretsmanager.Client
}

func NewAWSSecrets(region string) *AWSSecrets {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		//todo
	}
	smc := secretsmanager.NewFromConfig(cfg)

	return &AWSSecrets{smClient: smc}
}

func (s *AWSSecrets) Get(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := s.smClient.GetSecretValue(context.Background(), input)
	if err != nil {
		return "", err
	}

	// Assuming the secret is a string
	return *result.SecretString, nil
}
