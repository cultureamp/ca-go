package secrets

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
)

// DefaultAWSSecrets is a public *AWSSecretsManager used for package level methods.
var DefaultAWSSecrets = getInstance()

func getInstance() *AWSSecretsManager {
	var client AWSSecretsManagerClient

	if isTestMode() {
		client = newTestRunnerClient()
	} else {
		region := os.Getenv("AWS_REGION")
		cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
		if err != nil {
			err := fmt.Errorf("error loading aws sdk config, err='%w'\n", err)
			panic(err)
		}

		client = newSecretManagerClient(cfg)
	}
	return NewAWSSecretsManagerWithClient(client)
}

// Get retrives the secret from AWS SecretsManager.
func Get(ctx context.Context, secretKey string) (string, error) {
	return DefaultAWSSecrets.Get(ctx, secretKey)
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
