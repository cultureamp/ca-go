package cryptography_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/cryptography"
)

func TestPackageEncrypt(t *testing.T) {
	ctx := context.Background()

	cryptography.Encrypt(ctx, "test_plain_str")
}
