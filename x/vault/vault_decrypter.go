package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cultureamp/ca-go/x/vault/client"
	"github.com/cultureamp/glamplify/log"
	vaultapi "github.com/hashicorp/vault/api"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Client interface {
	RenewClient(ctx context.Context) error
	GetSecret(batch []interface{}, keyReference string, action string) (*vaultapi.Secret, error)
}

type Decrypter struct {
	vaultClient Client
}

func NewVaultDecrypter(vaultClient Client) *Decrypter {
	return &Decrypter{vaultClient}
}

func (v *Decrypter) Decrypt(keyReferences []string, encryptedData []string, ctx context.Context) ([]string, error) {
	var err error
	span, _ := tracer.StartSpanFromContext(ctx, "vault-decrypter")
	defer span.Finish(tracer.WithError(err))
	logger := log.NewFromCtx(ctx)
	result := encryptedData
	for _, keyReference := range reverse(keyReferences) {
		decryptedByKeyReference, err := v.decryptByKey(keyReference, result, logger, ctx)
		if err != nil {
			return nil, err
		}
		result = decryptedByKeyReference
	}
	if len(result) != len(encryptedData) {
		err := fmt.Errorf("incorrect number of decrypted values returned")
		logger.Error("decryption secret qty err", err)
		return nil, err
	}

	return result, nil
}

func (v *Decrypter) decryptByKey(keyReference string, encryptedData []string, logger *log.Logger, ctx context.Context) ([]string, error) {
	var batch []interface{}
	for _, field := range encryptedData {
		batch = append(batch, map[string]interface{}{
			"ciphertext": field,
		})
	}

	secret, err := v.decryptWithVault(keyReference, batch, logger, ctx)
	if err != nil {
		return nil, err
	}

	batchResults, ok := secret.Data["batch_results"].([]interface{})
	if !ok {
		errStr := "batch results of decryption secret could not be cast to []interface{}"
		err = fmt.Errorf(errStr)
		logger.Error(errStr, err)
		return nil, err
	}

	var result []string
	for _, r := range batchResults {
		rmap, ok := r.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("batch result decryption element is not map[string]interface{}")
			logger.Error("batch result casting error", err)
			return nil, err
		}
		plaintext := fmt.Sprintf("%v", rmap["plaintext"])
		base64Decoded, err := base64.StdEncoding.DecodeString(plaintext)
		if err != nil {
			logger.Error("Error base64 decoding", err)
			return nil, err
		}
		result = append(result, string(base64Decoded))
	}
	return result, nil
}

func (v *Decrypter) decryptWithVault(keyReference string, batch []interface{}, logger *log.Logger, ctx context.Context) (*vaultapi.Secret, error) {
	var secret *vaultapi.Secret
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		secret, err = v.vaultClient.GetSecret(batch, keyReference, client.DecryptionAction)
		if err != nil {
			if strings.Contains(err.Error(), client.VaultPermissionError) {
				err = v.vaultClient.RenewClient(ctx)
				if err != nil {
					logger.Info("unable to renew vault client", log.Fields{"err": err.Error()})
					return nil, err
				}
				continue
			}
			logger.Error("Vault client returned unhandled error", err)
			return nil, err
		} else {
			break
		}
	}

	return secret, err
}

func reverse(s []string) []string {
	var reversed []string

	for i := len(s) - 1; i >= 0; i-- {
		reversed = append(reversed, s[i])
	}

	return reversed
}
