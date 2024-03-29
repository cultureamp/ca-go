package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/cultureamp/ca-go/x/log"
	"github.com/cultureamp/ca-go/x/vault/client"
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
	if len(keyReferences) < 1 {
		return nil, client.ErrVaultMissingKeys
	}
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
		logger.Error().Err(err).Msg("decryption secret qty err")
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

	// There are scenarios in which secret will be nil, even if there is no error.
	// This can happen on the decryption endpoint for a 404 status code and an empty response body
	// https://github.com/hashicorp/vault/blob/601ad4823cb5b21ede5bf4fc6cbdf638a02feebd/api/logical.go#L241-L249
	// We need to handle this case otherwise the nil pointer will be dereferenced.
	if secret == nil {
		errMsg := fmt.Sprintf("tried to decrypt keyReference: %s but vault returned an empty body", keyReference)
		logger.Error().Msg(errMsg)
		return nil, errors.New(errMsg)
	}

	batchResults, ok := secret.Data["batch_results"].([]interface{})
	if !ok {
		err = fmt.Errorf("batch results casting error")
		logger.Error().Err(err).Msg("batch results of decryption secret could not be cast to []interface{}")
		return nil, err
	}

	var result []string
	for _, r := range batchResults {
		rmap, ok := r.(map[string]interface{})
		if !ok {
			err = fmt.Errorf("batch result casting error")
			logger.Error().Err(err).Msg("batch result decryption element is not map[string]interface{}")
			return nil, err
		}
		plaintext := fmt.Sprintf("%v", rmap["plaintext"])
		base64Decoded, err := base64.StdEncoding.DecodeString(plaintext)
		if err != nil {
			logger.Error().Err(err).Msg("Error base64 decoding")
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
					logger.Error().Err(err).Msg("unable to renew vault client")
					return nil, err
				}
				continue
			}
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
