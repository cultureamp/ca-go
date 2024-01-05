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
	client := newAWSSecrets("us-west-2")
	assert.NotNil(t, client)
}

func TestGetSecretSuccess(t *testing.T) {
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedAwsClient := new(mockedSMClient)
	mockedAwsClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)

	secrets := newAWSSecrets("us-west-2")
	secrets.client = mockedAwsClient

	result, err := secrets.Get("my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
}

func TestGetSecretOnError(t *testing.T) {
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedAwsClient := new(mockedSMClient)
	mockedAwsClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, errors.New("test-error"))

	secrets := newAWSSecrets("us-west-2")
	secrets.client = mockedAwsClient

	result, err := secrets.Get("my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
}

func TestGetSecretOnEmpty(t *testing.T) {
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: nil,
	}
	mockedAwsClient := new(mockedSMClient)
	mockedAwsClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)

	secrets := newAWSSecrets("us-west-2")
	secrets.client = mockedAwsClient

	result, err := secrets.Get("my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
}

type mockedSMClient struct {
	mock.Mock
}

func (m *mockedSMClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	args := m.Called(input)
	argZero, _ := args.Get(0).(*secretsmanager.GetSecretValueOutput)
	argOne, _ := args.Get(1).(error)
	return argZero, argOne
}
