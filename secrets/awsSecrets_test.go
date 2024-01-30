package secrets

import (
	"testing"
)

func TestNewAWSSecretsClient(t *testing.T) {
	client := NewAWSSecrets("us-west-2")
	if client == nil {
		t.Fatal("err")
	}
}
