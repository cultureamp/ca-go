# ca-go/cryptography

The `cryptography` package provides access to kms Encrpyt and Decrpyt. The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

The package creates a default cryptography that uses the `AWS_REGION` and `KMS_KEY_ID` environment variables. However, if you need to Encrpyt or Decrypt for another region or multiple KeyIDs then you can create a `NewCryptography("region", "keyID")` and manage the class yourself.

## Environment Variables

Here is the list of supported environment variables currently supported:
- AwsRegionEnv    = "AWS_REGION"
- AwsKmsKeyIDEnv  = "KMS_KEY_ID"

## Methods

func Encrypt(ctx context.Context, plainStr string) (string, error)
func Decrypt(ctx context.Context, encryptedStr string) (string, error)

## Examples
```
package cago

import (
	"context"
	"fmt"

	"github.com/cultureamp/ca-go/cryptography"
)

func BasicExamples() {
	ctx := context.Background()

	// this will automatically the environment variables "AWS_REGION" and KMS_KEY_ID
	encrypted, err := cryptography.Encrypt(ctx, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)

	decrypted, err := cryptography.Decrypt(ctx, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)

	// or if you need cryptogprahy for another region or keyID then use
	crypto := cryptography.NewCryptography("region", "keyID")
	encrypted, err = crypto.Encrypt(ctx, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)

	decrypted, err = crypto.Decrypt(ctx, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)
}
```
