package secrets

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWSSecretsManagerClient can be mocked by clients for testing purposes.
type AWSSecretsManagerClient interface {
	GetSecretValue(ctx context.Context, secretKey string) (string, error)
}

type awsSecretsManagerClient struct {
	smClient *secretsmanager.Client
}

func newSecretManagerClient(cfg aws.Config) *awsSecretsManagerClient {
	//	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	//	if err != nil {
	//		fmt.Printf("error loading aws sdk config, err='%v'\n", err)
	//	}
	smc := secretsmanager.NewFromConfig(cfg)
	return &awsSecretsManagerClient{smClient: smc}
}

func (c *awsSecretsManagerClient) GetSecretValue(ctx context.Context, secretKey string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretKey),
	}

	result, err := c.smClient.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

type testRunnerClient struct{}

func newTestRunnerClient() *testRunnerClient {
	return &testRunnerClient{}
}

func (c *testRunnerClient) GetSecretValue(_ context.Context, key string) (string, error) {
	return key, nil
}
