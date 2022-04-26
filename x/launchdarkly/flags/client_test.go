package flags

import (
	"os"
	"testing"
	"time"

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
	t.Run("errors if an SDK key is not supplied", func(t *testing.T) {
		_, err := NewClient()
		require.Error(t, err)
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
