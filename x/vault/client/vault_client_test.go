package client

import (
	"context"
	"fmt"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestNewVaultClient(t *testing.T) {
	settingsErr := fmt.Errorf("VaultClient settings incomplete, must provide RoleArn and VaultAddr")
	tests := []struct {
		name              string
		settings          VaultSettings
		returnedClientErr error
		returnedLoginErr  error
		expErr            error
	}{
		{
			"should not error when expected values are given",
			VaultSettings{
				RoleArn:   "arn",
				VaultAddr: "1234",
			},
			nil,
			nil,
			nil,
		},
		{
			"should error when DecryptionRoleArn is empty",
			VaultSettings{
				RoleArn:   "",
				VaultAddr: "1234",
			},
			nil,
			nil,
			settingsErr,
		},
		{
			"should error when VaultAddr is empty",
			VaultSettings{
				RoleArn:   "arn",
				VaultAddr: "",
			},
			nil,
			nil,
			settingsErr,
		},
		{
			"should error when client creator returns error",
			VaultSettings{
				RoleArn:   "arn",
				VaultAddr: "1234",
			},
			fmt.Errorf("error with client"),
			nil,
			fmt.Errorf("error with client"),
		},
		{
			"should error when login returns error",
			VaultSettings{
				RoleArn:   "arn",
				VaultAddr: "1234",
			},
			nil,
			fmt.Errorf("error with login"),
			fmt.Errorf("error with login"),
		},
		{
			"should error when login returns no secret",
			VaultSettings{
				RoleArn:   "arn",
				VaultAddr: "1234",
			},
			nil,
			nil,
			fmt.Errorf("no auth info was returned after login"),
		},
	}
	for _, tt := range tests {
		NewClient = func(vaultAddr string) (*vaultapi.Client, error) {
			client := &vaultapi.Client{}
			if tt.returnedClientErr != nil {
				client = nil
			}
			return client, tt.returnedClientErr
		}
		Login = func(client *vaultapi.Client, ctx context.Context, authMethod vaultapi.AuthMethod) (*vaultapi.Secret, error) {
			secret := &vaultapi.Secret{}
			if tt.returnedLoginErr != nil || tt.expErr != nil {
				secret = nil
			}
			return secret, tt.returnedLoginErr
		}
		ctx := context.Background()

		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(&tt.settings, ctx)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func TestNewTestingClient(t *testing.T) {
	t.Run("no error when setting up test client", func(t2 *testing.T) {
		_, closer, err := NewTestingClient(t)
		defer closer()
		assert.NoError(t2, err)
	})
}

func TestVaultClient_RenewClient(t *testing.T) {
	tests := []struct {
		name           string
		returnedClient *vaultapi.Client
		returnedErr    error
		expErr         error
	}{
		{
			"should change the client when no error occurs",
			&vaultapi.Client{},
			nil,
			nil,
		},
		{
			"should error when client errors",
			nil,
			fmt.Errorf("error with client"),
			fmt.Errorf("error with client"),
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		called := false
		Create = func(settings *VaultSettings, ctx context.Context) (*VaultClient, error) {
			called = true
			return &VaultClient{settings, tt.returnedClient}, tt.returnedErr
		}
		t.Run(tt.name, func(t *testing.T) {
			v := VaultClient{
				settings: &VaultSettings{
					RoleArn:   "arn",
					VaultAddr: "123",
				},
				client: nil,
			}
			err := v.RenewClient(ctx)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.returnedClient, v.client)
			assert.True(t, called)
		})
	}
}
