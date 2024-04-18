package secrets_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/cultureamp/ca-go/secrets"
)

func BasicExamples() {
	ctx := context.Background()

	// this will automatically use the AWS Region as per the environment variable "AWS_REGION"
	answer, err := secrets.Get(ctx, "my-test-secret")
	fmt.Printf("The answer to the secret is '%s' (err='%v')\n", answer, err)

	// or if you need secrets from another region other than the one you are running in use
	sm, err := secrets.NewAWSSecretsManager(ctx, "a-different-region")
	answer, err = sm.Get(ctx, "my-test-secret2")
	fmt.Printf("The answer to the secret2 is '%s' (err='%v')\n", answer, err)

	// of if you want to have a custom client that
	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	smc := secretsmanager.NewFromConfig(cfg)
	sm = secrets.NewAWSSecretsManagerWithClient(smc)

	// or if you want to be able to mock the behavior
	mockSM := newTestRunner()
	oldSM := secrets.DefaultAWSSecretsManager
	defer func() { secrets.DefaultAWSSecretsManager = oldSM }()
	secrets.DefaultAWSSecretsManager = mockSM
}

type testRunner struct{}

func newTestRunner() *testRunner {
	return &testRunner{}
}

// Get on the test runner returns the key as the secret.
func (c *testRunner) Get(_ context.Context, key string) (string, error) {
	// do whatever you want here
	return key, nil
}
