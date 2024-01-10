package secrets

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// smClient used for mock testing.
type smClient interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

// AWSSecrets supports the GetSecretValue method.
type AWSSecrets struct {
	client smClient
}

var defaultAWSSecrets = getInstance()

func getInstance() *AWSSecrets {
	// Should this take dependency on 'env' package and call env.AwsRegion()?
	region := os.Getenv("AWS_REGION")
	return NewAWSSecrets(region)
}

func NewAWSSecrets(region string) *AWSSecrets {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	client := secretsmanager.New(sess)
	return &AWSSecrets{
		client: client,
	}
}

// Get retrives the secret from AWS SecretsManager.
func Get(secretName string) (string, error) {
	return defaultAWSSecrets.Get(secretName)
}

func (s *AWSSecrets) Get(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := s.client.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve '%s': %w", secretName, err)
	}
	if result == nil || result.SecretString == nil {
		return "", fmt.Errorf("retrieved secret '%s' is empty", secretName)
	}
	return *result.SecretString, nil
}
