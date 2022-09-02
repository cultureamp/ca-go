package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cultureamp/ca-go/x/log"
	"github.com/cultureamp/ca-go/x/vault/client"
	vaultapi "github.com/hashicorp/vault/api"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Encrypter struct {
	vaultClient Client
}

func NewVaultEncrypter(vaultClient Client) *Encrypter {
	return &Encrypter{vaultClient}
}

func (v Encrypter) Encrypt(keyReferences []string, protectedData []string, ctx context.Context) ([]string, error) {
	var err error
	span, _ := tracer.StartSpanFromContext(ctx, "vault-encrypter")
	defer span.Finish(tracer.WithError(err))
	logger := log.NewFromCtx(ctx)
	if len(keyReferences) < 1 {
		return nil, client.ErrVaultMissingKeys
	}
	result := protectedData
	for _, keyReference := range keyReferences {
		encryptedByKeyReference, err := v.encrypt(keyReference, result, logger, ctx)
		if err != nil {
			return nil, err
		}
		result = encryptedByKeyReference
	}
	if len(result) != len(protectedData) {
		err := fmt.Errorf("encryption secret qty err")
		logger.Error().Err(err).Msg("incorrect number of encrypted values returned")
		return nil, err
	}

	return result, nil
}

func (v Encrypter) encrypt(keyReference string, protectedData []string, logger *log.Logger, ctx context.Context) ([]string, error) {
	var batch []interface{}
	for _, field := range protectedData {
		batch = append(batch, map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(field)),
		})
	}

	secret, err := v.encryptWithVault(keyReference, batch, logger, ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error encrypting with Vault")
		return nil, err
	}
	batchResults, ok := secret.Data["batch_results"].([]interface{})
	if !ok {
		errStr := "batch results of encryption secret could not be cast to []interface{}"
		err = fmt.Errorf(errStr)
		logger.Error().Err(err).Msg("batchResult casting error")
		return nil, err
	}
	var result []string
	for _, r := range batchResults {
		rmap, ok := r.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("encrypt batch result element is not map[string]interface{}")
			logger.Error().Err(err).Msg("encryption batch result casting error")
			return nil, err
		}
		ciphertext := fmt.Sprintf("%v", rmap["ciphertext"])
		result = append(result, ciphertext)
	}
	return result, nil
}

func (v Encrypter) encryptWithVault(keyReference string, batch []interface{}, logger *log.Logger, ctx context.Context) (*vaultapi.Secret, error) {
	maxRetries := 5
	var secret *vaultapi.Secret
	var err error
	for i := 0; i < maxRetries; i++ {
		secret, err = v.vaultClient.GetSecret(batch, keyReference, client.EncryptionAction)
		if err != nil {
			if strings.Contains(err.Error(), client.VaultPermissionError) {
				err = v.vaultClient.RenewClient(ctx)
				if err != nil {
					logger.Error().Err(err).Msg("unable to renew Vault encrypter")
					return nil, err
				}
				continue
			}
			logger.Error().Err(err).Msg("error calling vault encrypt API")
			return nil, err
		} else {
			break
		}
	}

	return secret, err
}
