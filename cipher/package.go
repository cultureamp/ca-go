package cipher

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/go-errors/errors"
)

// DefaultKMSCipher is the package level default implementation used by all package level methods.
// Package level methods are provided for ease of use.
// For testing you can replace the DefaultKMSCipher client with your own mock:
//
//	DefaultKMSCipher = newMockedClient()
var DefaultKMSCipher KMSCipher = getInstance()

func getInstance() *kmsCipher {
	var client KMSCipher

	region, ok := os.LookupEnv("AWS_REGION")
	if !ok || region == "" {
		if !isTestMode() {
			err := errors.Errorf("missing AWS_REGION environment variable")
			panic(err)
		}

		err := os.Setenv("AWS_REGION", "dev")
		if err != nil {
			panic(err)
		}
	}

	client = newAWSKMSClient(region)
	return NewKMSCipherWithClient(client)
}

// Encrypt will use env var AWS_REGION and the KMS keyId to encrypt the plainStr and return it as a base64 encoded string.
func Encrypt(ctx context.Context, keyId string, plainStr string) (string, error) {
	return DefaultKMSCipher.Encrypt(ctx, keyId, plainStr)
}

// Decrpyt will use env var AWS_REGION and the KMS keyId and the base64 encoded encryptedStr and return it decrypted as a plain string.
func Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error) {
	return DefaultKMSCipher.Decrypt(ctx, keyId, encryptedStr)
}

func isTestMode() bool {
	// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	argZero := os.Args[0]

	if strings.HasSuffix(argZero, ".test") ||
		strings.Contains(argZero, "/_test/") ||
		strings.Contains(argZero, "__debug_bin") || // vscode debug binary
		flag.Lookup("test.v") != nil {
		return true
	}

	return false
}
