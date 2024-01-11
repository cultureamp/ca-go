package secrets_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cultureamp/ca-go/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func BasicExamples() {
	// this will automatically use the AWS Region as per the environment variable "AWS_REGION"
	answer, err := secrets.Get("my-test-secret")
	fmt.Printf("The answer to the secret is '%s' (err='%v')\n", answer, err)

	// or if you need secrets from another region other than the one you are running in use
	sm := secrets.NewAWSSecretsManager("a-different-region")
	answer, err = sm.Get("my-test-secret2")
	fmt.Printf("The answer to the secret2 is '%s' (err='%v')\n", answer, err)

	// or if you want to be able to mock the behavior
	// create a mock that supports the AWSSecretsManagerClient interface
	// mockedClient := new(mockedAWSSecretsManagerClient)
	// mockedClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)
	// secrets := NewAWSSecretsManagerWithClient(mockedClient)
}

func TestExampleMockPackageLevelMethods(t *testing.T) {
	// Example if you want to be able to mock package level calls

	// 1. set up your mock
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)

	// 2. override the package level DefaultAWSSecrets.Client with your mock
	secrets.DefaultAWSSecrets.Client = mockedClient

	// 3. call the package methods which will call you mock
	result, err := secrets.Get("my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
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
