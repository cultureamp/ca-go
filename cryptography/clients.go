package cryptography

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// KMSClient used for mock testing.
type KMSClient interface {
	Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

type testRunnerClient struct{}

func newTestRunnerClient() *testRunnerClient {
	return &testRunnerClient{}
}

// Encrypt on the test runner just returns the "plainStr" as the encrypted CiphertextBlob.
func (c *testRunnerClient) Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error) {
	retval := &kms.EncryptOutput{
		CiphertextBlob: params.Plaintext,
		KeyId:          params.KeyId,
	}

	return retval, nil
}

// Decrypt on the test runner just returns the "CiphertextBlob" as the decrypted plainstr.
func (c *testRunnerClient) Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	retval := &kms.DecryptOutput{
		KeyId:     params.KeyId,
		Plaintext: params.CiphertextBlob,
	}

	return retval, nil
}
