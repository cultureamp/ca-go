package cryptography

import (
	"context"
)

// KMSCryptography supports basic Encrypt & Decrypt methods.
type KMSCryptography struct {
	client KMSClient
}

// NewKMSCryptography creates a new kms cryptography for the specific "region" and "keyid".
func NewKMSCryptography(region string, keyID string) *KMSCryptography {
	client := newAWSKMSClient(region, keyID)
	return &KMSCryptography{client}
}

// NewKMSCryptography creates a new kms cryptography for the specific "region" and "keyid".
func NewKMSCryptographyWithClient(client KMSClient) *KMSCryptography {
	return &KMSCryptography{client}
}

// Encrypt will encrypt the "plainStr" using the region and keyID of the cryptography.
func (c *KMSCryptography) Encrypt(ctx context.Context, plainStr string) (string, error) {
	return c.client.Encrypt(ctx, plainStr)
}

// Decrypt will decrypt the "encryptedStr" using the region and keyID of the cryptography.
func (c *KMSCryptography) Decrypt(ctx context.Context, encryptedStr string) (string, error) {
	return c.client.Decrypt(ctx, encryptedStr)
}
