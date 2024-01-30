package secrets_test

import (
	"fmt"
	"github.com/cultureamp/ca-go/secrets"
)

func Example() {

	answer, _ := secrets.Get("my-test-secret")
	fmt.Printf("test secret is '%s'\n", answer)

	// or if you need secrets from another region other than the one you are running in use
	//s := secrets.NewAWSSecrets("a-different-region")
	//answer, _ = s.Get("my-test-secret2")
	//fmt.Printf("The answer to my secret from a different region is '%s'\n", answer)
	ms := new(mockSecrets)
	secrets.SetImpl(ms)
	answer, _ = secrets.Get("my-test-secret2")
	fmt.Printf("mocked secret is '%s' \n", answer)

	//Output:
	// test secret is 'test-secret-value'
	// mocked secret is 'mock-secret-value'
}

type mockSecrets struct {
}

func (*mockSecrets) Get(key string) (string, error) {
	return "mock-secret-value", nil
}
