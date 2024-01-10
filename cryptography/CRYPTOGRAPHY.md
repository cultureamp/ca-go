# ca-go/cryptography

The `cryptography` package provides access to kms Encrpyt and Decrpyt. The design of this package is to provide a simple system that can be used in a variety of situations without requiring high cognitive load.

The package creates a default cryptography that uses the `AWS_REGION` and `KMS_KEY_ID` environment variables. However, if you need to Encrpyt or Decrypt for another region or multiple KeyIDs then you can create a `NewKMSCryptography("region", "keyID")` and manage the class life-cycle yourself.

The `KMS_KEY_ID` environment variable should be set to something like: `arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab`

The encrypted string will be a base64 string and look something like: `AQICAHgk4KLG1nZnyA8JokTKxExg+91EVz8GZMtgV5r0ImKJ2QFYCP9IuBbv1w4vduDowQYRAAABLTCCASkGCSqGSIb3DQEHBqCCARowggEWAgEAMIIBDwYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAwYdO9mDUMoKgH+9YACARCAgeEwVDZwFtIhdBL6JO2wNrcPyxdBEcTDbnqI81MyMSvNyGMEqZvZKQCHQElShUsHVqvIiW49KpCWvbbhzn6iPekYd+qaio59+mk4+AIMmQE8L43qMTKOobC/pUZeqQ1M/fqGqtzXpU0ezFhVMc7nDaVBj6VraQhCsaTuN4ZrJtRTD0c/SFcFXNvP0iN6wGaQAmU+TGIdK3Q9qOdCAp2k1254RrxM/A8Xtaw9cOJZea0e0d9O+IcET30vwLKNBy2ut96pPkAJCDuM6Gkvb8rHmjk69Ft7ClLKmSdKlYSS+WawPto=`

## Environment Variables

Here is the list of supported environment variables currently supported:
- AwsRegionEnv    = "AWS_REGION"
- AwsKmsKeyIDEnv  = "KMS_KEY_ID"


## Methods

- func Encrypt(ctx context.Context, plainStr string) (string, error)
- func Decrypt(ctx context.Context, encryptedStr string) (string, error)

## Examples
```
package cago

import (
	"context"
	"fmt"
	"os"

	"github.com/cultureamp/ca-go/cryptography"
)

func BasicExamples() {
	ctx := context.Background()

	os.SetEnv("AWS_REGION", "us-west-2")

	// Replace the following example key ARN with any valid key identfier
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab";
	os.Setenv("KMS_KEY_ID", keyId)

	// this will automatically the environment variables "AWS_REGION" and KMS_KEY_ID
	encrypted, err := cryptography.Encrypt(ctx, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)

	decrypted, err := cryptography.Decrypt(ctx, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)

	// or if you need cryptogprahy for another region or keyID then use
	crypto := cryptography.NewKMSCryptography("region", "keyID")
	encrypted, err = crypto.Encrypt(ctx, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)

	decrypted, err = crypto.Decrypt(ctx, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)
}
```
