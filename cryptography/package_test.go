package cryptography_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/cryptography"
	"github.com/stretchr/testify/assert"
)

func TestPackageEncrypt(t *testing.T) {
	ctx := context.Background()

	cipherText, err := cryptography.Encrypt(ctx, "test_plain_str")
	assert.Nil(t, err)

	plainText, err := cryptography.Decrypt(ctx, cipherText)
	assert.Nil(t, err)
	assert.Equal(t, "test_plain_str", plainText)
}
