package cipher

import (
	"context"
	"os"

	"github.com/go-errors/errors"
)

// KMSCipher used for mock testing.
type KMSCipher interface {
	Encrypt(ctx context.Context, keyId string, plainStr string) (string, error)
	Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error)
}

// DefaultKMSCipher is the package level default implementation used by all package level methods.
// Package level methods are provided for ease of use.
// For testing you can replace the DefaultKMSCipher client with your own mock:
//
//	DefaultKMSCipher = newMockedClient()
var DefaultKMSCipher KMSCipher = nil

// Encrypt will use env var AWS_REGION and the KMS keyId to encrypt the plainStr and return it as a base64 encoded string.
func Encrypt(ctx context.Context, keyId string, plainStr string) (string, error) {
	err := mustHaveDefaultKMSCipher()
	if err != nil {
		return "", err
	}

	return DefaultKMSCipher.Encrypt(ctx, keyId, plainStr)
}

// Decrpyt will use env var AWS_REGION and the KMS keyId and the base64 encoded encryptedStr and return it decrypted as a plain string.
func Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error) {
	err := mustHaveDefaultKMSCipher()
	if err != nil {
		return "", err
	}

	return DefaultKMSCipher.Decrypt(ctx, keyId, encryptedStr)
}

func mustHaveDefaultKMSCipher() error {
	if DefaultKMSCipher != nil {
		return nil // its set so we are good to go
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		return errors.Errorf("missing value for environment variable 'AWS_REGION'")
	}

	DefaultKMSCipher = NewKMSClient(region)
	return nil
}
