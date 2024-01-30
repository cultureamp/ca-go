package secrets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPackageLevelMethods(t *testing.T) {
	// Example when running in a test, the default client returns key as secret (no AWS call)
	ctx := context.Background()

	// 1. call the package methods which will call you mock
	result, err := Get(ctx, "my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-secret", result)
}

func TestMockPackageLevelMethods(t *testing.T) {
	// Example if you want to be able to mock package level calls
	ctx := context.Background()

	// 1. set up your mock
	expectedOutput := "my-super-secret-value"
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything, mock.Anything).Return(expectedOutput, nil)

	// 2. override the package level DefaultAWSSecrets.Client with your mock
	oldClient := DefaultAWSSecrets.Client
	defer func() { DefaultAWSSecrets.Client = oldClient }()
	DefaultAWSSecrets.Client = mockedClient

	// 3. call the package methods which will call you mock
	result, err := Get(ctx, "my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
	mockedClient.AssertExpectations(t)
}
