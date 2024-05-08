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
var DefaultAWSSecretsManager Secrets = nil //nolint:revive

// Get retrives the secret from AWS SecretsManager.
func Get(ctx context.Context, secretKey string) (string, error) {
	err := mustHaveSecretsManager(ctx)
	if err != nil {
		return "", err
	}

	return DefaultAWSSecretsManager.Get(ctx, secretKey)
}

func mustHaveSecretsManager(ctx context.Context) error {
	if DefaultAWSSecretsManager != nil {
		return nil // its set so we are good to go
	}

	region := os.Getenv("AWS_REGION")
	sm, err := NewAWSSecretsManager(ctx, region)
	if err != nil {
		return err
	}

	DefaultAWSSecretsManager = sm
	return nil
}
