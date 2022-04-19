package flags

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingletonInitialisation(t *testing.T) {
	t.Run(fmt.Sprintf("errors if %s is not present in environment", configurationEnvVar), func(t *testing.T) {
		err := Configure()
		require.Error(t, err)

		_, err = GetDefaultClient()
		require.Error(t, err)
	})

	t.Run("does not error if SDK key supplied as env var", func(t *testing.T) {
		os.Setenv(configurationEnvVar, validConfigJSON)
		defer os.Unsetenv(configurationEnvVar)
		err := Configure()
		require.NoError(t, err)
	})
}
