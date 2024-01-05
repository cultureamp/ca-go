package secrets_test

import (
	"fmt"

	"github.com/cultureamp/ca-go/secrets"
)

func BasicExamples() {
	// this will automatically use the AWS Region as per the environment variable "AWS_REGION"
	answer, err := secrets.Get("my-test-secret")
	fmt.Printf("The answer to the secret is '%s' (err='%v')\n", answer, err)
}
