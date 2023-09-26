// package encryption provides interface Encryptor for other packages i.e. "kms" which could implement
// their own "Encrypt" and "Decrypt". Also, it could be used in the consumer of this package for the
// returnning type.
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
