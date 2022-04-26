package flags

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingletonInitialisation(t *testing.T) {
	t.Run("does not error if SDK key supplied as env var", func(t *testing.T) {
		os.Setenv(configurationEnvVar, validConfigJSON)
		defer os.Unsetenv(configurationEnvVar)
		err := Configure()
		require.NoError(t, err)
	})
}
