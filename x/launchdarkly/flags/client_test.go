package flags_test

import (
	"os"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags"
	"github.com/stretchr/testify/require"
)

func TestClientInitialisation(t *testing.T) {
	t.Run("errors if no env var is supplied", func(t *testing.T) {
		_, err := flags.NewClient()
		require.Error(t, err)
	})

	t.Run("allows an initialisation wait time to be specified", func(t *testing.T) {
		os.Setenv("LAUNCHDARKLY_CONFIGURATION", "{\"sdkKey\": \"abc\"}")
		defer os.Unsetenv("LAUNCHDARKLY_CONFIGURATION")
		_, err := flags.NewClient(
			flags.WithInitWait(2 * time.Second))
		require.NoError(t, err)
	})
}
