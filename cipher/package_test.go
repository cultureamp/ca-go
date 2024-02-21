package cipher_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/cipher"
	"github.com/stretchr/testify/assert"
)

func TestPackageEncrypt(t *testing.T) {
	ctx := context.Background()
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"

	// replace the package level client with our mock
	stdClient := cipher.DefaultKMSCipher.Client
	cipher.DefaultKMSCipher.Client = newTestRunnerClient()
	defer func() {
		cipher.DefaultKMSCipher.Client = stdClient
	}()

	cipherText, err := cipher.Encrypt(ctx, keyId, "test_plain_str")
	assert.Nil(t, err)

	plainText, err := cipher.Decrypt(ctx, keyId, cipherText)
	assert.Nil(t, err)
	assert.Equal(t, "test_plain_str", plainText)
}

type testRunnerClient struct{}

func newTestRunnerClient() *testRunnerClient {
	return &testRunnerClient{}
}

// Encrypt on the test runner just returns the "plainStr" as the encrypted encryptedStr.
func (c *testRunnerClient) Encrypt(ctx context.Context, _ string, plainStr string) (string, error) {
	return plainStr, nil
}

// Decrypt on the test runner just returns the "encryptedStr" as the decrypted plainstr.
func (c *testRunnerClient) Decrypt(ctx context.Context, _ string, encryptedStr string) (string, error) {
	return encryptedStr, nil
}
