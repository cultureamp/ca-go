package secrets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSecretSuccess(t *testing.T) {

	result, err := Get("my-secret")
	assert.Nil(t, err)
	assert.Equal(t, "test-secret-value", result)
}

//
//func TestGetSecretOnError(t *testing.T) {
//	expectedOutput := &secretsmanager.GetOutput{
//		SecretString: aws.String("my-super-secret-value"),
//	}
//	mockedAwsClient := new(mockedSMClient)
//	mockedAwsClient.On("Get", mock.Anything).Return(expectedOutput, errors.New("test-error"))
//
//	secrets := NewAWSSecrets("us-west-2")
//	secrets.client = mockedAwsClient
//
//	result, err := secrets.Get("my-secret")
//	assert.NotNil(t, err)
//	assert.Equal(t, "", result)
//}
//
//func TestGetSecretOnEmpty(t *testing.T) {
//	expectedOutput := &secretsmanager.GetOutput{
//		SecretString: nil,
//	}
//	mockedAwsClient := new(mockedSMClient)
//	mockedAwsClient.On("Get", mock.Anything).Return(expectedOutput, nil)
//
//	secrets := NewAWSSecrets("us-west-2")
//	secrets.client = mockedAwsClient
//
//	result, err := secrets.Get("my-secret")
//	assert.NotNil(t, err)
//	assert.Equal(t, "", result)
//}
