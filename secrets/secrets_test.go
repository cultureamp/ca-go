package secrets

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAWSSecretsClient(t *testing.T) {
	secrets := NewAWSSecretsManager("us-west-2")
	assert.NotNil(t, secrets)
}

func TestGetSecretSuccess(t *testing.T) {
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)

	secrets := NewAWSSecretsManagerWithClient(mockedClient)
	result, err := secrets.Get("my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
	mockedClient.AssertExpectations(t)
}

func TestGetSecretOnError(t *testing.T) {
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, errors.New("test-error"))

	secrets := NewAWSSecretsManagerWithClient(mockedClient)
	result, err := secrets.Get("my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	mockedClient.AssertExpectations(t)
}

func TestGetSecretOnEmpty(t *testing.T) {
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: nil,
	}
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)

	secrets := NewAWSSecretsManagerWithClient(mockedClient)
	result, err := secrets.Get("my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	mockedClient.AssertExpectations(t)
}

type mockedAWSSecretsManagerClient struct {
	mock.Mock
}

func (m *mockedAWSSecretsManagerClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	args := m.Called(input)
	argZero, _ := args.Get(0).(*secretsmanager.GetSecretValueOutput)
	argOne, _ := args.Get(1).(error)
	return argZero, argOne
}
