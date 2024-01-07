package cryptography_test

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
