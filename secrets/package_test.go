package secrets

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockPackageLevelMethods(t *testing.T) {
	// Example if you want to be able to mock package level calls

	// 1. set up your mock
	expectedOutput := &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String("my-super-secret-value"),
	}
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)

	// 2. override the package level DefaultAWSSecrets.Client with your mock
	DefaultAWSSecrets.Client = mockedClient

	// 3. call the package methods which will call you mock
	result, err := Get("my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
	mockedClient.AssertExpectations(t)
}
