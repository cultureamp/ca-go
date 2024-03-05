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

// NewKMSCipher creates a new kms cipher for the specific "region" and "keyid".
func NewKMSCipher(region string) *kmsCipher {
	client := newAWSKMSClient(region)
	return NewKMSCipherWithClient(client)
}

// NewKMSCipherWithClient creates a new kms cipher for the specific "region" and "keyid".
func NewKMSCipherWithClient(client KMSCipher) *kmsCipher {
	return &kmsCipher{client}
}

// Encrypt will encrypt the "plainStr" using the region and keyID of the cipher.
func (c *kmsCipher) Encrypt(ctx context.Context, keyID string, plainStr string) (string, error) {
	return c.Client.Encrypt(ctx, keyID, plainStr)
}

// Decrypt will decrypt the "encryptedStr" using the region and keyID of the cipher.
func (c *kmsCipher) Decrypt(ctx context.Context, keyID string, encryptedStr string) (string, error) {
	return c.Client.Decrypt(ctx, keyID, encryptedStr)
}
