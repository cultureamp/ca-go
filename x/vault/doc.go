// Package vault adds helper methods for decrypting using vault
// To Decrypt, you must first create a client and decrypter:
// settings := client.VaultSettings{
////				RoleArn: <arn here>,
////				VaultAddr:        <vault address here>,
////			}
// client, err :=  client.NewVaultClient(&settings, ctx context.Context)
//
// decrypter := NewVaultDecrypter(client, &settings})
//
// decryptedSecret, err := decrypter.Decrypt(keyReferences, encryptedData, ctx)
//
// The default aws region for generating AWS login data is 'us-east-1' which
// can be changed by using the env var AWS_REGION
package vault
