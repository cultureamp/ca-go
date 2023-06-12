package flags

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	ld "github.com/launchdarkly/go-server-sdk/v6"
	"github.com/launchdarkly/go-server-sdk/v6/testhelpers/ldtestdata"
)

// Client is a wrapper around the LaunchDarkly client.
type Client struct {
	sdkKey             string
	initWait           time.Duration
	mode               mode
	bigSegmentsEnabled bool
	wrappedConfig      ld.Config
	wrappedClient      *ld.LDClient

	testModeConfig *TestModeConfig

	// Optional config overrides.
	proxyModeConfig  *ProxyModeConfig
	lambdaModeConfig *LambdaModeConfig
}

// The mode the SDK should be configured for.
type mode int

const (
	modeProxy  mode = iota // proxies requests through the LD Relay.
	modeLambda             // connects directly to DynamoDB.
	modeTest               // allows test data to be supplied.
)

// NewClient configures and returns an instance of the client. The client is
// configured automatically from the LAUNCHDARKLY_CONFIGURATION environment
// variable if it exists. Otherwise, the client falls back to test mode. See
// launchdarkly/flags/doc.go for more information.
func NewClient(opts ...ConfigOption) (*Client, error) {
	c := &Client{
		initWait:           5 * time.Second, // wait up to 5 seconds for LD to connect.
		mode:               modeProxy,       // defaults to proxying requests through the LD Relay.
		bigSegmentsEnabled: true,            // defaults to enable big segments
	}

	for _, opt := range opts {
		opt(c)
	}

	parsedConfig := configurationJSON{}
	_, ok := os.LookupEnv(configurationEnvVar)
	if ok {
		config, err := configFromEnvironment()
		if err != nil {
			return nil, fmt.Errorf("configure from environment variable: %w", err)
		}
		parsedConfig = config
	}

	// Use test mode if LAUNCHDARKLY_CONFIGURATION isn't set OR if the user
	// explicitly configured the client for test mode.
	if !ok || c.mode == modeTest {
		c.mode = modeTest
		if c.testModeConfig == nil {
			c.testModeConfig = &TestModeConfig{}
		}
		c.wrappedConfig = configForTestMode(c.testModeConfig)

		// Short-circuit the rest of the configuration.
		return c, nil
	}

	c.sdkKey = parsedConfig.SDKKey

	if parsedConfig.Proxy != nil && c.mode == modeProxy {
		c.wrappedConfig = configForProxyMode(parsedConfig, c.proxyModeConfig)
	}

	if parsedConfig.Storage != nil && c.mode == modeLambda {
		c.wrappedConfig = configForLambdaMode(parsedConfig, c.lambdaModeConfig)
	}

	// Configure big segments if the storage table name is present
	if c.bigSegmentsEnabled &&
		parsedConfig.Storage != nil &&
		parsedConfig.Storage.TableName != "" {
		c.wrappedConfig.BigSegments = configForBigSegments(parsedConfig).BigSegments
	}

	return c, nil
}

// Connect attempts to establish the initial connection to LaunchDarkly. An
// error is returned if a connection has already been established, or a
// connection error occurs.
func (c *Client) Connect() error {
	if c.wrappedClient != nil {
		return errors.New("attempted to call Connect on a connected client")
	}

	wrappedClient, err := ld.MakeCustomClient(c.sdkKey, c.wrappedConfig, c.initWait)
	if err != nil {
		return fmt.Errorf("create LaunchDarkly client: %w", err)
	}

	c.wrappedClient = wrappedClient

	return nil
}

// QueryBool retrieves the value of a boolean flag. User attributes are
// extracted from the context. The supplied fallback value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryBool(ctx context.Context, key FlagName, fallbackValue bool) (bool, error) {
	user, err := evaluationcontext.EvaluationContextFromContext(ctx)
	if err != nil {
		return fallbackValue, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.BoolVariation(string(key), user.ToLDContext(), fallbackValue)
}

// QueryBoolWithEvaluationContext retrieves the value of a boolean flag. An evaluation context
// must be supplied manually. The supplied fallback value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryBoolWithEvaluationContext(key FlagName, evalContext evaluationcontext.Context, fallbackValue bool) (bool, error) {
	return c.wrappedClient.BoolVariation(string(key), evalContext.ToLDContext(), fallbackValue)
}

// QueryString retrieves the value of a string flag. User attributes are
// extracted from the context. The supplied fallback value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryString(ctx context.Context, key FlagName, fallbackValue string) (string, error) {
	user, err := evaluationcontext.EvaluationContextFromContext(ctx)
	if err != nil {
		return fallbackValue, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.StringVariation(string(key), user.ToLDContext(), fallbackValue)
}

// QueryStringWithEvaluationContext retrieves the value of a string flag. An evaluation context
// must be supplied manually. The supplied fallback value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryStringWithEvaluationContext(key FlagName, evalContext evaluationcontext.Context, fallbackValue string) (string, error) {
	return c.wrappedClient.StringVariation(string(key), evalContext.ToLDContext(), fallbackValue)
}

// QueryInt retrieves the value of an integer flag. User attributes are
// extracted from the context. The supplied fallback value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryInt(ctx context.Context, key FlagName, fallbackValue int) (int, error) {
	user, err := evaluationcontext.EvaluationContextFromContext(ctx)
	if err != nil {
		return fallbackValue, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.IntVariation(string(key), user.ToLDContext(), fallbackValue)
}

// QueryIntWithEvaluationContext retrieves the value of an integer flag. An evaluation context
// must be supplied manually. The supplied fallback value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryIntWithEvaluationContext(key FlagName, evalContext evaluationcontext.Context, fallbackValue int) (int, error) {
	return c.wrappedClient.IntVariation(string(key), evalContext.ToLDContext(), fallbackValue)
}

// RawClient returns the wrapped LaunchDarkly client. The return value should be
// casted to an *ld.LDClient instance.
func (c *Client) RawClient() interface{} {
	return c.wrappedClient
}

// Shutdown instructs the wrapped LaunchDarkly client to close any open
// connections and flush any flag evaluation events.
func (c *Client) Shutdown() error {
	return c.wrappedClient.Close()
}

// TestDataSource returns the dynamic test data source used by the client, or an
// error if:
// - the client wasn't configured in test mode.
// - the client was configured to read test data from a JSON file.
//
// See https://docs.launchdarkly.com/sdk/features/test-data-sources for more
// information on using the test data source returned by this method.
func (c *Client) TestDataSource() (*ldtestdata.TestDataSource, error) {
	if c.testModeConfig == nil || c.testModeConfig.datasource == nil {
		return nil, errors.New("client not initialised with dynamic test data source")
	}

	return c.testModeConfig.datasource, nil
}
