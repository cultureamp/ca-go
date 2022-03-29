package flags_test

import (
	"os"
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags"
	"github.com/stretchr/testify/require"
)

func TestSingletonInitialisation(t *testing.T) {
	t.Run("errors if LAUNCHDARKLY_CONFIGURATION is not present in environment", func(t *testing.T) {
		err := flags.Configure()
		require.Error(t, err)

		_, err = flags.GetDefaultClient()
		require.Error(t, err)
	})

	t.Run("does not error if SDK key supplied as env var", func(t *testing.T) {
		os.Setenv("LAUNCHDARKLY_CONFIGURATION", "{\"sdkKey\": \"abc\"}")
		defer os.Unsetenv("LAUNCHDARKLY_CONFIGURATION")
		err := flags.Configure()
		require.NoError(t, err)
	})
}
