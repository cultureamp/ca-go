package flags

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	lddynamodb "github.com/launchdarkly/go-server-sdk-dynamodb"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

var errClientNotConfigured = errors.New("client not configured")

const configurationEnvVar = "LAUNCHDARKLY_CONFIGURATION"

// proxyModeConfig declares the structure of the Proxy Option in configurationJSON
type proxyModeConfig struct {
	RelayProxyURL string `json:"url"`
}

// daemonModeConfig declares the structure of the DaemonMode Option in configurationJSON
type daemonModeConfig struct {
	DynamoTableName string `json:"DynamoTableName"`
	DynamoBaseURL   string `json:"DynamoBaseUrl"`
	CacheTTLSeconds int64  `json:"dynamoCacheTTLSeconds"`
}

// configurationJSON declares the structure of the LAUNCHDARKLY_CONFIGURATION
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

// WithLambdaMode configures the client to connect to Dynamo for feature flags
func WithLambdaMode() ConfigOption {
	return func(c *Client) {
		c.mode = modeLambda
	}
}

func configFromEnvironment() (parsedConfig configurationJSON, err error) {
	configEnvVar, ok := os.LookupEnv(configurationEnvVar)
	if !ok {
		return parsedConfig, fmt.Errorf("the %s environment variable was not found", configurationEnvVar)
	}

	if err := json.Unmarshal([]byte(configEnvVar), &parsedConfig); err != nil {
		return parsedConfig, fmt.Errorf("parse %s: %w", configurationEnvVar, err)
	}

	// At a minimum the JSON should have an SDK key.
	if parsedConfig.SDKKey == "" {
		return parsedConfig, fmt.Errorf("%s did not contain an SDK key", configurationEnvVar)
	}

	return parsedConfig, nil
}

func configForProxyMode(cfg *proxyModeConfig) ld.Config {
	return ld.Config{
		ServiceEndpoints: ldcomponents.RelayProxyEndpoints(cfg.RelayProxyURL),
	}
}

func configForLambdaMode(cfg *daemonModeConfig) ld.Config {
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
