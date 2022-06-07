package kafkaclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultProducerConfiguration(t *testing.T) {
	// This client exercises the underlying xdg-go/scram library. This test just
	// ensures that it looks basically correct, it's not trying to be a
	// behavioural test.

	client := xDGSCRAMClient{}

	require.NoError(t, client.Begin("user", "password", "authzID"))

	str, err := client.Step("challenge")
	require.NoError(t, err)
	require.Regexp(t, "^n,authzID,n=user,r=", str)

	require.False(t, client.Done())
}
