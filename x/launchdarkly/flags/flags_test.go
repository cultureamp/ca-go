package flags

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingletonInitialisation(t *testing.T) {
	t.Run("does not error if SDK key supplied as env var", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)
		err := Configure(WithTestMode(nil))
		require.NoError(t, err)
	})
}
