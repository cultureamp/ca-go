# ca-go/cipher

The `cipher` package provides access to kms Encrpyt and Decrpyt. The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

The package creates a default cipher that uses the `AWS_REGION` environment variable. For ease of use, it is recommended that you use the package level methods Encrypt and Decrypt. However, if you need to support another region then you can create a `NewKMSCipher("region")` and manage the class life-cycle yourself.

The `keyId` parameter should be set to something like: `arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab`

The encrypted string will be a base64 string and look something like: `AQICAHgk4KLG1nZnyA8JokTKxExg+91EVz8GZMtgV5r0ImKJ2QFYCP9IuBbv1w4vduDowQYRAAABLTCCASkGCSqGSIb3DQEHBqCCARowggEWAgEAMIIBDwYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAwYdO9mDUMoKgH+9YACARCAgeEwVDZwFtIhdBL6JO2wNrcPyxdBEcTDbnqI81MyMSvNyGMEqZvZKQCHQElShUsHVqvIiW49KpCWvbbhzn6iPekYd+qaio59+mk4+AIMmQE8L43qMTKOobC/pUZeqQ1M/fqGqtzXpU0ezFhVMc7nDaVBj6VraQhCsaTuN4ZrJtRTD0c/SFcFXNvP0iN6wGaQAmU+TGIdK3Q9qOdCAp2k1254RrxM/A8Xtaw9cOJZea0e0d9O+IcET30vwLKNBy2ut96pPkAJCDuM6Gkvb8rHmjk69Ft7ClLKmSdKlYSS+WawPto=`

## Environment Variables

Here is the list of supported environment variables currently supported:
- AwsRegionEnv    = "AWS_REGION"

## Methods

- func Encrypt(ctx context.Context, keyId string, plainStr string) (string, error)
- func Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error)

## Examples

```
package cago

import (
	"context"
	"fmt"
	"os"

	"github.com/cultureamp/ca-go/cipher"
)

func Example() {
	ctx := context.Background()

	os.SetEnv("AWS_REGION", "us-west-2")

	// Replace the following example key ARN with any valid key identfier
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab";

	// this will automatically use the environment variable "AWS_REGION" 
	encrypted, err := cipher.Encrypt(ctx, keyId, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)

	decrypted, err := cipher.Decrypt(ctx, keyId, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)

	// or if you need cipher for another region use
	crypto := cipher.NewKMSCipher("region")
	encrypted, err = crypto.Encrypt(ctx, keyId, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)

	decrypted, err = crypto.Decrypt(ctx, keyId, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)
}
```
## Testing and Mocks

During tests you can override the package level `DefaultKMSCipher.Client` with a mock that supports the `KMSClient` interface.

- Encrypt(ctx context.Context, keyId string, plainStr string) (string, error)
- Decrypt(ctx context.Context, keyId string, encryptedStr string) (string, error)

```
import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/cipher"
	"github.com/stretchr/testify/assert"
)

func TestPackageEncrypt(t *testing.T) {
	ctx := context.Background()
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"

	// replace the package level client with our mock
	stdClient := cipher.DefaultKMSCipher.Client
	cipher.DefaultKMSCipher.Client = newMockedCipherClient()
	defer func() {
		cipher.DefaultKMSCipher.Client = stdClient
	}()

	cipherText, err := cipher.Encrypt(ctx, keyId, "test_plain_str")
	assert.Nil(t, err)

	plainText, err := cipher.Decrypt(ctx, keyId, cipherText)
	assert.Nil(t, err)
	assert.Equal(t, "test_plain_str", plainText)
}

type mockedCipherClient struct{}

func newMockedCipherClient() *mockedCipherClient {
	return &mockedCipherClient{}
}

// Encrypt on the test runner just returns the "plainStr" as the encrypted encryptedStr.
func (c *mockedCipherClient) Encrypt(ctx context.Context, _ string, plainStr string) (string, error) {
	return plainStr, nil
}

// Decrypt on the test runner just returns the "encryptedStr" as the decrypted plainstr.
func (c *mockedCipherClient) Decrypt(ctx context.Context, _ string, encryptedStr string) (string, error) {
	return encryptedStr, nil
}
```
