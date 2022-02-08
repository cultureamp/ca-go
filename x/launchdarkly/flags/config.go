package flags

import (
	"errors"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	lddynamodb "github.com/launchdarkly/go-server-sdk-dynamodb"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

var errClientNotConfigured = errors.New("client not configured")

type proxyModeConfig struct {
	relayProxyURL string
}

type daemonModeConfig struct {
	dynamoTableName string
	dynamoBaseURL   string
	cacheTTL        time.Duration
}

// ConfigOption are functions that can be supplied to Configure and NewClient to
// configure the flags client.
type ConfigOption func(c *Client)

// WithSDKKey configures the client to use the given SDK key to authenticate
// against LaunchDarkly.
func WithSDKKey(key string) ConfigOption {
	return func(c *Client) {
		c.sdkKey = key
	}
}

// WithInitWait configures the client to wait for the given duration for the
// LaunchDarkly client to connect.
// If you don't provide this option, the client will wait up to 5 seconds by
// default.
func WithInitWait(t time.Duration) ConfigOption {
	return func(c *Client) {
		c.initWait = t
	}
}

// WithProxyMode configures the client to establish a connection to LaunchDarkly
// via a Relay Proxy.
func WithProxyMode(proxyURL *url.URL) ConfigOption {
	return func(c *Client) {
		c.proxyModeConfig = &proxyModeConfig{
			relayProxyURL: proxyURL.String(),
		}
	}
}

// WithDaemonMode configures the client to source flag data from a DynamoDB
// table, bypassing a direct connection to LaunchDarkly completely.
func WithDaemonMode(dynamoTableName string, cacheTTL time.Duration) ConfigOption {
	return func(c *Client) {
		if c.daemonModeConfig == nil {
			c.daemonModeConfig = &daemonModeConfig{}
		}

		c.daemonModeConfig = &daemonModeConfig{
			dynamoTableName: dynamoTableName,
			cacheTTL:        cacheTTL,
		}
	}
}

// WithDynamoBaseURL configures the client to use the given base URL for
// DyanmoDB, overriding any AWS configuration implicit in the environment.
//
// This will typically only be used in local development or testing, where you
// might supply the URL of a local DynamoDB instance.
func WithDynamoBaseURL(baseURL *url.URL) ConfigOption {
	return func(c *Client) {
		if c.daemonModeConfig == nil {
			c.daemonModeConfig = &daemonModeConfig{}
		}

		c.daemonModeConfig.dynamoBaseURL = baseURL.String()
	}
}

func configForProxyMode(cfg *proxyModeConfig) ld.Config {
	return ld.Config{
		DataSource: ldcomponents.StreamingDataSource().BaseURI(cfg.relayProxyURL),
	}
}

func configForDaemonMode(cfg *daemonModeConfig) ld.Config {
	datastoreBuilder := lddynamodb.DataStore(cfg.dynamoTableName)

	if cfg.dynamoBaseURL != "" {
		datastoreBuilder.ClientConfig(aws.NewConfig().WithEndpoint(cfg.dynamoBaseURL))
	}

	return ld.Config{
		DataSource: ldcomponents.ExternalUpdatesOnly(),
		DataStore: ldcomponents.PersistentDataStore(
			datastoreBuilder,
		).CacheTime(cfg.cacheTTL),
	}
}
