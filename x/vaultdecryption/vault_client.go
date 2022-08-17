package vaultdecryption

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

const (
	vaultDecrypterRole   = "decrypter"
	vaultPermissionError = "Code: 403"
	maxRetries           = 5
)

type VaultSettings struct {
	DecrypterRoleArn string
	VaultAddr        string
}

type VaultClients struct {
	settings *VaultSettings
}

func DefaultVaultClients(settings *VaultSettings) *VaultClients {
	return &VaultClients{settings}
}

func (v *VaultClients) NewAwsIamVaultDecrypterClient(ctx context.Context) (*vaultapi.Client, error) {
	decrypterRoleArn := v.settings.DecrypterRoleArn
	if decrypterRoleArn == "" {
		return nil, fmt.Errorf("decrypter role ARN is not set")
	}
	return v.newVaultClient(NewAWSIamAuth(vaultDecrypterRole, decrypterRoleArn), ctx)
}

func (v *VaultClients) newVaultClient(authMethod vaultapi.AuthMethod, ctx context.Context) (*vaultapi.Client, error) {
	vaultAddr := v.settings.VaultAddr
	if vaultAddr == "" {
		return nil, fmt.Errorf("vault address is not set")
	}
	client, err := vaultapi.NewClient(&vaultapi.Config{
		Address: vaultAddr,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	secret, err := client.Auth().Login(ctx, authMethod)
	if err != nil {
		return nil, fmt.Errorf("unable to login with auth method: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return client, nil
}
