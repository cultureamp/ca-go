// package encryption provides interface Encryptor for other packages i.e. "kms" which could implement
// thier own functions for encryption and decryption
//
// ```
//
//	type Encryptor interface {
//		Encrypt(ctx context.Context, plainStr string) (encryptedStr *string, err error)
//		Decrypt(ctx context.Context, encryptedStr string) (decryptedStr *string, err error)
//	}
//
// ````
package encryption
