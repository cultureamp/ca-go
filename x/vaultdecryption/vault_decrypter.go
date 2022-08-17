package vaultdecryption

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cultureamp/glamplify/log"
	vaultapi "github.com/hashicorp/vault/api"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Client interface {
	RenewClient(ctx context.Context) error
	GetSecret(batch []interface{}, keyReference string) (*vaultapi.Secret, error)
}

type VaultDecrypter interface {
	Decrypt(keyReferences []string, encryptedData []string) ([]string, error)
}

type vaultDecrypter struct {
	vaultClient Client
	settings    *VaultSettings
}

func DefaultVaultDecrypter(ctx context.Context, settings *VaultSettings, client *VaultClient) (*vaultDecrypter, error) {
	client, err := NewVaultClient(settings, ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault decrypter: %w", err)
	}
	return &vaultDecrypter{client, settings}, nil
}

func NewVaultDecrypter(vaultClient Client, settings *VaultSettings) *vaultDecrypter {
	return &vaultDecrypter{vaultClient, settings}
}

func (v vaultDecrypter) Decrypt(keyReferences []string, encryptedData []string, ctx context.Context) ([]string, error) {
	logger := log.NewFromCtx(ctx)
	var err error
	span, _ := tracer.StartSpanFromContext(ctx, "vault-decrypter")
	defer span.Finish(tracer.WithError(err))

	result := encryptedData
	for _, keyReference := range reverse(keyReferences) {
		decryptedByKeyReference, err := v.decryptByKey(keyReference, result, *logger, ctx)
		if err != nil {
			return nil, fmt.Errorf("error decrypting with key reference %w", err)
		}
		result = decryptedByKeyReference
	}

	return result, nil
}

func (v vaultDecrypter) decryptByKey(keyReference string, encryptedData []string, logger log.Logger, ctx context.Context) ([]string, error) {
	var batch []interface{}
	for _, field := range encryptedData {
		batch = append(batch, map[string]interface{}{
			"ciphertext": field,
		})
	}

	secret, err := v.decryptWithVault(keyReference, batch, logger, ctx)
	if err != nil {
		return nil, fmt.Errorf("error decrypting with Vault %w", err)
	}

	batchResults, ok := secret.Data["batch_results"].([]interface{})
	var result []string
	if ok {
		for _, r := range batchResults {
			rmap, ok := r.(map[string]interface{})
			if !ok {
				err = fmt.Errorf("batch result is not map[string]interface{}")
				logger.Error("batch result is not map[string]interface{}", err)
				return nil, err
			}
			plaintext := fmt.Sprintf("%v", rmap["plaintext"])
			base64Decoded, err := base64.StdEncoding.DecodeString(plaintext)
			if err != nil {
				logger.Error("Error base64 decoding", err)
				return nil, fmt.Errorf("error base64 decoding %w", err)
			}
			result = append(result, string(base64Decoded))
		}
	} else {
		errStr := "batch results of secret could not be cast to []interface{}"
		err = fmt.Errorf(errStr)
		logger.Error(errStr, err)
		return nil, err
	}

	return result, nil
}

func (v vaultDecrypter) decryptWithVault(keyReference string, batch []interface{}, logger log.Logger, ctx context.Context) (*vaultapi.Secret, error) {
	var secret *vaultapi.Secret
	var err error
	for i := 0; i < maxRetries; i++ {
		secret, err = v.vaultClient.GetSecret(batch, keyReference)
		if err != nil {
			if strings.Contains(err.Error(), vaultPermissionError) {
				err = v.vaultClient.RenewClient(ctx)
				if err != nil {
					logger.Info("unable to renew vault client", log.Fields{"err": err.Error()})
					return nil, fmt.Errorf("unable to initialize Vault decrypter: %w", err)
				}
				continue
			}
			logger.Error("Vault client returned unhandled error", err)
			return nil, fmt.Errorf("error calling vault decryptByKey API %w", err)
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
