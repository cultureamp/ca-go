package secrets_test

import (
	"fmt"
	"github.com/cultureamp/ca-go/secrets"
)

func Example() {
	// this will automatically use the AWS Region as per the environment variable "AWS_REGION"
	answer, err := secrets.Get("my-test-secret")
	fmt.Printf("The answer to the secret is '%s' (err='%v')\n", answer, err)

	// or if you need secrets from another region other than the one you are running in use
	//s := secrets.NewAWSSecrets("a-different-region")
	//answer, err = s.Get("my-test-secret2")
	//fmt.Printf("The answer to the secret2 is '%s' (err='%v')\n", answer, err)
	ms := new(mockSecrets)
	secrets.SetImpl(ms)
	answer, err = secrets.Get("my-test-secret2")
	fmt.Printf("The answer to the secret2 is '%s' (err='%v')\n", answer, err)

	//Output:
	// The answer to the secret is test-secret-value
}

type mockSecrets struct {
}

func (*mockSecrets) Get(key string) (string, error) {
	return "my-mock-secret", nil
}
