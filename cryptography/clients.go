package cryptography

import (
	"context"
	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/pkg/errors"
)

// KMSClient used for mock testing.
type KMSClient interface {
	Encrypt(ctx context.Context, keyId string, plainStr string) (string, error)
	Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error)
}

type awsKMSClient struct {
	kmsClient *kms.Client
}

func newAWSKMSClient(region string) *awsKMSClient {
	client := kms.New(kms.Options{Region: region})
	return &awsKMSClient{
		kmsClient: client,
	}
}

func (c *awsKMSClient) Encrypt(ctx context.Context, keyId string, plainStr string) (string, error) {
	input := &kms.EncryptInput{
		KeyId:     &keyId,
		Plaintext: []byte(plainStr),
	}

	result, err := c.kmsClient.Encrypt(ctx, input)
	if err != nil {
		return "", err
	}

	blobString := b64.StdEncoding.EncodeToString(result.CiphertextBlob)
	return blobString, nil
}

func (c *awsKMSClient) Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error) {
	blob, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode")
	}

	input := &kms.DecryptInput{
		KeyId:          &keyId,
		CiphertextBlob: blob,
	}

	result, err := c.kmsClient.Decrypt(ctx, input)
	if err != nil {
		return "", err
	}

	decStr := string(result.Plaintext)
	return decStr, nil
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
