package vaultdecryption

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	renewClient func(ctx context.Context) error
	getSecret   func(batch []interface{}, keyReference string) (*vaultapi.Secret, error)
}

func (m MockClient) RenewClient(ctx context.Context) error {
	return m.renewClient(ctx)
}
func (m MockClient) GetSecret(batch []interface{}, keyReference string) (*vaultapi.Secret, error) {
	return m.getSecret(batch, keyReference)
}

const (
	decryptedString = "abc123!?$*&()'-=@~"
)

func TestNewVaultDecrypter(t *testing.T) {
	secretReturn := vaultapi.Secret{
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
	}

	tests := []struct {
		name            string
		decryptedSecret []string
		secret          *vaultapi.Secret
		shouldRenew     bool
		err             error
	}{
		{
			"should not error when secret is returned as usual",
			[]string{decryptedString},
			&secretReturn,
			false,
			nil},
		{
			"should error when getSecret errors",
			nil,
			nil,
			false,
			fmt.Errorf("secretError"),
		},
		{
			"should renew then continue when getSecret returns permission error",
			[]string{decryptedString},
			&secretReturn,
			true,
			fmt.Errorf(vaultPermissionError),
		},
		{
			"should error when renewClient returns error",
			nil,
			nil,
			false,
			fmt.Errorf(vaultPermissionError),
		},
	}
	keyReferences := []string{"keyRef1", "keyRef2"}
	encryptedData := []string{"encrypted1", "encrypted2"}
	ctx := context.Background()

	for _, tt := range tests {
		renewed := false
		mockClient := MockClient{
			renewClient: func(ctx context.Context) error {
				if !renewed && tt.shouldRenew {
					renewed = true
				} else if !tt.shouldRenew {
					return tt.err
				}
				return nil
			},
			getSecret: func(batch []interface{}, keyReference string) (*vaultapi.Secret, error) {
				if renewed {
					tt.err = nil
				}
				return tt.secret, tt.err
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			v := NewVaultDecrypter(mockClient, &VaultSettings{
				DecrypterRoleArn: "arn:1234",
				VaultAddr:        "1234",
			})
			decryptedSecret, err := v.Decrypt(keyReferences, encryptedData, ctx)
			assert.Equal(t, tt.decryptedSecret, decryptedSecret)
			assert.Equal(t, tt.shouldRenew, renewed)
			assert.Equal(t, tt.err, err)
		})
	}
}
