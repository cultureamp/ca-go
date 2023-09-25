package encryption

import "context"

type Encryptor interface {
	Encrypt(ctx context.Context, plainStr string) (encryptedStr *string, err error)
	Decrypt(ctx context.Context, encryptedStr string) (decryptedStr *string, err error)
}
