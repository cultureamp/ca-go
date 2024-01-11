package secrets

import (
	"os"
)

var DefaultAWSSecrets = getInstance()

func getInstance() *AWSSecretsManager {
	// Should this take dependency on 'env' package and call env.AwsRegion()?
	region := os.Getenv("AWS_REGION")
	client := NewSecretManagerClient(region)
	return NewAWSSecretsManagerWithClient(client)
}

// Get retrives the secret from AWS SecretsManager.
func Get(secretName string) (string, error) {
	return DefaultAWSSecrets.Get(secretName)
}
