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

type VaultDecrypter interface {
	Decrypt(keyReferences []string, encryptedData []string) ([]string, error)
}

type vaultDecrypter struct {
	ctx         context.Context
	vaultClient *vaultapi.Client
	logger      *log.Logger
	settings    *VaultSettings
}

func DefaultVaultDecrypter(ctx context.Context, settings *VaultSettings) (*vaultDecrypter, error) {
	client, err := DefaultVaultClients(ctx, settings).NewAwsIamVaultDecrypterClient()
	logger := log.NewFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault decrypter: %w", err)
	}
	return &vaultDecrypter{ctx, client, logger, settings}, nil
}

func NewVaultDecrypter(ctx context.Context, vaultClient *vaultapi.Client, settings *VaultSettings) *vaultDecrypter {
	logger := log.NewFromCtx(ctx)
	return &vaultDecrypter{ctx, vaultClient, logger, settings}
}

func (v vaultDecrypter) Decrypt(keyReferences []string, encryptedData []string) (decryptedData []string, err error) {
	span, _ := tracer.StartSpanFromContext(v.ctx, "vault-decrypter")
	defer span.Finish(tracer.WithError(err))

	result := encryptedData
	for _, keyReference := range reverse(keyReferences) {
		decryptedByKeyReference, err := v.decrypt(keyReference, result)
		if err != nil {
			return nil, fmt.Errorf("error decrypting with key reference %w", err)
		}
		result = decryptedByKeyReference
	}

	return result, nil
}

func (v vaultDecrypter) decrypt(keyReference string, encryptedData []string) ([]string, error) {
	var batch []interface{}
	for _, field := range encryptedData {
		batch = append(batch, map[string]interface{}{
			"ciphertext": field,
		})
	}

	secret, err := v.decryptWithVault(keyReference, batch)
	if err != nil {
		return nil, fmt.Errorf("error decrypting with Vault %w", err)
	}

	batchResults := secret.Data["batch_results"]
	var result []string
	for _, r := range batchResults.([]interface{}) {
		plaintext := fmt.Sprintf("%v", r.(map[string]interface{})["plaintext"])
		base64Decoded, err := base64.StdEncoding.DecodeString(plaintext)
		if err != nil {
			return nil, fmt.Errorf("error base64 decoding %w", err)
		}
		result = append(result, string(base64Decoded))
	}

	return result, nil
}

func (v vaultDecrypter) decryptWithVault(keyReference string, batch []interface{}) (secret *vaultapi.Secret, err error) {
	for i := 0; i < maxRetries; i++ {
		secret, err = v.vaultClient.Logical().Write(fmt.Sprintf("transit/decrypt/%s", keyReference), map[string]interface{}{
			"batch_input": batch,
		})

		if err != nil {
			if strings.Contains(err.Error(), vaultPermissionError) {
				client, e := DefaultVaultClients(v.ctx, v.settings).NewAwsIamVaultDecrypterClient()
				if e != nil {
					return nil, fmt.Errorf("unable to initialize Vault decrypter: %w", e)
				}
				v.vaultClient = client
				v.logger.Info("Renewing vault client", log.Fields{
					"err": err.Error(),
				})
				continue
			}

			v.logger.Error("Vault client returned unhandled error", err)
			return nil, fmt.Errorf("error calling vault decrypt API %w", err)
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
