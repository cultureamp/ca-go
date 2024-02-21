package cipher

import (
	"context"
)

// KMSCipher supports basic Encrypt & Decrypt methods.
type KMSCipher struct {
	Client KMSClient
}

// NewKMSCipher creates a new kms cipher for the specific "region" and "keyid".
func NewKMSCipher(region string) *KMSCipher {
	client := newAWSKMSClient(region)
	return &KMSCipher{client}
}

// NewKMSCipherWithClient creates a new kms cipher for the specific "region" and "keyid".
func NewKMSCipherWithClient(client KMSClient) *KMSCipher {
	return &KMSCipher{client}
}

// Encrypt will encrypt the "plainStr" using the region and keyID of the cipher.
func (c *KMSCipher) Encrypt(ctx context.Context, keyID string, plainStr string) (string, error) {
	return c.Client.Encrypt(ctx, keyID, plainStr)
}

// Decrypt will decrypt the "encryptedStr" using the region and keyID of the cipher.
func (c *KMSCipher) Decrypt(ctx context.Context, keyID string, encryptedStr string) (string, error) {
	return c.Client.Decrypt(ctx, keyID, encryptedStr)
}
