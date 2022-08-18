package vaultdecryption

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVaultClientErrors(t *testing.T) {
	tests := []struct {
		name     string
		settings VaultSettings
	}{
		{
			"should error when DecryptionRoleArn is empty",
			VaultSettings{
				DecrypterRoleArn: "",
				VaultAddr:        "1234",
			},
		},
		{
			"should error when VaultAddr is empty",
			VaultSettings{
				DecrypterRoleArn: "arn",
				VaultAddr:        "",
			},
		},
	}
	for _, tt := range tests {
		ctx := context.Background()

		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVaultClient(&tt.settings, ctx)
			assert.NotNil(t, err)
		})
	}
}
