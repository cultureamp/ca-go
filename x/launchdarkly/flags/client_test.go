package flags

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validConfigJSON = `
{
    "sdkKey":"super-secret-key",
    "options":{
        "daemonMode":{
            "dynamo_base_url":"url-here",
            "DynamoTableName":"my-dynamo-table",
            "dynamoCacheTTLSeconds":30
        },
        "proxyMode":{
            "url":"https://relay-proxy.cultureamp.net"
        }
    }
}
`

func TestInitialisationClient(t *testing.T) {
	t.Run("configures for Test mode if LAUNCHDARKLY_CONFIGURATION is not set", func(t *testing.T) {
		c, err := NewClient()
		require.NoError(t, err)

		require.NoError(t, c.Connect())

		td, err := c.TestDataSource()
		require.NoError(t, err)

		td.Update(td.Flag("test-flag").VariationForAllUsers(true))
		res, err := c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewAnonymousUser(""), false)
		require.NoError(t, err)
		assert.Equal(t, true, res)

		td.Update(td.Flag("test-flag").VariationForAllUsers(false))
		res, err = c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewAnonymousUser(""), true)
		require.NoError(t, err)
		assert.Equal(t, false, res)
	})

	t.Run("configures for Test mode when explicitly told to", func(t *testing.T) {
		c, err := NewClient(WithTestMode(nil))
		require.NoError(t, err)

		require.NoError(t, c.Connect())

		td, err := c.TestDataSource()
		require.NoError(t, err)

		td.Update(td.Flag("test-flag").VariationForAllUsers(true))
		res, err := c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewAnonymousUser(""), false)
		require.NoError(t, err)
		assert.Equal(t, true, res)

		td.Update(td.Flag("test-flag").VariationForAllUsers(false))
		res, err = c.QueryBoolWithEvaluationContext("test-flag", evaluationcontext.NewAnonymousUser(""), true)
		require.NoError(t, err)
		assert.Equal(t, false, res)
	})

	t.Run("configures for Test mode with a JSON file data source", func(t *testing.T) {
		jsonFilename, err := ioutil.TempFile("", "test-flags.json")
		require.NoError(t, err)

		_, err = jsonFilename.Write([]byte(`
		{
			"flagValues": {
			  "my-string-flag-key": "value-1",
			  "my-boolean-flag-key": true,
			  "my-integer-flag-key": 3
			}
		}	
		`))
		require.NoError(t, err)

		c, err := NewClient(WithTestMode(&TestModeConfig{FlagFilename: jsonFilename.Name()}))
		require.NoError(t, err)

		require.NoError(t, c.Connect())

		res, err := c.QueryStringWithEvaluationContext("my-string-flag-key", evaluationcontext.NewAnonymousUser(""), "value-2")
		require.NoError(t, err)
		assert.Equal(t, "value-1", res)

		res2, err := c.QueryBoolWithEvaluationContext("my-boolean-flag-key", evaluationcontext.NewAnonymousUser(""), false)
		require.NoError(t, err)
		assert.Equal(t, true, res2)

		res3, err := c.QueryIntWithEvaluationContext("my-integer-flag-key", evaluationcontext.NewAnonymousUser(""), 1)
		require.NoError(t, err)
		assert.Equal(t, 3, res3)
	})

	t.Run("allows an initialisation wait time to be specified", func(t *testing.T) {
		os.Setenv(configurationEnvVar, validConfigJSON)
		defer os.Unsetenv(configurationEnvVar)

		client, err := NewClient(
			WithInitWait(2 * time.Second))
		require.NoError(t, err)
		assert.Equal(t, client.initWait, 2*time.Second)
	})

	t.Run("configures for Lambda (daemon) mode", func(t *testing.T) {
		os.Setenv(configurationEnvVar, validConfigJSON)
		defer os.Unsetenv(configurationEnvVar)

		client, err := NewClient(WithLambdaMode(nil))
		require.NoError(t, err)

		err = client.Connect()
		require.NoError(t, err)

		assert.True(t, client.wrappedClient.GetDataStoreStatusProvider().GetStatus().Available)
	})

	t.Run("configures for Lambda mode with optional overrides", func(t *testing.T) {
		os.Setenv(configurationEnvVar, validConfigJSON)
		defer os.Unsetenv(configurationEnvVar)

		client, err := NewClient(WithLambdaMode(&LambdaModeConfig{
			DynamoCacheTTL: 10 * time.Second,
			DynamoBaseURL:  "https://dynamo.us-east-1.amazonaws.com",
		}))
		require.NoError(t, err)

		assert.Equal(t, 10*time.Second, client.lambdaModeConfig.DynamoCacheTTL)
		assert.Equal(t, "https://dynamo.us-east-1.amazonaws.com", client.lambdaModeConfig.DynamoBaseURL)

		err = client.Connect()
		require.NoError(t, err)

		assert.True(t, client.wrappedClient.GetDataStoreStatusProvider().GetStatus().Available)
	})

	t.Run("configures for Proxy mode", func(t *testing.T) {
		os.Setenv(configurationEnvVar, validConfigJSON)
		defer os.Unsetenv(configurationEnvVar)
		client, err := NewClient()
		require.NoError(t, err)

		assert.Equal(t, "https://relay-proxy.cultureamp.net", client.wrappedConfig.ServiceEndpoints.Streaming)
		assert.Equal(t, "https://relay-proxy.cultureamp.net", client.wrappedConfig.ServiceEndpoints.Events)
		assert.Equal(t, "https://relay-proxy.cultureamp.net", client.wrappedConfig.ServiceEndpoints.Polling)
	})

	t.Run("configures for Proxy mode with optional overrides", func(t *testing.T) {
		os.Setenv(configurationEnvVar, validConfigJSON)
		defer os.Unsetenv(configurationEnvVar)
		client, err := NewClient(WithProxyMode(&ProxyModeConfig{
			RelayProxyURL: "https://foo.bar",
		}))
		require.NoError(t, err)

		assert.Equal(t, "https://foo.bar", client.wrappedConfig.ServiceEndpoints.Streaming)
		assert.Equal(t, "https://foo.bar", client.wrappedConfig.ServiceEndpoints.Events)
		assert.Equal(t, "https://foo.bar", client.wrappedConfig.ServiceEndpoints.Polling)
	})
}
