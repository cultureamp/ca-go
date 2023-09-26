package kms

import (
	"context"
	b64 "encoding/base64"

	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/cultureamp/ca-go/x/encryption"
	"github.com/pkg/errors"
)

type KMSClient interface {
	Encrypt(ctx context.Context, params *awskms.EncryptInput, optFns ...func(*awskms.Options)) (*awskms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *awskms.DecryptInput, optFns ...func(*awskms.Options)) (*awskms.DecryptOutput, error)
}

type Encryptor struct {
	client KMSClient
	keyID  *string
}

func NewEncryptor(keyID string, client KMSClient) (encryption.Encryptor, error) {
	if client == nil {
		return nil, errors.New("failed to get kms client")
	}
	return &Encryptor{client, &keyID}, nil
}

func (k *Encryptor) Encrypt(ctx context.Context, plainStr string) (*string, error) {
	input := &awskms.EncryptInput{
		KeyId:     k.keyID,
		Plaintext: []byte(plainStr),
	}

	result, err := k.client.Encrypt(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt with kms")
	}

	blobString := b64.StdEncoding.EncodeToString(result.CiphertextBlob)
	return &blobString, nil
}

func (k *Encryptor) Decrypt(ctx context.Context, encryptedStr string) (*string, error) {
	blob, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	input := &awskms.DecryptInput{
		CiphertextBlob: blob,
	}

	result, err := k.client.Decrypt(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt with kms")
	}

	decStr := string(result.Plaintext)

	return &decStr, nil
}
