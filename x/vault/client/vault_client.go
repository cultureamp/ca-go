package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/cultureamp/ca-go/x/log"
	"github.com/cultureamp/ca-go/x/vault/auth"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/benchhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	vaultDecrypterRole   = "decrypter"
	VaultPermissionError = "Code: 403"
	EncryptionAction     = "encrypt"
	DecryptionAction     = "decrypt"
)

var (
	Create    = NewVaultClient
	Login     = wrappedLogin
	NewClient = wrappedNewVClient
)

type VaultSettings struct {
	RoleArn   string
	VaultAddr string
}

type VaultClient struct {
	settings *VaultSettings
	client   *vaultapi.Client
}

// NewTestingClient creates a test vault cluster and returns a configured API
// client and closer function.
func NewTestingClient(t testing.TB) (*VaultClient, func(), error) {
	t.Helper()
	client, closer, err := testVaultServerCoreConfig(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
	})
	return &VaultClient{nil, client}, closer, err

}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper.
func testVaultServerCoreConfig(t testing.TB, coreConfig *vault.CoreConfig) (*vaultapi.Client, func(), error) {
	t.Helper()

	cluster := vault.NewTestCluster(benchhelpers.TBtoT(t), coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(benchhelpers.TBtoT(t), core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	err := client.Sys().Mount("transit", &vaultapi.MountInput{
		Type: "transit",
	})
	if err != nil {
		return nil, nil, err
	}

	return client, func() { defer cluster.Cleanup() }, nil
}

func NewVaultClient(settings *VaultSettings, ctx context.Context) (*VaultClient, error) {
	logger := log.NewFromCtx(ctx)
	decrypterRoleArn := settings.RoleArn
	if decrypterRoleArn == "" || settings.VaultAddr == "" {
		err := fmt.Errorf("VaultClient settings incomplete, must provide RoleArn and VaultAddr")
		logger.Error().Err(err).Msgf("VaultClient settings incomplete: %+v", settings)
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
		err = fmt.Errorf("login auth error")
		logger.Error().Err(err).Msg("no auth info was returned after login")
		return nil, err
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
		logger.Error().Err(err).Msg("unable to renew vault client")
		return err
	}
	v.client = newClient.client
	logger.Info().Msg("Renewed vault client")
	return nil
}

func (v VaultClient) GetSecret(batch []interface{}, keyReference string, action string) (*vaultapi.Secret, error) {
	secret, err := v.client.Logical().Write(fmt.Sprintf("transit/%s/%s", action, keyReference), map[string]interface{}{
		"batch_input": batch,
	})
	if err != nil {
		return nil, err
	}
	return secret, nil
}
