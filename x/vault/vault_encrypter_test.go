package vault

import (
	"context"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/x/vault/client"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	const encStr = "encrypted string"
	tests := []struct {
		name        string
		secretData  map[string]interface{}
		shouldRenew bool
		returnErr   error
		expErr      error
	}{
		{
			"should not error when one secret is returned as usual",
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"ciphertext": encStr,
					},
				},
			},
			false,
			nil,
			nil,
		},
		{
			"should error when secret is returns wrong number of results",
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"ciphertext": encStr,
					},
					map[string]interface{}{
						"plaintext": encStr,
					},
				},
			},
			false,
			nil,
			fmt.Errorf("encryption secret qty err"),
		},
		{
			"should not error when multiple secrets are returned as usual",
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"ciphertext": encStr,
					},
				},
			},
			false,
			nil,
			nil,
		},
		{
			"should error when getSecret errors",
			nil,
			false,
			fmt.Errorf("secretError"),
			fmt.Errorf("secretError"),
		},
		{
			"should renew then continue when getSecret returns permission error",
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"ciphertext": encStr,
					},
				},
			},
			true,
			fmt.Errorf(client.VaultPermissionError),
			nil,
		},
		{
			"should error when renewClient returns error",
			nil,
			false,
			fmt.Errorf(client.VaultPermissionError),
			fmt.Errorf(client.VaultPermissionError),
		},
		{
			"should error when batch_results is not []interface{}",
			map[string]interface{}{
				"batch_results": map[string]interface{}{
					"cyphertext": decryptedString,
				},
			},
			false,
			nil,
			fmt.Errorf("batch results of encryption secret could not be cast to []interface{}"),
		},
		{
			"should error when batch_results entries are not map[string]interface{}",
			map[string]interface{}{
				"batch_results": []interface{}{"ciphertext"},
			},
			false,
			nil,
			fmt.Errorf("encrypt batch result element is not map[string]interface{}"),
		},
	}
	keyReferences := []string{"keyRef1"}
	decryptedData := []string{"decrypted1"}
	ctx := context.Background()
	for _, tt := range tests {
		renewed := false
		mockClient := MockClient{
			renewClient: func(ctx context.Context) error {
				if !renewed && tt.shouldRenew {
					renewed = true
				} else if !tt.shouldRenew {
					return tt.returnErr
				}
				return nil
			},
			getSecret: func(batch []interface{}, keyReference string, action string) (*vaultapi.Secret, error) {
				secretReturn := vaultapi.Secret{
					RequestID:     "",
					LeaseID:       "",
					LeaseDuration: 0,
					Renewable:     false,
					Data:          tt.secretData,
					Warnings:      nil,
					Auth:          nil,
					WrapInfo:      nil,
				}
				if renewed {
					tt.returnErr = nil
				}
				return &secretReturn, tt.returnErr
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			v := NewVaultEncrypter(mockClient)
			_, err := v.Encrypt(keyReferences, decryptedData, ctx)
			assert.Equal(t, tt.shouldRenew, renewed)
			fmt.Printf("tt err: %v, err: %v\n", tt.expErr, err)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
