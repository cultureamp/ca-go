package secrets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// AWSSecretsManagerClient can be mocked by clients for testing purposes.
type AWSSecretsManagerClient interface {
	GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

func newSecretManagerClient(region string) *secretsmanager.SecretsManager {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return secretsmanager.New(sess)
}

type testClient struct {

}

func newTestClient() *testClient {
	return &testClient{}
}

func (c *testClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	retVal := &secretsmanager.GetSecretValueOutput{
		SecretString: input.SecretId,
	}

	return retVal, nil
}