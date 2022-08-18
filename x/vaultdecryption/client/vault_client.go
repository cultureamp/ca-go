package client

import (
	"context"
	"fmt"

	"github.com/cultureamp/ca-go/x/vaultdecryption/auth"
	"github.com/cultureamp/glamplify/log"
	vaultapi "github.com/hashicorp/vault/api"
)

const (
	vaultDecrypterRole   = "decrypter"
	VaultPermissionError = "Code: 403"
)

var (
	Create    = NewVaultClient
	Login     = wrappedLogin
	NewClient = wrappedNewVClient
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
	client, err := NewClient(settings.VaultAddr)
	if err != nil {
		return nil, err
	}

	secret, err := Login(client, ctx, auth.NewAWSIamAuth(vaultDecrypterRole, decrypterRoleArn))
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return &VaultClient{settings, client}, nil
}

func wrappedLogin(client *vaultapi.Client, ctx context.Context, authMethod vaultapi.AuthMethod) (*vaultapi.Secret, error) {
	return client.Auth().Login(ctx, authMethod)
}

func wrappedNewVClient(vaultAddr string) (*vaultapi.Client, error) {
	return vaultapi.NewClient(&vaultapi.Config{
		Address: vaultAddr,
	})
}

func (v *VaultClient) RenewClient(ctx context.Context) error {
	logger := log.NewFromCtx(ctx)
	newClient, err := Create(v.settings, ctx)
	if err != nil {
		logger.Info("unable to renew vault client", log.Fields{"err": err.Error()})
		return err
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
