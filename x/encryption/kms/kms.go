package kms

import (
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"

	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/cultureamp/ca-go/x/encryption"
)

type KMSClient interface {
	Encrypt(ctx context.Context, params *awskms.EncryptInput, optFns ...func(*awskms.Options)) (*awskms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *awskms.DecryptInput, optFns ...func(*awskms.Options)) (*awskms.DecryptOutput, error)
}

type Encryptor struct {
	client KMSClient
	keyID  string
}

func NewEncryptor(keyID string, client KMSClient) (encryption.Encryptor, error) {
	if client == nil {
		return nil, errors.New("failed to get kms client")
	}
	return &Encryptor{client, keyID}, nil
}

func (e *Encryptor) Encrypt(ctx context.Context, plainStr string) (string, error) {
	input := &awskms.EncryptInput{
		KeyId:     &e.keyID,
		Plaintext: []byte(plainStr),
	}

	result, err := e.client.Encrypt(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt with kms: %w", err)
	}

	blobString := b64.StdEncoding.EncodeToString(result.CiphertextBlob)
	return blobString, nil
}

func (e *Encryptor) Decrypt(ctx context.Context, encryptedStr string) (string, error) {
	blob, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode: %w", err)
	}

	input := &awskms.DecryptInput{
		CiphertextBlob: blob,
	}

	result, err := e.client.Decrypt(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt with kms: %w", err)
	}

	decStr := string(result.Plaintext)

	return decStr, nil
}
