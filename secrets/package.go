package secrets

import (
	"context"
	"os"
)

// Secrets can be mocked by clients for testing purposes.
type Secrets interface {
	Get(ctx context.Context, secretKey string) (string, error)
}

// DefaultAWSSecretsManager is a public *AWSSecretsManager used for package level methods.
var DefaultAWSSecretsManager Secrets = nil

// Get retrives the secret from AWS SecretsManager.
func Get(ctx context.Context, secretKey string) (string, error) {
	mustHaveSecretsManager(ctx)

	return DefaultAWSSecretsManager.Get(ctx, secretKey)
}

func mustHaveSecretsManager(ctx context.Context) {
	if DefaultAWSSecretsManager != nil {
		return // its set so we are good to go
	}

	region := os.Getenv("AWS_REGION")
	sm, err := NewAWSSecretsManager(ctx, region)
	if err != nil {
		panic(err)
	}

	DefaultAWSSecretsManager = sm
}
