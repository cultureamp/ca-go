package secrets

import (
	"time"

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

type testRunnerClient struct{}

func newTestRunnerClient() *testRunnerClient {
	return &testRunnerClient{}
}

func (c *testRunnerClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	arn := "arn:aws:secretmanager:eu-west-2:abc123:secret/id"
	now := time.Now()

	retVal := &secretsmanager.GetSecretValueOutput{
		ARN:          &arn,
		CreatedDate:  &now,
		Name:         input.SecretId,
		SecretString: input.SecretId, // just echo back the key as the secret when running in a test
		VersionId:    input.VersionId,
	}
	return retVal, nil
}
