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
		vaultResponse   *vaultapi.Secret
		data            map[string]interface{}
		keyRefs         []string
		shouldRenew     bool
		returnErr       error
		expErr          error
	}{
		{
			name:            "should not error when secret is returned as usual",
			decryptedSecret: []string{decryptedString},
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data: map[string]interface{}{
					"batch_results": []interface{}{
						map[string]interface{}{
							"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
						},
					},
				},
				Warnings: nil,
				Auth:     nil,
				WrapInfo: nil,
			},
			keyRefs:     []string{"keyRef1"},
			shouldRenew: false,
			returnErr:   nil,
			expErr:      nil,
		},
		{
			name:            "should error when no keys are given",
			decryptedSecret: nil,
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data:          nil,
				Warnings:      nil,
				Auth:          nil,
				WrapInfo:      nil,
			},
			data:        nil,
			keyRefs:     []string{},
			shouldRenew: false,
			returnErr:   nil,
			expErr:      client.ErrVaultMissingKeys,
		},
		{
			name:            "should error when secret is returns wrong number of results",
			decryptedSecret: nil,
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data: map[string]interface{}{
					"batch_results": []interface{}{
						map[string]interface{}{
							"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
						},
						map[string]interface{}{
							"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
						},
					},
				},
				Warnings: nil,
				Auth:     nil,
				WrapInfo: nil,
			},
			keyRefs:     []string{"keyRef1"},
			shouldRenew: false,
			returnErr:   nil,
			expErr:      fmt.Errorf("incorrect number of decrypted values returned"),
		},
		{
			name:            "should error when getSecret errors",
			decryptedSecret: nil,
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data:          nil,
				Warnings:      nil,
				Auth:          nil,
				WrapInfo:      nil,
			},
			data:        nil,
			keyRefs:     []string{"keyRef1"},
			shouldRenew: false,
			returnErr:   fmt.Errorf("secretError"),
			expErr:      fmt.Errorf("secretError"),
		},
		{
			name:            "should renew then continue when getSecret returns permission error",
			decryptedSecret: []string{decryptedString},
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data: map[string]interface{}{
					"batch_results": []interface{}{
						map[string]interface{}{
							"plaintext": base64.StdEncoding.EncodeToString([]byte(decryptedString)),
						},
					},
				},
				Warnings: nil,
				Auth:     nil,
				WrapInfo: nil,
			},
			keyRefs:     []string{"keyRef1"},
			shouldRenew: true,
			returnErr:   fmt.Errorf(client.VaultPermissionError),
			expErr:      nil,
		},
		{
			name:            "should error when renewClient returns error",
			decryptedSecret: nil,
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data:          nil,
				Warnings:      nil,
				Auth:          nil,
				WrapInfo:      nil,
			},
			keyRefs:     []string{"keyRef1"},
			shouldRenew: false,
			returnErr:   fmt.Errorf(client.VaultPermissionError),
			expErr:      fmt.Errorf(client.VaultPermissionError),
		},
		{
			name:            "should error when batch_results is not []interface{}",
			decryptedSecret: nil,
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data: map[string]interface{}{
					"batch_results": map[string]interface{}{
						"plaintext": decryptedString,
					},
				},
				Warnings: nil,
				Auth:     nil,
				WrapInfo: nil,
			},
			keyRefs:     []string{"keyRef1"},
			shouldRenew: false,
			returnErr:   nil,
			expErr:      fmt.Errorf("batch results casting error"),
		},
		{
			name:            "should error when batch_results entries are not map[string]interface{}",
			decryptedSecret: nil,
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data: map[string]interface{}{
					"batch_results": []interface{}{"plaintext"},
				},
				Warnings: nil,
				Auth:     nil,
				WrapInfo: nil,
			},
			keyRefs:     []string{"keyRef1"},
			shouldRenew: false,
			returnErr:   nil,
			expErr:      fmt.Errorf("batch result casting error"),
		},
		{
			name:            "should error when not base64 encoded",
			decryptedSecret: nil,
			vaultResponse: &vaultapi.Secret{
				RequestID:     "",
				LeaseID:       "",
				LeaseDuration: 0,
				Renewable:     false,
				Data: map[string]interface{}{
					"batch_results": []interface{}{
						map[string]interface{}{
							"plaintext": decryptedString,
						},
					},
				},
				Warnings: nil,
				Auth:     nil,
				WrapInfo: nil,
			},
			keyRefs:     []string{"keyRef1"},
			shouldRenew: false,
			returnErr:   nil,
			expErr:      base64.CorruptInputError(6),
		},
		{
			name:            "should error when Vault returns an empty response",
			decryptedSecret: nil,
			vaultResponse:   nil,
			keyRefs:         []string{"keyRef1"},
			shouldRenew:     false,
			returnErr:       nil,
			expErr:          fmt.Errorf("tried to decrypt keyReference: keyRef1 but vault returned an empty body"),
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
				if renewed {
					tt.returnErr = nil
				}
				return tt.vaultResponse, tt.returnErr
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			assertThat := assert.New(t)
			v := NewVaultDecrypter(mockClient)
			decryptedSecret, err := v.Decrypt(tt.keyRefs, encryptedData, ctx)

			assertThat.Equal(tt.decryptedSecret, decryptedSecret)
			assertThat.Equal(tt.shouldRenew, renewed)
			assertThat.Equal(tt.expErr, err)
		})
	}
}
