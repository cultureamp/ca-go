package secrets

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAWSSecretsClient(t *testing.T) {
	client, err := NewAWSSecretsClient("us-west-2")

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
}

type MockClient struct {
	get func(input *secretsmanager.GetInput) (*secretsmanager.GetOutput, error)
}

func (m MockClient) Get(input *secretsmanager.GetInput) (*secretsmanager.GetOutput, error) {
	return m.get(input)
}

func TestAWSSecretsGet(t *testing.T) {
	var actualRequestedSecretName string
	secretManager := AWSSecrets{
		Client: MockClient{
			get: func(input *secretsmanager.GetInput) (*secretsmanager.GetOutput, error) {
				actualRequestedSecretName = *input.SecretId

				return &secretsmanager.GetOutput{
					SecretString: aws.String("my-super-secret-value"),
				}, nil
			},
		},
	}

	secretString, err := secretManager.Get("my-secret-name")

	require.NoError(t, err)
	assert.Equal(t, "my-super-secret-value", secretString)
	assert.Equal(t, "my-secret-name", actualRequestedSecretName)
}

func TestAWSSecretsGetEmptySecret(t *testing.T) {
	secretManager := AWSSecrets{
		Client: MockClient{
			get: func(input *secretsmanager.GetInput) (*secretsmanager.GetOutput, error) {
				return &secretsmanager.GetOutput{
					SecretString: nil,
				}, nil
			},
		},
	}

	secretString, err := secretManager.Get("my-secret-name")

	require.Error(t, err)
	assert.Equal(t, "", secretString)
}
