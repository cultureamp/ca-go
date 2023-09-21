package encryption

import (
	"context"
	"errors"
	"testing"

	"github.com/cultureamp/ca-go/ref"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type KMSEncrypt struct {
	mock.Mock
}

// Decrypt provides a mock function with given fields: ctx, encryptedStr
func (_m *KMSEncrypt) Decrypt(ctx context.Context, encryptedStr string) (*string, error) {
	args := _m.Called(ctx, encryptedStr)
	output, ok := args.Get(0).(*string)
	if ok {
		return output, nil
	} else {
		return nil, args.Error(1)
	}
}

// Encrypt provides a mock function with given fields: ctx, plainStr
func (_m *KMSEncrypt) Encrypt(ctx context.Context, plainStr string) (*string, error) {
	args := _m.Called(ctx, plainStr)
	output, ok := args.Get(0).(*string)
	if ok {
		return output, nil
	} else {
		return nil, args.Error(1)
	}
}

func TestDecrypt(t *testing.T) {
	t.Run("should return an error on failing to decrypt", func(t *testing.T) {
		ctx := context.Background()
		encryptedStr := "@#trtrtrtrt!"
		mockKMS := newKMSEncrypt(t)
		mockKMS.On("Decrypt", ctx, encryptedStr).
			Return(nil, errors.New("decrypt failed"))

		encryptionSvc := NewEncryption(mockKMS)

		_, err := encryptionSvc.Decrypt(ctx, encryptedStr)

		assert.Error(t, err)
		assert.Equal(t, "failed to decrypt the string: decrypt failed", err.Error())
	})

	t.Run("should return decrypted error on success", func(t *testing.T) {
		ctx := context.Background()
		encryptedStr := "@#trtrtrtrt!"
		decryptedStr := "test1"
		mockKMS := newKMSEncrypt(t)
		mockKMS.On("Decrypt", ctx, encryptedStr).
			Return(ref.String(decryptedStr), nil)

		encryptionSvc := NewEncryption(mockKMS)

		ds, err := encryptionSvc.Decrypt(ctx, encryptedStr)

		assert.NotNil(t, ds)
		assert.Equal(t, decryptedStr, *ds)
		assert.Nil(t, err)
	})
}

func TestEncrypt(t *testing.T) {
	t.Run("should return an error when failed to marshall the data", func(t *testing.T) {
		ctx := context.Background()
		// setting a channel to fail the json marshalling
		data := make(chan int)

		encryptionSvc := NewEncryption(nil)

		_, err := encryptionSvc.Encrypt(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, "failed to marshal data: json: unsupported type: chan int", err.Error())
	})

	t.Run("should return an error when failing to encrypt", func(t *testing.T) {
		ctx := context.Background()
		data := "test2"
		mockKMS := newKMSEncrypt(t)
		mockKMS.On("Encrypt", ctx, mock.AnythingOfType("string")).
			Return(nil, errors.New("encrypt failed"))

		encryptionSvc := NewEncryption(mockKMS)

		_, err := encryptionSvc.Encrypt(ctx, data)

		assert.Error(t, err)
		assert.Equal(t, "failed to encrypt data: encrypt failed", err.Error())
	})

	t.Run("should return encrypted string on success", func(t *testing.T) {
		ctx := context.Background()
		data := "test3"
		encryptedStr := ref.String("123@#$$!")
		mockKMS := newKMSEncrypt(t)
		mockKMS.On("Encrypt", ctx, mock.AnythingOfType("string")).
			Return(encryptedStr, nil)

		encryptionSvc := NewEncryption(mockKMS)

		es, err := encryptionSvc.Encrypt(ctx, data)

		assert.NotNil(t, es)
		assert.Equal(t, encryptedStr, es)
		assert.Nil(t, err)
	})
}

func newKMSEncrypt(t interface {
	mock.TestingT
	Cleanup(func())
}) *KMSEncrypt {
	mock := &KMSEncrypt{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
