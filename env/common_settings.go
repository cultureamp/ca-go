package env

import (
	senv "github.com/caarlos0/env/v10"
)

// commonSettings that drive behavior used by at least 90% of apps.
type commonSettings struct {
	App        string `env:"APP,required,notEmpty"`
	AppVersion string `env:"APP_VERSION" envDefault:"1.0.0"`
	AppEnv     string `env:"APP_ENV"     envDefault:"development"`
	Farm       string `env:"FARM"        envDefault:"local"`
	Product    string `env:"PRODUCT"`
}

func newCommonSettings() *commonSettings {
	settings := commonSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// GetAppName returns the application name from the "APP" environment variable.
func (s *commonSettings) GetAppName() string {
	return s.App
}

// GetAppVersion returns the application version from the "APP_VER" environment variable.
// Default: "1.0.0".
func (s *commonSettings) GetAppVersion() string {
	return s.AppVersion
}

// GetAppEnv returns the application environment from the "APP_ENV" environment variable.
// Examples: "development", "production".
func (s *commonSettings) GetAppEnv() string {
	return s.AppEnv
}

// GetFarm returns the farm running the application from the "FARM" environment variable.
// Examples: "local", "dolly", "production".
func (s *commonSettings) GetFarm() string {
	return s.Farm
}

// GetProductSuite returns the product suite this application belongs to from the "PRODUCT" environment variable.
// Examples: "engagement", "performance".
func (s *commonSettings) GetProductSuite() string {
	return s.Product
}

// IsProduction returns true if "APP_ENV" == "production".
func (s *commonSettings) IsProduction() bool {
	return s.AppEnv == "production"
}

// IsRunningInAWS returns true if "APP_ENV" != "local".
func (s *commonSettings) IsRunningInAWS() bool {
	return !s.IsRunningLocal()
}

// IsRunningLocal returns true if FARM" == "local".
func (s *commonSettings) IsRunningLocal() bool {
	return s.Farm == "local"
}
