package flags

import (
	"fmt"
)

var flagsClient *Client

// FlagName establishes a type for flag names.
type FlagName string

// Configure configures the client as a managed singleton.
func Configure(opts ...ConfigOption) error {
	c, err := NewClient(opts...)
	if err != nil {
		return fmt.Errorf("configure client: %w", err)
	}

	flagsClient = c
	return nil
}

// Connect attempts to connect the managed singleton to LaunchDarkly. An error
// is returned if the singleton is not yet configured, a connection has already
// been established, or a connection error occurs.
func Connect() error {
	if flagsClient == nil {
		return errClientNotConfigured
	}

	return flagsClient.Connect()
}

// GetDefaultClient returns the managed singleton client. An error is returned
// if the client is not yet configured.
func GetDefaultClient() (*Client, error) {
	if flagsClient == nil {
		return nil, errClientNotConfigured
	}

	return flagsClient, nil
}
