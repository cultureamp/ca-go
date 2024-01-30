# ca-go/secrets

The `secrets` package wraps the AWS SecretManager in a simple to use sington pattern that you can call directly.

## Environment Variables

You MUST set these:
- AWS_REGION = The AWS region this code is running in (eg. "us-west-1")

## FAQ

Question: I need to load secrets from another region? How do I do that?
Answer: You can create your own secrets with sm := NewAWSSecrets("region") and then call sm.Get("secret")

## Examples
```
package cago

import (
	"fmt"

	"github.com/cultureamp/ca-go/secrets"
)

func BasicExamples() {
	ctx := context.Background()

	// this will automatically use the AWS Region as per the environment variable "AWS_REGION"
	answer, err := secrets.Get(ctx, "my-test-secret")
	fmt.Printf("The answer to the secret is '%s' (err='%v')\n", answer, err)

	// or if you need secrets from another region other than the one you are running in use
	sm := secrets.NewAWSSecretsManager("a-different-region")
	answer, err = sm.Get(ctx, "my-test-secret2")
	fmt.Printf("The answer to the secret2 is '%s' (err='%v')\n", answer, err)
}
```
