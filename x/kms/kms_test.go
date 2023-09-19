package kms

import (
	"context"
	b64 "encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKMSClient struct {
	mock.Mock
}

func (_m *MockKMSClient) Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *kms.EncryptOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *kms.EncryptInput, ...func(*kms.Options)) (*kms.EncryptOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *kms.EncryptInput, ...func(*kms.Options)) *kms.EncryptOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.EncryptOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *kms.EncryptInput, ...func(*kms.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockKMSClient) Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *kms.DecryptOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *kms.DecryptInput, ...func(*kms.Options)) (*kms.DecryptOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *kms.DecryptInput, ...func(*kms.Options)) *kms.DecryptOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kms.DecryptOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *kms.DecryptInput, ...func(*kms.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

const testID = "inputStr"

func TestEncrypt(t *testing.T) {
	strInput := "inputStr"
	ctx := context.Background()

	t.Run("Green Path - Should encrypt", func(t *testing.T) {
		// arrange
		mockKmsClient := newKMSClient(t)
		kms := NewKMSWithClient(testID, &mockKmsClient)
		expectedOutput := "SGVsbG8sIHBsYXlncm91bmQ="
		blob, err := b64.StdEncoding.DecodeString(expectedOutput)
		assert.NoError(t, err)
		awsOutput := awskms.EncryptOutput{
			CiphertextBlob: blob,
		}
		mockKmsClient.On("Encrypt", ctx, mock.Anything).
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
		kms := NewKMSWithClient(testID, &mockKmsClient)
		mockKmsClient.On("Encrypt", ctx, mock.Anything).
			Return(nil, errors.New("error"))
		// act
		output, err := kms.Encrypt(ctx, strInput)
		// assert
		assert.Error(t, err)
		assert.Nil(t, output)
	})

}

// Decrypt(ctx context.Context, encryptedStr string) (decryptedStr *string, err error)
func TestDecrypt(t *testing.T) {

	strInput := "inputStr"
	ctx := context.Background()

	t.Run("Green Path - Should decrypt", func(t *testing.T) {
		// arrange
		mockKmsClient := newKMSClient(t)
		kms := NewKMSWithClient(testID, &mockKmsClient)
		expectedOutput := "Decrypted"
		awsOutput := awskms.DecryptOutput{
			Plaintext: []byte(expectedOutput),
		}
		mockKmsClient.On("Decrypt", ctx, mock.Anything).
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
		kms := NewKMSWithClient(testID, &mockKmsClient)
		mockKmsClient.On("Decrypt", ctx, mock.Anything).
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
}) MockKMSClient {
	mock := MockKMSClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
