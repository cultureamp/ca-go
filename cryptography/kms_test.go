package cryptography

import (
	"context"
	b64 "encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
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
		crypto := NewCryptography(testRegion, testKeyID)
		crypto.client = mockKmsClient

		expectedOutput := "SGVsbG8sIHBsYXlncm91bmQ="
		blob, err := b64.StdEncoding.DecodeString(expectedOutput)
		require.NoError(t, err)

		awsOutput := kms.EncryptOutput{
			CiphertextBlob: blob,
		}
		mockKmsClient.On("Encrypt", ctx, mock.Anything, mock.Anything).
			Return(&awsOutput, nil)

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
		crypto := NewCryptography(testRegion, testKeyID)
		crypto.client = mockKmsClient

		mockKmsClient.On("Encrypt", ctx, mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))

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
		crypto := NewCryptography(testRegion, testKeyID)
		crypto.client = mockKmsClient

		expectedOutput := "Decrypted"
		awsOutput := kms.DecryptOutput{
			Plaintext: []byte(expectedOutput),
		}
		mockKmsClient.On("Decrypt", ctx, mock.Anything, mock.Anything).
			Return(&awsOutput, nil)

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
		crypto := NewCryptography(testRegion, testKeyID)
		crypto.client = mockKmsClient

		mockKmsClient.On("Decrypt", ctx, mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))

		// act
		output, err := crypto.Decrypt(ctx, strInput)

		// assert
		require.Error(t, err)
		assert.Equal(t, "", output)
		mockKmsClient.AssertExpectations(t)
	})

	t.Run("When non-base64 input should error", func(t *testing.T) {
		// arrange
		mockKmsClient := &mockKMSClient{}
		crypto := NewCryptography(testRegion, testKeyID)
		crypto.client = mockKmsClient

		// act
		output, err := crypto.Decrypt(ctx, "59216167-f9c0-4b1b-b1db-1babd1209f10@ABC")

		// assert
		require.Error(t, err)
		assert.Equal(t, "", output)
		mockKmsClient.AssertExpectations(t)
	})
}

type mockKMSClient struct {
	mock.Mock
}

func (_m *mockKMSClient) Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error) {
	args := _m.Called(ctx, params, optFns)
	output, ok := args.Get(0).(*kms.EncryptOutput)
	if ok {
		return output, nil
	} else {
		return nil, args.Error(1)
	}
}

func (_m *mockKMSClient) Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	args := _m.Called(ctx, params, optFns)
	output, ok := args.Get(0).(*kms.DecryptOutput)
	if ok {
		return output, nil
	} else {
		return nil, args.Error(1)
	}
}
