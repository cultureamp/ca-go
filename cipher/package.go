package cipher

import (
	"context"
	"os"
)

var DefaultKMSCipher KMSCipher = getInstance()

func getInstance() *kmsCipher {
	var client KMSCipher

	region := os.Getenv("AWS_REGION")
	client = newAWSKMSClient(region)
	return NewKMSCipherWithClient(client)
}

// Encrypt uses the default AWS_REGION to kms encrypt "plainStr".
func Encrypt(ctx context.Context, keyId string, plainStr string) (string, error) {
	return DefaultKMSCipher.Encrypt(ctx, keyId, plainStr)
}

// Decrypt uses the default AWS_REGION and KMS_KEY_ID to kms decrypt "encryptedStr".
func Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error) {
	return DefaultKMSCipher.Decrypt(ctx, keyId, encryptedStr)
}
