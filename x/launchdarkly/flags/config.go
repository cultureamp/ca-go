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

// configurationJSON declares the structure of the LAUNCHDARKLY_CONFIGURATION
// environment variable.
type configurationJSON struct {
	SDKKey  string `json:"sdkKey"`
	Options struct {
		DaemonMode *struct {
			DynamoTableName string `json:"DynamoTableName"`
		} `json:"daemonMode"`
		Proxy *struct {
			RelayProxyURL string `json:"url"`
		} `json:"proxyMode"`
	} `json:"options"`
}

// ProxyModeConfig declares optional overrides for configuring the client
// in Proxy mode.
type ProxyModeConfig struct {
	RelayProxyURL string
}

// LambdaModeConfig declares optional overrides for configuring the client
// in Lambda mode.
type LambdaModeConfig struct {
	DynamoCacheTTL time.Duration
	DynamoBaseURL  string
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

// WithLambdaMode configures the client to connect to Dynamo for flags.
func WithLambdaMode(cfg *LambdaModeConfig) ConfigOption {
	return func(c *Client) {
		c.mode = modeLambda
		c.lambdaModeConfig = cfg
	}
}

// WithProxyMode configures the client to connect to LaunchDarkly via the
// Relay Proxy. This is typically set automatically based on the LAUNCHDARKLY_CONFIGURATION
// environment variable. Only use this ConfigOption if you need to override
// the URL of the Relay Proxy to connect to.
func WithProxyMode(cfg *ProxyModeConfig) ConfigOption {
	return func(c *Client) {
		c.mode = modeProxy
		c.proxyModeConfig = cfg
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

func configForProxyMode(env configurationJSON, cfg *ProxyModeConfig) ld.Config {
	urlToUse := env.Options.Proxy.RelayProxyURL
	// Override the Relay URL from the environment variable if one was provided
	// explicitly.
	if cfg != nil && cfg.RelayProxyURL != "" {
		urlToUse = cfg.RelayProxyURL
	}

	return ld.Config{
		ServiceEndpoints: ldcomponents.RelayProxyEndpoints(urlToUse),
	}
}

func configForLambdaMode(env configurationJSON, cfg *LambdaModeConfig) ld.Config {
	datastoreBuilder := lddynamodb.DataStore(env.Options.DaemonMode.DynamoTableName)

	// Set the Dynamo base URL if one was provided explicitly.
	if cfg != nil && cfg.DynamoBaseURL != "" {
		datastoreBuilder.ClientConfig(aws.NewConfig().WithEndpoint(cfg.DynamoBaseURL))
	}

	datastore := ldcomponents.PersistentDataStore(
		datastoreBuilder,
	)

	// Override the default cache TTL if one was provided explicitly.
	if cfg != nil && cfg.DynamoCacheTTL != 0 {
		datastore.CacheTime(cfg.DynamoCacheTTL)
	}

	return ld.Config{
		DataSource: ldcomponents.ExternalUpdatesOnly(),
		DataStore:  datastore,
	}
}
