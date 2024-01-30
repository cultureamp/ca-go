package secrets

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewAWSSecretsClient(t *testing.T) {
	secrets, err := NewAWSSecretsManager("us-west-2")
	assert.Nil(t, err)
	assert.NotNil(t, secrets)
}

func TestGetSecretSuccess(t *testing.T) {
	ctx := context.Background()
	expectedOutput := "my-super-secret-value"
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything, mock.Anything).Return(expectedOutput, nil)

	secrets := NewAWSSecretsManagerWithClient(mockedClient)
	result, err := secrets.Get(ctx, "my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
	mockedClient.AssertExpectations(t)
}

func TestGetSecretOnError(t *testing.T) {
	ctx := context.Background()
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything, mock.Anything).Return("", errors.New("test-error"))

	secrets := NewAWSSecretsManagerWithClient(mockedClient)
	result, err := secrets.Get(ctx, "my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	mockedClient.AssertExpectations(t)
}

func TestGetSecretOnEmpty(t *testing.T) {
	ctx := context.Background()
	mockedClient := new(mockedAWSSecretsManagerClient)
	mockedClient.On("GetSecretValue", mock.Anything, mock.Anything).Return("", nil)

	secrets := NewAWSSecretsManagerWithClient(mockedClient)
	result, err := secrets.Get(ctx, "my-secret")
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	mockedClient.AssertExpectations(t)
}

type mockedAWSSecretsManagerClient struct {
	mock.Mock
}

func (m *mockedAWSSecretsManagerClient) GetSecretValue(ctx context.Context, secretKey string) (string, error) {
	args := m.Called(ctx, secretKey)
	argZero, _ := args.Get(0).(string)
	argOne, _ := args.Get(1).(error)
	return argZero, argOne
}
