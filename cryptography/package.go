package cryptography

import (
	"context"
	"flag"
	"os"
	"strings"
)

var DefaultKMSCryptogrphy *KMSCryptography = getInstance()

func getInstance() *KMSCryptography {
	var client KMSClient

	if isTestMode() {
		client = newTestRunnerClient()
	} else {
		region := os.Getenv("AWS_REGION")
		keyID := os.Getenv("KMS_KEY_ID")
		client = newAWSKMSClient(region, keyID)
	}
	return NewKMSCryptographyWithClient(client)
}

// Encrypt uses the default AWS_REGION and KMS_KEY_ID to kms encrypt "plainStr".
func Encrypt(ctx context.Context, plainStr string) (string, error) {
	return DefaultKMSCryptogrphy.Encrypt(ctx, plainStr)
}

// Decrypt uses the default AWS_REGION and KMS_KEY_ID to kms decrypt "encryptedStr".
func Decrypt(ctx context.Context, encryptedStr string) (string, error) {
	return DefaultKMSCryptogrphy.Decrypt(ctx, encryptedStr)
}

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}
