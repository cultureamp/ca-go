package encryption

import (
	"context"

	"github.com/pkg/errors"
)

type encryption struct {
	encryptor Encryptor
}

type Encryptor interface {
	Encrypt(ctx context.Context, plainStr string) (encryptedStr *string, err error)
	Decrypt(ctx context.Context, encryptedStr string) (decryptedStr *string, err error)
}

func NewEncryption(encryptor Encryptor) Encryptor {
	return &encryption{encryptor}
}

func (a *encryption) Decrypt(ctx context.Context, encriptedStr string) (*string, error) {
	deryptedStr, err := a.encryptor.Decrypt(ctx, encriptedStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt the string")
	}

	return deryptedStr, nil
}

func (a *encryption) Encrypt(ctx context.Context, plainStr string) (*string, error) {
	encryptedString, err := a.encryptor.Encrypt(ctx, plainStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt data")
	}
	return encryptedString, nil
}
