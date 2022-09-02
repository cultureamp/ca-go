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
		keyRefs         []string
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
			[]string{"keyRef1"},
			false,
			nil,
			nil,
		},
		{
			"should error when no keys are given",
			nil,
			nil,
			[]string{},
			false,
			nil,
			client.ErrVaultMissingKeys,
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
			[]string{"keyRef1"},
			false,
			nil,
			fmt.Errorf("incorrect number of decrypted values returned"),
		},
		{
			"should error when getSecret errors",
			nil,
			nil,
			[]string{"keyRef1"},
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
			[]string{"keyRef1"},
			true,
			fmt.Errorf(client.VaultPermissionError),
			nil,
		},
		{
			"should error when renewClient returns error",
			nil,
			nil,
			[]string{"keyRef1"},
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
			[]string{"keyRef1"},
			false,
			nil,
			fmt.Errorf("batch results casting error"),
		},
		{
			"should error when batch_results entries are not map[string]interface{}",
			nil,
			map[string]interface{}{
				"batch_results": []interface{}{"plaintext"},
			},
			[]string{"keyRef1"},
			false,
			nil,
			fmt.Errorf("batch result casting error"),
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
			[]string{"keyRef1"},
			false,
			nil,
			base64.CorruptInputError(6),
		},
	}
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
			v := NewVaultDecrypter(mockClient)
			decryptedSecret, err := v.Decrypt(tt.keyRefs, encryptedData, ctx)
			assert.Equal(t, tt.decryptedSecret, decryptedSecret)
			assert.Equal(t, tt.shouldRenew, renewed)
			fmt.Printf("tt err: %v, err: %v\n", tt.expErr, err)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
