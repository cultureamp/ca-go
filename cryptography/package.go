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
		client = newAWSKMSClient(region)
	}
	return NewKMSCryptographyWithClient(client)
}

// Encrypt uses the default AWS_REGION to kms encrypt "plainStr".
func Encrypt(ctx context.Context, keyId string, plainStr string) (string, error) {
	return DefaultKMSCryptogrphy.Encrypt(ctx, keyId, plainStr)
}

// Decrypt uses the default AWS_REGION and KMS_KEY_ID to kms decrypt "encryptedStr".
func Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error) {
	return DefaultKMSCryptogrphy.Decrypt(ctx, keyId, encryptedStr)
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
