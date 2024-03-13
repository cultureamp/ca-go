package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommonSettings(t *testing.T) {
	settings := newCommonSettings()
	assert.NotNil(t, settings)
}

func TestCommonSettings(t *testing.T) {
	t.Setenv(AppNameEnv, "ca-go-unit-tests")
	t.Setenv(AppVerEnv, "1.2.3")
	t.Setenv(AppEnvironmentEnv, "local")
	t.Setenv(AppFarmEnv, "local")
	t.Setenv(ProductEnv, "standard_library")

	settings := newCommonSettings()
	assert.Equal(t, "ca-go-unit-tests", settings.CSApp)
	assert.Equal(t, "1.2.3", settings.CSAppVersion)
	assert.Equal(t, "local", settings.CSAppEnv)
	assert.Equal(t, "local", settings.CSFarm)
	assert.Equal(t, "standard_library", settings.CSProduct)
}

func TestSettingsHelpers(t *testing.T) {
	t.Setenv(AppEnvironmentEnv, "production")
	settings := newCommonSettings()
	isProd := settings.IsProduction()
	assert.True(t, isProd)

	t.Setenv(AppEnvironmentEnv, "dev")
	settings = newCommonSettings()
	isProd = settings.IsProduction()
	assert.False(t, isProd)
}

func Test_Settings_Env_IsAws_IsLocal(t *testing.T) {
	t.Setenv(AppFarmEnv, "local")
	settings := newCommonSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "local", settings.CSFarm)
	assert.True(t, settings.IsRunningLocal())
	assert.False(t, settings.IsRunningInAWS())

	t.Setenv(AppFarmEnv, "falcon")
	settings = newCommonSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "falcon", settings.CSFarm)
	assert.False(t, settings.IsRunningLocal())
	assert.True(t, settings.IsRunningInAWS())

	t.Setenv(AppFarmEnv, "production")
	settings = newCommonSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "production", settings.CSFarm)
	assert.False(t, settings.IsRunningLocal())
	assert.True(t, settings.IsRunningInAWS())
}
