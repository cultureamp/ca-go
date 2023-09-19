package awsconfig

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func envSetter() func() {
	envs := os.Environ()
	originalEnvs := map[string]string{}

	for _, v := range envs {
		arr := strings.Split(v, "=")
		key := arr[0]
		value := arr[1]
		originalEnvs[key] = value
	}

	return func() {
		os.Clearenv()
		for name, value := range originalEnvs {
			_ = os.Setenv(name, value)
		}
	}
}

func TestGetAwsConfig(t *testing.T) {
	t.Cleanup(func() { envSetter() })
	os.Clearenv()
	ctx := context.Background()
	localstackhost := "localstack"
	t.Setenv("LOCALSTACK_HOST", localstackhost)
	t.Setenv("HOME", "~")
	t.Setenv("AWS_ACCESS_KEY_ID", "fake-access-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "fake-secret-access-key")

	t.Run("Should default when env doesn't have local dev", func(t *testing.T) {
		// arrange
		// act
		config, err := GetAwsConfig(ctx)
		// assert
		require.NoError(t, err)
		require.NotNil(t, config)
		assert.Nil(t, config.EndpointResolverWithOptions)
	})

	t.Run("Should override endpoint when env has local dev", func(t *testing.T) {
		// arrange
		t.Setenv("IS_LOCAL_DEV", "true")
		// act
		config, err := GetAwsConfig(ctx)
		// assert
		require.NoError(t, err)
		require.NotNil(t, config)
		assert.NotNil(t, config.EndpointResolverWithOptions)
		resovledEndpoint, err := config.EndpointResolverWithOptions.ResolveEndpoint("", "")
		require.NoError(t, err)
		hostName := net.JoinHostPort(localstackhost, "4566")
		assert.Equal(t, fmt.Sprintf("http://%s", hostName), resovledEndpoint.URL)
	})

	t.Run("Should not override endpoint when env has local dev false", func(t *testing.T) {
		// arrange
		t.Setenv("IS_LOCAL_DEV", "false")
		// act
		config, err := GetAwsConfig(ctx)
		// assert
		require.NoError(t, err)
		require.NotNil(t, config)
		assert.Nil(t, config.EndpointResolverWithOptions)
	})
}
