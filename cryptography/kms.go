package cryptography

import (
	"context"
	b64 "encoding/base64"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/pkg/errors"
)

// kmsClient used for mock testing.
type kmsClient interface {
	Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

// KMSCryptography supports basic Encrypt & Decrypt methods.
type KMSCryptography struct {
	client kmsClient
	keyID  string
}

var defaultKMSCrypto *KMSCryptography = getInstance()

func getInstance() *KMSCryptography {
	// Should this take dependency on 'env' package and call env.AwsRegion()?
	region := os.Getenv("AWS_REGION")
	keyID := os.Getenv("KMS_KEY_ID")
	return NewKMSCryptography(region, keyID)
}

// NewKMSCryptography creates a new kms cryptography for the specific "region" and "keyid".
func NewKMSCryptography(region string, keyID string) *KMSCryptography {
	client := kms.New(kms.Options{Region: region})
	return &KMSCryptography{client, keyID}
}

// Encrypt uses the default AWS_REGION and KMS_KEY_ID to kms encrypt "plainStr".
func Encrypt(ctx context.Context, plainStr string) (string, error) {
	return defaultKMSCrypto.Encrypt(ctx, plainStr)
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

// Decrypt uses the default AWS_REGION and KMS_KEY_ID to kms decrypt "encryptedStr".
func Decrypt(ctx context.Context, encryptedStr string) (string, error) {
	return defaultKMSCrypto.Decrypt(ctx, encryptedStr)
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
