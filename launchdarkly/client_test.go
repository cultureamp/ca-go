package flags

import (
	"os"
	"path"
	"testing"
	"time"

	//"github.com/cultureamp/ca-go/x/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cultureamp/ca-go/launchdarkly/evaluationcontext"
)

const validConfigJSON = `
{
	"sdkKey":"super-secret-key",
	"storage":{
		"dynamo_base_url":"url-here",
		"tableName":"my-dynamo-table",
		"dynamoCacheTTLSeconds":30
	},
	"proxy":{
		"url":"https://relay-proxy.cultureamp.net"
	}
}
`

const validFlagsJSON = `
{
	"flagValues": {
		"my-string-flag-key": "value-1",
		"my-boolean-flag-key": true,
		"my-integer-flag-key": 3
	}
}
`

func TestClientTestMode(t *testing.T) {
	t.Run("configures for Test mode if LAUNCHDARKLY_CONFIGURATION is not set", func(t *testing.T) {
		c, err := NewClient()
		require.NoError(t, err)

		require.NoError(t, c.Connect())

		td, err := c.TestDataSource()
		require.NoError(t, err)

		td.Update(td.Flag("test-flag").VariationForAll(true))
		res, err := c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewEvaluationContext(), false)
		require.NoError(t, err)
		assert.True(t, res)

		td.Update(td.Flag("test-flag").VariationForAll(false))
		res, err = c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewEvaluationContext(), true)
		require.NoError(t, err)
		assert.False(t, res)
	})

	t.Run("configures for Test mode with data set at runtime", func(t *testing.T) {
		c, err := NewClient(WithTestMode(nil))
		require.NoError(t, err)

		require.NoError(t, c.Connect())

		td, err := c.TestDataSource()
		require.NoError(t, err)

		td.Update(td.Flag("test-flag").VariationForAll(true))
		res, err := c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewEvaluationContext(), false)
		require.NoError(t, err)
		assert.True(t, res)

		td.Update(td.Flag("test-flag").VariationForAll(false))
		res, err = c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewEvaluationContext(), true)
		require.NoError(t, err)
		assert.False(t, res)
	})

	t.Run("configures for Test mode data sourced from a local JSON file", func(t *testing.T) {
		jsonFilename, err := os.CreateTemp("", "test-flags.json")
		require.NoError(t, err)

		_, err = jsonFilename.WriteString(validFlagsJSON)
		require.NoError(t, err)

		c, err := NewClient(WithTestMode(&TestModeConfig{FlagFilename: jsonFilename.Name()}))
		require.NoError(t, err)

		require.NoError(t, c.Connect())

		assertTestJSONFlags(t, c)
	})

	t.Run("configures for Test mode with data sourced from the default JSON file", func(t *testing.T) {
		testDir, err := os.Getwd()
		require.NoError(t, err)

		flagsFilename := path.Join(testDir, flagsJSONFilename)

		require.NoError(t, os.WriteFile(flagsFilename, []byte(validFlagsJSON), 0o666)) //nolint:gosec
		defer func() {
			require.NoError(t, os.Unsetenv(flagsFilename))
			os.Remove(flagsFilename)
		}()

		c, err := NewClient()
		require.NoError(t, err)

		require.NoError(t, c.Connect())

		assertTestJSONFlags(t, c)
	})

	t.Run("returns an error when getting the test data source if not configured in test mode", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)

		client, err := NewClient()
		require.NoError(t, err)

		_, err = client.TestDataSource()
		require.Error(t, err)
	})
}

func TestClientLogger(t *testing.T) {
	t.Run("Client creates logger from context when one isn't provided", func(t *testing.T) {
		c, err := NewClient()
		require.NoError(t, err)
		require.NoError(t, c.Connect())
	})
}

func assertTestJSONFlags(t *testing.T, c *Client) {
	t.Helper()

	evalContext := evaluationcontext.NewEvaluationContext()

	res, err := c.QueryStringWithEvaluationContext("my-string-flag-key", evalContext, "value-2")
	require.NoError(t, err)
	assert.Equal(t, "value-1", res)

	res2, err := c.QueryBoolWithEvaluationContext("my-boolean-flag-key", evalContext, false)
	require.NoError(t, err)
	assert.True(t, res2)

	res3, err := c.QueryIntWithEvaluationContext("my-integer-flag-key", evalContext, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, res3)
}

func TestClientLambdaMode(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "fake-access-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")

	t.Run("configures for Lambda (daemon) mode", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)

		client, err := NewClient(WithLambdaMode(nil))
		require.NoError(t, err)

		err = client.Connect()
		require.NoError(t, err)

		assert.True(t, client.wrappedClient.GetDataStoreStatusProvider().GetStatus().Available)
	})

	t.Run("configures for Lambda mode with optional overrides", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)

		client, err := NewClient(WithLambdaMode(&LambdaModeConfig{
			DynamoCacheTTL: 10 * time.Second,
		}))
		require.NoError(t, err)

		assert.Equal(t, 10*time.Second, client.lambdaModeConfig.DynamoCacheTTL)

		err = client.Connect()
		require.NoError(t, err)

		assert.True(t, client.wrappedClient.GetDataStoreStatusProvider().GetStatus().Available)
	})
}

func TestClientInitialisation(t *testing.T) {
	t.Run("allows an initialisation wait time to be specified", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)

		client, err := NewClient(
			WithInitWait(2 * time.Second))
		require.NoError(t, err)
		assert.Equal(t, 2*time.Second, client.initWait)
	})

	t.Run("configures for Proxy mode", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)
		client, err := NewClient()
		require.NoError(t, err)

		assert.Equal(t, "https://relay-proxy.cultureamp.net", client.wrappedConfig.ServiceEndpoints.Streaming)
		assert.Equal(t, "https://relay-proxy.cultureamp.net", client.wrappedConfig.ServiceEndpoints.Events)
		assert.Equal(t, "https://relay-proxy.cultureamp.net", client.wrappedConfig.ServiceEndpoints.Polling)
	})

	t.Run("configures for Proxy mode with optional overrides", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)
		client, err := NewClient(WithProxyMode(&ProxyModeConfig{
			RelayProxyURL: "https://foo.bar",
		}))
		require.NoError(t, err)

		assert.Equal(t, "https://foo.bar", client.wrappedConfig.ServiceEndpoints.Streaming)
		assert.Equal(t, "https://foo.bar", client.wrappedConfig.ServiceEndpoints.Events)
		assert.Equal(t, "https://foo.bar", client.wrappedConfig.ServiceEndpoints.Polling)
	})

	t.Run("allows big segments to be disabled", func(t *testing.T) {
		t.Setenv(configurationEnvVar, validConfigJSON)

		client, err := NewClient(
			WithBigSegmentsDisabled())
		require.NoError(t, err)
		assert.False(t, client.bigSegmentsEnabled)
	})
}
