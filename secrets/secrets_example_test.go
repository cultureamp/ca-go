package secrets_test

import (
	"fmt"

	"github.com/cultureamp/ca-go/secrets"
)

func BasicExamples() {
	// this will automatically use the AWS Region as per the environment variable "AWS_REGION"
	answer, err := secrets.Get("my-test-secret")
	fmt.Printf("The answer to the secret is '%s' (err='%v')\n", answer, err)

	// or if you need secrets from another region other than the one you are running in use
	sm := secrets.NewAWSSecretsManager("a-different-region")
	answer, err = sm.Get("my-test-secret2")
	fmt.Printf("The answer to the secret2 is '%s' (err='%v')\n", answer, err)

	// of if you want to have a custom client that
	// supports the AWSSecretsManagerClient interface
	// myClient := newMyCustomClient()
	// secrets := NewAWSSecretsManagerWithClient(myClient)
	// or overwrite the default package implementation with
	// DefaultAWSSecrets.Client = myClient

	// or if you want to be able to mock the behavior
	// create a mock that supports the AWSSecretsManagerClient interface
	// mockedClient := new(mockedAWSSecretsManagerClient)
	// mockedClient.On("GetSecretValue", mock.Anything).Return(expectedOutput, nil)
	// re-assign the client
	// secrets := NewAWSSecretsManagerWithClient(mockedClient)
	// or overwrite the default package implementation with
	// DefaultAWSSecrets.Client = mockedClient
}
