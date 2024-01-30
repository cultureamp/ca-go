package secrets

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type awsSecrets struct {
	smClient *secretsmanager.Client
}

// NewAWSSecrets returns an instance of the Secrets interface using
// AWS secret manager client configured using the given region
func NewAWSSecrets(region string) Secrets {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		fmt.Printf("error loading aws sdk config, err='%v'\n", err)
	}
	smc := secretsmanager.NewFromConfig(cfg)

	return &awsSecrets{smClient: smc}
}

// Get returns the secret for the given id or the error encountered when trying to retrieve it
func (s *awsSecrets) Get(secretID string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	result, err := s.smClient.GetSecretValue(context.Background(), input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}
