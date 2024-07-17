package cipher

import (
	"context"
	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/go-errors/errors"
)

type KMSClient struct {
	aws *kms.Client
}

func NewKMSClient(region string, optFns ...func(*kms.Options)) *KMSClient {
	client := kms.New(kms.Options{Region: region}, optFns...)
	return &KMSClient{
		aws: client,
	}
}

// Encrypt will use the KMS keyId to encrypt the plainStr and return it as a base64 encoded string.
func (c *KMSClient) Encrypt(ctx context.Context, keyID string, plainStr string) (string, error) {
	input := &kms.EncryptInput{
		KeyId:     &keyID,
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
func (c *KMSClient) Decrypt(ctx context.Context, keyID string, encryptedStr string) (string, error) {
	blob, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", errors.Errorf("failed to decode: %w", err)
	}

	input := &kms.DecryptInput{
		KeyId:          &keyID,
		CiphertextBlob: blob,
	}

	result, err := c.aws.Decrypt(ctx, input)
	if err != nil {
		return "", err
	}

	decStr := string(result.Plaintext)
	return decStr, nil
}
