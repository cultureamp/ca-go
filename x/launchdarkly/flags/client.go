package flags

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
)

// Client is a wrapper around the LaunchDarkly client.
type Client struct {
	sdkKey           string
	initWait         time.Duration
	proxyModeConfig  *proxyModeConfig
	daemonModeConfig *daemonModeConfig
	wrappedConfig    ld.Config
	wrappedClient    *ld.LDClient
}

// NewClient configures and returns an instance of the client. An error is
// returned if mandatory ConfigOptions are not supplied, or an invalid
// combination of options is provided.
func NewClient(opts ...ConfigOption) (*Client, error) {
	c := &Client{
		initWait: 5 * time.Second, // wait up to 5 seconds for LD to connect
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.sdkKey == "" {
		defaultSDKKey, ok := os.LookupEnv(defaultSDKKeyEnvironmentVariable)
		if !ok {
			return nil, errors.New("LaunchDarkly SDK key not supplied via config option and the LAUNCHDARKLY_SDK_KEY environment variable does not exist")
		}
		c.sdkKey = defaultSDKKey
	}

	if c.proxyModeConfig != nil && c.daemonModeConfig != nil {
		return nil, errors.New("cannot configure the SDK for Proxy and Daemon modes simultaneously")
	}

	if c.proxyModeConfig != nil {
		c.wrappedConfig = configForProxyMode(c.proxyModeConfig)
	} else if c.daemonModeConfig != nil {
		c.wrappedConfig = configForDaemonMode(c.daemonModeConfig)
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

	flagsClient.wrappedClient = wrappedClient

	return nil
}

// QueryBool retrieves the value of a boolean flag. User attributes are
// extracted from the context. The supplied fallback value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryBool(ctx context.Context, key FlagName, fallbackValue bool) (bool, error) {
	user, err := evaluationcontext.UserFromContext(ctx)
	if err != nil {
		return fallbackValue, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.BoolVariation(string(key), user.ToLDUser(), fallbackValue)
}

// QueryBoolWithEvaluationContext retrieves the value of a boolean flag. An evaluation context
// must be supplied manually. The supplied fallback value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryBoolWithEvaluationContext(key FlagName, evalContext evaluationcontext.Context, fallbackValue bool) (bool, error) {
	return c.wrappedClient.BoolVariation(string(key), evalContext.ToLDUser(), fallbackValue)
}

// QueryString retrieves the value of a string flag. User attributes are
// extracted from the context. The supplied fallback value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryString(ctx context.Context, key FlagName, fallbackValue string) (string, error) {
	user, err := evaluationcontext.UserFromContext(ctx)
	if err != nil {
		return fallbackValue, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.StringVariation(string(key), user.ToLDUser(), fallbackValue)
}

// QueryStringWithEvaluationContext retrieves the value of a string flag. An evaluation context
// must be supplied manually. The supplied fallback value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryStringWithEvaluationContext(key FlagName, evalContext evaluationcontext.Context, fallbackValue string) (string, error) {
	return c.wrappedClient.StringVariation(string(key), evalContext.ToLDUser(), fallbackValue)
}

// QueryInt retrieves the value of an integer flag. User attributes are
// extracted from the context. The supplied fallback value is always reflected in
// the returned value regardless of whether an error occurs.
func (c *Client) QueryInt(ctx context.Context, key FlagName, fallbackValue int) (int, error) {
	user, err := evaluationcontext.UserFromContext(ctx)
	if err != nil {
		return fallbackValue, fmt.Errorf("get user from context: %w", err)
	}

	return c.wrappedClient.IntVariation(string(key), user.ToLDUser(), fallbackValue)
}

// QueryIntWithEvaluationContext retrieves the value of an integer flag. An evaluation context
// must be supplied manually. The supplied fallback value is always reflected in the
// returned value regardless of whether an error occurs.
func (c *Client) QueryIntWithEvaluationContext(key FlagName, evalContext evaluationcontext.Context, fallbackValue int) (int, error) {
	return c.wrappedClient.IntVariation(string(key), evalContext.ToLDUser(), fallbackValue)
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
