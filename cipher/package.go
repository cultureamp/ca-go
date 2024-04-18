package cipher

import (
	"context"
	"os"
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
	mustHaveDefaultKMSCipher()

	return DefaultKMSCipher.Encrypt(ctx, keyId, plainStr)
}

// Decrpyt will use env var AWS_REGION and the KMS keyId and the base64 encoded encryptedStr and return it decrypted as a plain string.
func Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error) {
	mustHaveDefaultKMSCipher()

	return DefaultKMSCipher.Decrypt(ctx, keyId, encryptedStr)
}

func mustHaveDefaultKMSCipher() {
	if DefaultKMSCipher != nil {
		return // its set so we are good to go
	}

	region := os.Getenv("AWS_REGION")
	kmsClient := NewKMSClient(region)

	DefaultKMSCipher = kmsClient
}
