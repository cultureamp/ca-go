package encryption

import (
	"context"
	"encoding/json"

	"github.com/cultureamp/ca-go/x/kms"

	"github.com/pkg/errors"
)

type Encryption interface {
	Decrypt(ctx context.Context, encriptedStr string) (decryptedStr *string, err error)
	Encrypt(ctx context.Context, data interface{}) (encriptedStr *string, err error)
}

type encryption struct {
	encryptor kms.KMSEncrypt
}

func NewEncryption(encryptor kms.KMSEncrypt) (encsrv Encryption) {
	return &encryption{encryptor}
}

func (a *encryption) Decrypt(ctx context.Context, encriptedStr string) (*string, error) {
	deryptedStr, err := a.encryptor.Decrypt(ctx, encriptedStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt the string")
	}

	return deryptedStr, nil
}

func (a *encryption) Encrypt(ctx context.Context, data interface{}) (*string, error) {
	var encryptedString *string
	if data != nil {
		dataByte, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal data")
		}

		encryptedString, err = a.encryptor.Encrypt(ctx, string(dataByte))
		if err != nil {
			return nil, errors.Wrap(err, "failed to encrypt data")
		}
	}
	return encryptedString, nil
}
