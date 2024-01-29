package secrets

import (
	"flag"
	"os"
	"strings"
)

// DefaultAWSSecrets is a public *AWSSecretsManager used for package level methods.
var DefaultAWSSecrets = getInstance()

func getInstance() *AWSSecretsManager {
	var client AWSSecretsManagerClient

	if isTestMode() {
		client = newTestRunnerClient()
	} else {
		// Should this take dependency on 'env' package and call env.AwsRegion()?
		region := os.Getenv("AWS_REGION")
		client = newSecretManagerClient(region)
	}
	return NewAWSSecretsManagerWithClient(client)
}

// Get retrives the secret from AWS SecretsManager.
func Get(secretName string) (string, error) {
	return DefaultAWSSecrets.Get(secretName)
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
