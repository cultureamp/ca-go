package kms

import (
	"context"
	b64 "encoding/base64"
	"errors"
	"testing"

	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKMSClient struct {
	mock.Mock
}

func (_m *MockKMSClient) Encrypt(ctx context.Context, params *awskms.EncryptInput, optFns ...func(*awskms.Options)) (*awskms.EncryptOutput, error) {
	args := _m.Called(ctx, params, optFns)
	output, ok := args.Get(0).(*awskms.EncryptOutput)
	if ok {
		return output, nil
	} else {
		return nil, args.Error(1)
	}
}

func (_m *MockKMSClient) Decrypt(ctx context.Context, params *awskms.DecryptInput, optFns ...func(*awskms.Options)) (*awskms.DecryptOutput, error) {
	args := _m.Called(ctx, params, optFns)
	output, ok := args.Get(0).(*awskms.DecryptOutput)
	if ok {
		return output, nil
	} else {
		return nil, args.Error(1)
	}
}

const testID = "inputStr"

func TestEncrypt(t *testing.T) {
	strInput := "inputStr"
	ctx := context.Background()

	t.Run("Green Path - Should encrypt", func(t *testing.T) {
		// arrange
		mockKmsClient := newKMSClient(t)
		kms := NewKMSWithClient(testID, mockKmsClient)
		expectedOutput := "SGVsbG8sIHBsYXlncm91bmQ="
		blob, err := b64.StdEncoding.DecodeString(expectedOutput)
		assert.NoError(t, err)
		awsOutput := awskms.EncryptOutput{
			CiphertextBlob: blob,
		}
		mockKmsClient.On("Encrypt", ctx, mock.Anything, mock.Anything).
			Return(&awsOutput, nil)
		// act
		output, err := kms.Encrypt(ctx, strInput)
		// assert
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, expectedOutput, *output)
	})

	t.Run("When unable to use AWS client should error", func(t *testing.T) {
		// arrange
		kms := NewKMS(testID)
		// act
		output, err := kms.Encrypt(ctx, strInput)
		// assert
		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("When unable to encrypt should error", func(t *testing.T) {
		// arrange
		mockKmsClient := newKMSClient(t)
		kms := NewKMSWithClient(testID, mockKmsClient)
		mockKmsClient.On("Encrypt", ctx, mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))
		// act
		output, err := kms.Encrypt(ctx, strInput)
		// assert
		assert.Error(t, err)
		assert.Nil(t, output)
	})
}

func TestDecrypt(t *testing.T) {
	strInput := "inputStr"
	ctx := context.Background()

	t.Run("Green Path - Should decrypt", func(t *testing.T) {
		// arrange
		mockKmsClient := newKMSClient(t)
		kms := NewKMSWithClient(testID, mockKmsClient)
		expectedOutput := "Decrypted"
		awsOutput := awskms.DecryptOutput{
			Plaintext: []byte(expectedOutput),
		}
		mockKmsClient.On("Decrypt", ctx, mock.Anything, mock.Anything).
			Return(&awsOutput, nil)
		// act
		output, err := kms.Decrypt(ctx, strInput)
		// assert
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, expectedOutput, *output)
	})

	t.Run("When unable to use AWS client should error", func(t *testing.T) {
		// arrange
		kms := NewKMS(testID)
		// act
		output, err := kms.Decrypt(ctx, strInput)
		// assert
		assert.Error(t, err)
		assert.Nil(t, output)
	})

	t.Run("When unable to decrypt should error", func(t *testing.T) {
		// arrange
		mockKmsClient := newKMSClient(t)
		kms := NewKMSWithClient(testID, mockKmsClient)
		mockKmsClient.On("Decrypt", ctx, mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))
		// act
		output, err := kms.Decrypt(ctx, strInput)
		// assert
		assert.Error(t, err)
		assert.Nil(t, output)
	})
}

func newKMSClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKMSClient {
	mock := &MockKMSClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
