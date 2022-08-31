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

	if err := c.Connect(); err != nil {
		return fmt.Errorf("connect client: %w", err)
	}

	flagsClient = c
	return nil
}

// GetDefaultClient returns the managed singleton client. An error is returned
// if the client is not yet configured.
func GetDefaultClient() (*Client, error) {
	if flagsClient == nil {
		return nil, errClientNotConfigured
	}

	return flagsClient, nil
}
