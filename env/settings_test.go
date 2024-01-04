package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Settings_New(t *testing.T) {
	settings := NewSettings()
	assert.NotNil(t, settings)
}

func Test_Settings_Defaults(t *testing.T) {
	defer os.Unsetenv(AppNameEnv)
	t.Setenv(AppNameEnv, "ca-go-unit-tests")

	defer os.Unsetenv(AppFarmEnv)
	t.Setenv(AppFarmEnv, "local")

	settings := NewSettings()

	assert.Equal(t, "ca-go-unit-tests", settings.App)
	assert.Equal(t, "local", settings.Farm)
}

func Test_Settings_IsProduction(t *testing.T) {
	defer os.Unsetenv(AppEnv)

	t.Setenv(AppEnv, "production")
	settings := NewSettings()
	isProd := settings.IsProduction()
	assert.True(t, isProd)

	t.Setenv(AppEnv, "dev")
	settings = NewSettings()
	isProd = settings.IsProduction()
	assert.False(t, isProd)
}

func Test_Settings_Env_IsAws_IsLocal(t *testing.T) {
	defer os.Unsetenv(AppFarmEnv)

	t.Setenv(AppFarmEnv, "local")
	settings := NewSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "local", settings.Farm)
	assert.True(t, settings.IsRunningLocal())
	assert.False(t, settings.IsRunningInAWS())

	t.Setenv(AppFarmEnv, "falcon")
	settings = NewSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "falcon", settings.Farm)
	assert.False(t, settings.IsRunningLocal())
	assert.True(t, settings.IsRunningInAWS())

	t.Setenv(AppFarmEnv, "production")
	settings = NewSettings()
	assert.NotNil(t, settings)
	assert.Equal(t, "production", settings.Farm)
	assert.False(t, settings.IsRunningLocal())
	assert.True(t, settings.IsRunningInAWS())
}

func Test_Settings_JSON(t *testing.T) {
	defer os.Unsetenv(AppFarmEnv)
	defer os.Unsetenv(AppEnv)
	defer os.Unsetenv(DatadogAPIEnvVar)
	defer os.Unsetenv(SentryDsnEnv)

	t.Setenv(AppFarmEnv, "farm")
	t.Setenv(AppEnv, "development")
	t.Setenv(DatadogAPIEnvVar, "1234567890")
	t.Setenv(SentryDsnEnv, "1234567890")

	settings := NewSettings()
	json := settings.ToJSON()
	assert.Equal(t, "{\"app\":\"authz-api\",\"app_version\":\"1.0.0\",\"app_env\":\"development\",\"farm\":\"farm\",\"product\":\"\",\"aws_profile\":\"default\",\"aws_region\":\"us-west-2\",\"aws_account_id\":\"\",\"xray_logging\":true,\"dd_api_key\":\"1234567890\",\"sentry_dsn\":\"1234567890\",\"sentry_flush_ms\":50}", json)
}

func Test_Settings_Redacted_JSON(t *testing.T) {
	defer os.Unsetenv(AppFarmEnv)
	defer os.Unsetenv(AppEnv)
	defer os.Unsetenv(DatadogAPIEnvVar)
	defer os.Unsetenv(SentryDsnEnv)

	t.Setenv(AppFarmEnv, "farm")
	t.Setenv(AppEnv, "development")
	t.Setenv(DatadogAPIEnvVar, "1234567890")
	t.Setenv(SentryDsnEnv, "1234567890")

	settings := NewSettings()
	json := settings.ToRedactedJSON()
	assert.Equal(t, "{\"app\":\"authz-api\",\"app_version\":\"1.0.0\",\"app_env\":\"development\",\"farm\":\"farm\",\"product\":\"\",\"aws_profile\":\"default\",\"aws_region\":\"us-west-2\",\"aws_account_id\":\"\",\"xray_logging\":true,\"dd_api_key\":\"******7890\",\"sentry_dsn\":\"******7890\",\"sentry_flush_ms\":50}", json)
}

func Test_Settings_Redacted_String(t *testing.T) {
	defer os.Unsetenv(AppFarmEnv)
	defer os.Unsetenv(AppEnv)
	defer os.Unsetenv(DatadogAPIEnvVar)
	defer os.Unsetenv(SentryDsnEnv)

	t.Setenv(AppFarmEnv, "farm")
	t.Setenv(AppEnv, "development")
	t.Setenv(DatadogAPIEnvVar, "1234567890")
	t.Setenv(SentryDsnEnv, "1234567890")

	settings := NewSettings()
	s := settings.ToRedactedString()
	assert.Equal(t, "{app:authz-api,app_version:1.0.0,app_env:development,farm:farm,product:,aws_profile:default,aws_region:us-west-2,aws_account_id:,xray_logging:true,dd_api_key:******7890,sentry_dsn:******7890,sentry_flush_ms:50}", s)
}
