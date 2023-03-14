// Package vault adds helper methods for decrypting using vault
// To Decrypt or Encrypt, you must first create a client:
//
//	settings := client.VaultSettings{
//					RoleArn: <arn here>,
//					VaultAddr:        <vault address here>,
//				}
//
// client, err :=  client.NewVaultClient(&settings, ctx context.Context)
//
// Please note that you should use a different client for decrypting and encrypting
// as they will have different RoleArns
//
// decrypter := NewVaultDecrypter(client, &settings})
//
// decryptedSecret, err := decrypter.Decrypt(keyReferences, encryptedData, ctx)
//
// encrypter := NewVaultEncrypter(client, &settings})
//
// encryptedSecret, err := encrypter.Encrypt(keyReferences, decryptedKeys, ctx)
//
// # If no keys are passed to Encrypt or Decrypt a VaultMissingKeysErr will be thrown
//
// The default aws region for generating AWS login data is 'us-east-1'. This
// is fixed for all environments.
package vault
