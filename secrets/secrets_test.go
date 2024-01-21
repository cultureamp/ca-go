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
	client := NewAWSSecrets("us-west-2")
	assert.NotNil(t, client)
}

func TestGetSecretSuccess(t *testing.T) {
	expectedOutput := &secretsmanager.GetOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedAwsClient := new(mockedSMClient)
	mockedAwsClient.On("Get", mock.Anything).Return(expectedOutput, nil)

	secrets := NewAWSSecrets("us-west-2")
	secrets.client = mockedAwsClient

	result, err := secrets.Get("my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
}

func TestGetSecretOnError(t *testing.T) {
	expectedOutput := &secretsmanager.GetOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedAwsClient := new(mockedSMClient)
	mockedAwsClient.On("Get", mock.Anything).Return(expectedOutput, errors.New("test-error"))

	secrets := NewAWSSecrets("us-west-2")
	secrets.client = mockedAwsClient

	result, err := secrets.Get("my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
}

func TestGetSecretOnEmpty(t *testing.T) {
	expectedOutput := &secretsmanager.GetOutput{
		SecretString: nil,
	}
	mockedAwsClient := new(mockedSMClient)
	mockedAwsClient.On("Get", mock.Anything).Return(expectedOutput, nil)

	secrets := NewAWSSecrets("us-west-2")
	secrets.client = mockedAwsClient

	result, err := secrets.Get("my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
}

type mockedSMClient struct {
	mock.Mock
}

func (m *mockedSMClient) Get(input *secretsmanager.GetInput) (*secretsmanager.GetOutput, error) {
	args := m.Called(input)
	argZero, _ := args.Get(0).(*secretsmanager.GetOutput)
	argOne, _ := args.Get(1).(error)
	return argZero, argOne
}
