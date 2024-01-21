package secrets

import (
	"os"
)

type Secrets interface {
	Get(name string) (string, error)
}

var defaultSecrets = getInstance()

func getInstance() Secrets {
	region := os.Getenv("AWS_REGION")
	return NewAWSSecrets(region)
}

// Get retrives the secret from AWS SecretsManager.
func Get(secretName string) (string, error) {
	return defaultSecrets.Get(secretName)
}
