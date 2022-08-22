package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/x/vault/client"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	renewClient func(ctx context.Context) error
	getSecret   func(batch []interface{}, keyReference string, action string) (*vaultapi.Secret, error)
}

func (m MockClient) RenewClient(ctx context.Context) error {
	return m.renewClient(ctx)
}
func (m MockClient) GetSecret(batch []interface{}, keyReference string, action string) (*vaultapi.Secret, error) {
	return m.getSecret(batch, keyReference, action)
}

const (
	decryptedString = "abc123!?$*&()'-=@~"
)

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name            string
		decryptedSecret []string
		data            map[string]interface{}
		shouldRenew     bool
		returnErr       error
		expErr          error
	}{
		{
			"should not error when secret is returned as usual",
			[]string{decryptedString},
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
					},
				},
			},
			false,
			nil,
			nil,
		},
		{
			"should error when secret is returns wrong number of results",
			nil,
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
					},
					map[string]interface{}{
						"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
					},
				},
			},
			false,
			nil,
			fmt.Errorf("incorrect number of decrypted values returned"),
		},
		{
			"should error when getSecret errors",
			nil,
			nil,
			false,
			fmt.Errorf("secretError"),
			fmt.Errorf("secretError"),
		},
		{
			"should renew then continue when getSecret returns permission error",
			[]string{decryptedString},
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
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
			nil,
			false,
			fmt.Errorf(client.VaultPermissionError),
			fmt.Errorf(client.VaultPermissionError),
		},
		{
			"should error when batch_results is not []interface{}",
			nil,
			map[string]interface{}{
				"batch_results": map[string]interface{}{
					"plaintext": decryptedString,
				},
			},
			false,
			nil,
			fmt.Errorf("batch results of decryption secret could not be cast to []interface{}"),
		},
		{
			"should error when batch_results entries are not map[string]interface{}",
			nil,
			map[string]interface{}{
				"batch_results": []interface{}{"plaintext"},
			},
			false,
			nil,
			fmt.Errorf("batch result decryption element is not map[string]interface{}"),
		},
		{
			"should error when not base64 encoded",
			nil,
			map[string]interface{}{
				"batch_results": []interface{}{
					map[string]interface{}{
						"plaintext": decryptedString,
					},
				},
			},
			false,
			nil,
			base64.CorruptInputError(6),
		},
	}
	keyReferences := []string{"keyRef1"}
	encryptedData := []string{"encrypted1"}
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
					Data:          tt.data,
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
			v := NewVaultDecrypter(mockClient, client.VaultSettings{
				RoleArn:   "arn:1234",
				VaultAddr: "1234",
			})
			decryptedSecret, err := v.Decrypt(keyReferences, encryptedData, ctx)
			assert.Equal(t, tt.decryptedSecret, decryptedSecret)
			assert.Equal(t, tt.shouldRenew, renewed)
			fmt.Printf("tt err: %v, err: %v\n", tt.expErr, err)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
