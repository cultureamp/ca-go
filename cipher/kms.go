package cipher

import (
	"context"
	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/go-errors/errors"
)

type awsKMSClient struct {
	aws *kms.Client
}

func NewKMSClient(region string, optFns ...func(*kms.Options)) *awsKMSClient {
	client := kms.New(kms.Options{Region: region}, optFns...)
	return &awsKMSClient{
		aws: client,
	}
}

// Encrypt will use the KMS keyId to encrypt the plainStr and return it as a base64 encoded string.
func (c *awsKMSClient) Encrypt(ctx context.Context, keyId string, plainStr string) (string, error) {
	input := &kms.EncryptInput{
		KeyId:     &keyId,
		Plaintext: []byte(plainStr),
	}

	result, err := c.aws.Encrypt(ctx, input)
	if err != nil {
		return "", err
	}

	blobString := b64.StdEncoding.EncodeToString(result.CiphertextBlob)
	return blobString, nil
}

// Decrpyt will use the KMS keyId and the base64 encoded encryptedStr and return it decrypted as a plain string.
func (c *awsKMSClient) Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error) {
	blob, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", errors.Errorf("failed to decode: %w", err)
	}

	input := &kms.DecryptInput{
		KeyId:          &keyId,
		CiphertextBlob: blob,
	}

	result, err := c.aws.Decrypt(ctx, input)
	if err != nil {
		return "", err
	}

	decStr := string(result.Plaintext)
	return decStr, nil
}
