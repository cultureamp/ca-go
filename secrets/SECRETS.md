# ca-go/secrets

The `secrets` package wraps the AWS SecretManager in a simple to use sington pattern that you can call directly.

## Environment Variables

You MUST set these:
- AWS_REGION = The AWS region this code is running in (eg. "us-west-1")

## Examples
```
package cago

import (
	"fmt"

	"github.com/cultureamp/ca-go/secrets"
)

func BasicExamples() {
	// this will automatically use the AWS Region as per the environment variable "AWS_REGION"
	answer, err := secrets.Get("my-test-secret")
	fmt.Printf("The answer to the secret is '%s' (err='%v')\n", answer, err)
}
```
