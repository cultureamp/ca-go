package cryptography

import (
	"context"
	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/pkg/errors"
)


// KMSCryptography supports basic Encrypt & Decrypt methods.
type KMSCryptography struct {
	client KMSClient
	keyID  string
}

// NewKMSCryptography creates a new kms cryptography for the specific "region" and "keyid".
func NewKMSCryptography(region string, keyID string) *KMSCryptography {
	client := kms.New(kms.Options{Region: region})
	return &KMSCryptography{client, keyID}
}

// NewKMSCryptography creates a new kms cryptography for the specific "region" and "keyid".
func NewKMSCryptographyWithClient(client KMSClient, keyID string) *KMSCryptography {
	return &KMSCryptography{client, keyID}
}

// Encrypt will encrypt the "plainStr" using the region and keyID of the cryptography.
func (c *KMSCryptography) Encrypt(ctx context.Context, plainStr string) (string, error) {
	input := &kms.EncryptInput{
		KeyId:     &c.keyID,
		Plaintext: []byte(plainStr),
	}

	result, err := c.client.Encrypt(ctx, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt with kms")
	}

	blobString := b64.StdEncoding.EncodeToString(result.CiphertextBlob)
	return blobString, nil
}

// Decrypt will decrypt the "encryptedStr" using the region and keyID of the cryptography.
func (c *KMSCryptography) Decrypt(ctx context.Context, encryptedStr string) (string, error) {
	blob, err := b64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode")
	}

	input := &kms.DecryptInput{
		CiphertextBlob: blob,
	}

	result, err := c.client.Decrypt(ctx, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt with kms")
	}

	decStr := string(result.Plaintext)
	return decStr, nil
}
