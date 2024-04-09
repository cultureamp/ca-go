package cipher

import (
	"context"
)

// KMSCipher used for mock testing.
type KMSCipher interface {
	Encrypt(ctx context.Context, keyId string, plainStr string) (string, error)
	Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error)
}

// kmsCipher supports basic Encrypt & Decrypt methods.
type kmsCipher struct {
	Client KMSCipher
}

// NewKMSCipher creates a new kms cipher for the specific "region".
func NewKMSCipher(region string) *kmsCipher {
	client := newAWSKMSClient(region)
	return NewKMSCipherWithClient(client)
}

// NewKMSCipherWithClient creates a new kms cipher using a specific client that supports the KMSCipher interface.
// This is provided mostly for testing purposes.
func NewKMSCipherWithClient(client KMSCipher) *kmsCipher {
	return &kmsCipher{client}
}

// Encrypt will use the KMS keyId to encrypt the plainStr and return it as a base64 encoded string.
func (c *kmsCipher) Encrypt(ctx context.Context, keyID string, plainStr string) (string, error) {
	return c.Client.Encrypt(ctx, keyID, plainStr)
}

// Decrpyt will use the KMS keyId and the base64 encoded encryptedStr and return it decrypted as a plain string.
func (c *kmsCipher) Decrypt(ctx context.Context, keyID string, encryptedStr string) (string, error) {
	return c.Client.Decrypt(ctx, keyID, encryptedStr)
}
