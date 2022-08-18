// Package vaultdecryption adds helper methods for decrypting using vault
// To Decrypt, you must first create a client and decrypter:
// settings := client.VaultSettings{
////				DecrypterRoleArn: <arn here>,
////				VaultAddr:        <vault address here>,
////			}
// client, err :=  client.NewVaultClient(&settings, ctx context.Context)
//
// decrypter := NewVaultDecrypter(client, &settings})
//
// decryptedSecret, err := decrypter.Decrypt(keyReferences, encryptedData, ctx)
package vaultdecryption
