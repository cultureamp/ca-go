package vaultdecryption

import (
	"context"
	"fmt"

	"github.com/cultureamp/glamplify/log"
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

type VaultClient struct {
	settings *VaultSettings
	client   *vaultapi.Client
}

func NewVaultClient(settings *VaultSettings, ctx context.Context) (*VaultClient, error) {
	logger := log.NewFromCtx(ctx)
	decrypterRoleArn := settings.DecrypterRoleArn
	if decrypterRoleArn == "" || settings.VaultAddr == "" {
		err := fmt.Errorf("VaultClient settings incomplete, must provide DecrypterRoleArn and VaultAddr")
		logger.Error("VaultClient settings incomplete", err, log.Fields{"VaultSettings": settings})
		return nil, err
	}
	client, err := newVaultClient(NewAWSIamAuth(vaultDecrypterRole, decrypterRoleArn), ctx, *settings)
	if err != nil {
		logger.Error("client could not be initialised", err)
		return nil, err
	}
	return &VaultClient{settings, client}, nil
}

func newVaultClient(authMethod vaultapi.AuthMethod, ctx context.Context, settings VaultSettings) (*vaultapi.Client, error) {
	client, err := vaultapi.NewClient(&vaultapi.Config{
		Address: settings.VaultAddr,
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

func (v VaultClient) RenewClient(ctx context.Context) error {
	newClient, err := NewVaultClient(v.settings, ctx)
	logger := log.NewFromCtx(ctx)
	if err != nil {
		logger.Info("unable to renew vault client", log.Fields{"err": err.Error()})
		return fmt.Errorf("unable to renew vault client")
	}
	v.client = newClient.client
	logger.Info("Renewed vault client")
	return nil
}

func (v VaultClient) GetSecret(batch []interface{}, keyReference string) (*vaultapi.Secret, error) {
	secret, err := v.client.Logical().Write(fmt.Sprintf("transit/decryptByKey/%s", keyReference), map[string]interface{}{
		"batch_input": batch,
	})
	if err != nil {
		return nil, err
	}
	return secret, nil
}
