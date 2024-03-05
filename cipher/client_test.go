package cipher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientEncrypt(t *testing.T) {
	ctx := context.Background()
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
	region := "region"

	t.Run("Success: With Mocked HTTP Client", func(t *testing.T) {
		bodyReader := io.NopCloser(strings.NewReader(`{
			"CiphertextBlob": "SGVsbG8sIHBsYXlncm91bmQ=",
			"EncryptionAlgorithm": "SYMMETRIC_DEFAULT",
			"KeyId": "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
		}`))

		expectedResp := &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
			Body:       bodyReader,
		}
		mockHTTPClient := mockKMSClient{}
		mockHTTPClient.On("Do", mock.Anything).Return(expectedResp, nil)

		client := newAWSKMSClient(region, func(opt *kms.Options) { opt.HTTPClient = &mockHTTPClient })

		cipherText, err := client.Encrypt(ctx, keyId, "Hello, playground")
		assert.Nil(t, err)
		assert.Equal(t, "SGVsbG8sIHBsYXlncm91bmQ=", cipherText)
	})

	t.Run("Error: Mocked HTTP Client returns error", func(t *testing.T) {
		mockHTTPClient := mockKMSClient{}
		mockHTTPClient.On("Do", mock.Anything).Return(nil, fmt.Errorf("expected error"))

		client := newAWSKMSClient(region, func(opt *kms.Options) {
			opt.HTTPClient = &mockHTTPClient
			opt.Retryer = &aws.NopRetryer{}
		})

		cipherText, err := client.Encrypt(ctx, keyId, "Hello, playground")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "expected error")
		assert.Equal(t, "", cipherText)
	})
}

func TestClientDecrypt(t *testing.T) {
	ctx := context.Background()
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
	region := "region"

	t.Run("Success: With Mocked HTTP Client", func(t *testing.T) {
		bodyReader := io.NopCloser(strings.NewReader(`{
			"CiphertextForRecipient": "SGVsbG8sIHBsYXlncm91bmQ=",
			"EncryptionAlgorithm": "SYMMETRIC_DEFAULT",
			"KeyId": "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab",
			"Plaintext": "SGVsbG8sIHBsYXlncm91bmQ="
		}`))

		expectedResp := &http.Response{
			StatusCode: 200,
			Header:     http.Header{},
			Body:       bodyReader,
		}
		mockHTTPClient := mockKMSClient{}
		mockHTTPClient.On("Do", mock.Anything).Return(expectedResp, nil)

		client := newAWSKMSClient(region, func(opt *kms.Options) { opt.HTTPClient = &mockHTTPClient })

		plainText, err := client.Decrypt(ctx, keyId, "SGVsbG8sIHBsYXlncm91bmQ=")
		assert.Nil(t, err)
		assert.Equal(t, "Hello, playground", plainText)
	})

	t.Run("Error: Mocked HTTP Client returns error", func(t *testing.T) {
		mockHTTPClient := mockKMSClient{}
		mockHTTPClient.On("Do", mock.Anything).Return(nil, fmt.Errorf("expected error"))

		client := newAWSKMSClient(region, func(opt *kms.Options) {
			opt.HTTPClient = &mockHTTPClient
			opt.Retryer = &aws.NopRetryer{}
		})

		plainText, err := client.Decrypt(ctx, keyId, "SGVsbG8sIHBsYXlncm91bmQ=")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "expected error")
		assert.Equal(t, "", plainText)
	})

	t.Run("Error: Bad Base64 CipherText", func(t *testing.T) {
		mockHTTPClient := mockKMSClient{}
		mockHTTPClient.On("Do", mock.Anything).Return(nil, nil)

		client := newAWSKMSClient(region, func(opt *kms.Options) {
			opt.HTTPClient = &mockHTTPClient
			opt.Retryer = &aws.NopRetryer{}
		})

		plainText, err := client.Decrypt(ctx, keyId, "%^&*()")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "illegal base64 data")
		assert.Equal(t, "", plainText)
	})
}

type mockHTTPClient struct {
	mock.Mock
}

func (_m *mockKMSClient) Do(req *http.Request) (*http.Response, error) {
	args := _m.Called(req)
	output, _ := args.Get(0).(*http.Response)
	return output, args.Error(1)
}
