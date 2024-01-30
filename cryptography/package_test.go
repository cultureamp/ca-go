package cryptography_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/cryptography"
	"github.com/stretchr/testify/assert"
)

func TestPackageEncrypt(t *testing.T) {
	ctx := context.Background()
	keyId := "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"

	cipherText, err := cryptography.Encrypt(ctx, keyId, "test_plain_str")
	assert.Nil(t, err)

	plainText, err := cryptography.Decrypt(ctx, keyId, cipherText)
	assert.Nil(t, err)
	assert.Equal(t, "test_plain_str", plainText)
}
