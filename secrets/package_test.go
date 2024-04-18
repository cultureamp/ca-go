package secrets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockPackageLevelMethods(t *testing.T) {
	// Example if you want to be able to mock package level calls
	ctx := context.Background()

	// 1. set up your mock
	expectedOutput := "my-super-secret-value"
	mockSM := new(mockedAWSSecretsManager)
	mockSM.On("Get", mock.Anything, mock.Anything).Return(expectedOutput, nil)

	// 2. override the package level DefaultAWSSecrets.Client with your mock
	oldSM := DefaultAWSSecretsManager
	defer func() { DefaultAWSSecretsManager = oldSM }()
	DefaultAWSSecretsManager = mockSM

	// 3. call the package methods which will call you mock
	result, err := Get(ctx, "my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "my-super-secret-value", result)
	mockSM.AssertExpectations(t)
}

type mockedAWSSecretsManager struct {
	mock.Mock
}

func (m *mockedAWSSecretsManager) Get(ctx context.Context, secretKey string) (string, error) {
	args := m.Called(ctx, secretKey)
	argZero, _ := args.Get(0).(string)
	argOne, _ := args.Get(1).(error)
	return argZero, argOne
}
