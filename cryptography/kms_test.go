package cryptography

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testKeyID  = "inputStr"
	testRegion = "us-west-1"
)

func TestEncrypt(t *testing.T) {
	strInput := "inputStr"
	ctx := context.Background()

	t.Run("Green Path - Should encrypt", func(t *testing.T) {
		// arrange
		mockKmsClient := &mockKMSClient{}
		crypto := NewKMSCryptographyWithClient(mockKmsClient)

		expectedOutput := "SGVsbG8sIHBsYXlncm91bmQ="
		mockKmsClient.On("Encrypt", ctx, mock.Anything, mock.Anything).Return(expectedOutput, nil)

		// act
		output, err := crypto.Encrypt(ctx, strInput)

		// assert
		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, expectedOutput, output)
		mockKmsClient.AssertExpectations(t)
	})

	t.Run("When unable to encrypt should error", func(t *testing.T) {
		// arrange
		mockKmsClient := &mockKMSClient{}
		crypto := NewKMSCryptographyWithClient(mockKmsClient)
		crypto.client = mockKmsClient

		mockKmsClient.On("Encrypt", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		// act
		output, err := crypto.Encrypt(ctx, strInput)

		// assert
		require.Error(t, err)
		assert.Equal(t, "", output)
		mockKmsClient.AssertExpectations(t)
	})
}

func TestDecrypt(t *testing.T) {
	strInput := "inputStr"
	ctx := context.Background()

	t.Run("Green Path - Should decrypt", func(t *testing.T) {
		// arrange
		mockKmsClient := &mockKMSClient{}
		crypto := NewKMSCryptographyWithClient(mockKmsClient)

		expectedOutput := "Decrypted"
		mockKmsClient.On("Decrypt", ctx, mock.Anything, mock.Anything).Return(expectedOutput, nil)

		// act
		output, err := crypto.Decrypt(ctx, strInput)

		// assert
		require.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, expectedOutput, output)
		mockKmsClient.AssertExpectations(t)
	})

	t.Run("When unable to decrypt should error", func(t *testing.T) {
		// arrange
		mockKmsClient := &mockKMSClient{}
		crypto := NewKMSCryptographyWithClient(mockKmsClient)
		crypto.client = mockKmsClient

		mockKmsClient.On("Decrypt", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		// act
		output, err := crypto.Decrypt(ctx, strInput)

		// assert
		require.Error(t, err)
		assert.Equal(t, "", output)
		mockKmsClient.AssertExpectations(t)
	})
}

type mockKMSClient struct {
	mock.Mock
}

func (_m *mockKMSClient) Encrypt(ctx context.Context, plainStr string) (string, error) {
	args := _m.Called(ctx, plainStr)
	output, _ := args.Get(0).(string)
	return output, args.Error(1)
}

func (_m *mockKMSClient) Decrypt(ctx context.Context, encryptedStr string) (string, error) {
	args := _m.Called(ctx, encryptedStr)
	output, _ := args.Get(0).(string)
	return output, args.Error(1)
}
