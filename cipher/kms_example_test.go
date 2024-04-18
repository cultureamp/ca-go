package cipher_test

import (
	"context"
	"fmt"

	"github.com/cultureamp/ca-go/cipher"
)

func Example() {
	// Note: Make sure AWS_REGION is set in the environment

	ctx := context.Background()

	// Replace the following example key ARN with any valid key identfier
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"

	// this will automatically use the environment variable "AWS_REGION"
	encrypted, err := cipher.Encrypt(ctx, keyId, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)
	// encrypted will be base64 string and look something like this:
	// "AQICAHgk4KLG1nZnyA8JokTKxExg+91EVz8GZMtgV5r0ImKJ2QFYCP9IuBbv1w4vduDowQYRAAABLTCCASkGCSqGSIb3DQEHBqCCARowggEWAgEAMIIBDwYJKoZIhvcNAQcBMB4GCWCGSAFlAwQBLjARBAwYdO9mDUMoKgH+9YACARCAgeEwVDZwFtIhdBL6JO2wNrcPyxdBEcTDbnqI81MyMSvNyGMEqZvZKQCHQElShUsHVqvIiW49KpCWvbbhzn6iPekYd+qaio59+mk4+AIMmQE8L43qMTKOobC/pUZeqQ1M/fqGqtzXpU0ezFhVMc7nDaVBj6VraQhCsaTuN4ZrJtRTD0c/SFcFXNvP0iN6wGaQAmU+TGIdK3Q9qOdCAp2k1254RrxM/A8Xtaw9cOJZea0e0d9O+IcET30vwLKNBy2ut96pPkAJCDuM6Gkvb8rHmjk69Ft7ClLKmSdKlYSS+WawPto="

	decrypted, err := cipher.Decrypt(ctx, keyId, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)

	// or if you need cipher for another region or keyID then use
	crypto := cipher.NewKMSClient("region")
	encrypted, err = crypto.Encrypt(ctx, keyId, "plain-string")
	fmt.Printf("The encrypted string is '%s' (err='%v')\n", encrypted, err)

	decrypted, err = crypto.Decrypt(ctx, keyId, encrypted)
	fmt.Printf("The decrypted string is '%s' (err='%v')\n", decrypted, err)
}
