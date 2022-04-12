package flags

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	lddynamodb "github.com/launchdarkly/go-server-sdk-dynamodb"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

var errClientNotConfigured = errors.New("client not configured")

type proxyModeConfig struct {
	RelayProxyURL string `json:"url"`
}

type daemonModeConfig struct {
	DynamoTableName string `json:"DynamoTableName"`
	DynamoBaseURL   string `json:"DynamoBaseUrl"`
	CacheTTLSeconds int64  `json:"dynamoCacheTTLSeconds"`
}

//  The EnvConfig class holds the configuration values. It is expecting
//  attributes to be in the following structure:
//  '{
//     "sdkKey": "super-secret-key",
//     "options": {
//       "daemonMode": {
//         "dynamo_base_url": "url-here"
//         "DynamoTableName": "my-dynamo-table",
//         "dynamoCacheTTLSeconds": 30
//       },
//       "proxyMode": {
//         "url": "https://relay-proxy.cultureamp.net"
//       }
//    }'
// configurationJSON is the structure of the LAUNCHDARKLY_CONFIGURATION
// environment variable.
type configurationJSON struct {
	SDKKey  string `json:"sdkKey"`
	Options struct {
		DaemonMode *daemonModeConfig `json:"daemonMode"`
		Proxy      *proxyModeConfig  `json:"proxyMode"`
	} `json:"options"`
}

// ConfigOption are functions that can be supplied to Configure and NewClient to
// configure the flags client.
type ConfigOption func(c *Client)

// WithInitWait configures the client to wait for the given duration for the
// LaunchDarkly client to connect.
// If you don't provide this option, the client will wait up to 5 seconds by
// default.
func WithInitWait(t time.Duration) ConfigOption {
	return func(c *Client) {
		c.initWait = t
	}
}

func configForProxyMode(cfg *proxyModeConfig) ld.Config {
	return ld.Config{
		ServiceEndpoints: ldcomponents.RelayProxyEndpoints(cfg.RelayProxyURL),
	}
}

func configForDaemonMode(cfg *daemonModeConfig) ld.Config {
	datastoreBuilder := lddynamodb.DataStore(cfg.DynamoTableName)

	if cfg.DynamoBaseURL != "" {
		datastoreBuilder.ClientConfig(aws.NewConfig().WithEndpoint(cfg.DynamoBaseURL))
	}

	return ld.Config{
		DataSource: ldcomponents.ExternalUpdatesOnly(),
		DataStore: ldcomponents.PersistentDataStore(
			datastoreBuilder,
		).CacheTime(time.Duration(cfg.CacheTTLSeconds) * time.Second),
	}
}
