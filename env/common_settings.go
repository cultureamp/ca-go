package env

import (
	senv "github.com/caarlos0/env/v11"
)

// CommonSettings implements common settings used in 90% of all our apps.
// This is an interface so that clients can mock out this behaviour in tests.
type CommonSettings interface {
	AppName() string
	AppVersion() string
	AppEnv() string
	Farm() string
	ProductSuite() string
	IsProduction() bool
	IsRunningInAWS() bool
	IsRunningLocal() bool
}

// commonSettings that drive behavior used by at least 90% of apps.
type commonSettings struct {
	// These have to be public so that "github.com/caarlos0/env/v10" can populate them
	AppNameEnv      string `env:"APP"         envDefault:"<unknown>"`
	AppVersionEnv   string `env:"APP_VERSION" envDefault:"1.0.0"`
	AppTypeEnv      string `env:"APP_ENV"     envDefault:"development"`
	FarmNameEnv     string `env:"FARM"        envDefault:"local"`
	ProductSuiteEnv string `env:"PRODUCT"     envDefault:"<unknown>"`
}

func newCommonSettings() *commonSettings {
	settings := commonSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// GetAppName returns the application name from the "APP" environment variable.
func (s *commonSettings) AppName() string {
	return s.AppNameEnv
}

// tAppVersion returns the application version from the "APP_VER" environment variable.
// Default: "1.0.0".
func (s *commonSettings) AppVersion() string {
	return s.AppVersionEnv
}

// AppEnv returns the application environment from the "APP_ENV" environment variable.
// Examples: "development", "production".
func (s *commonSettings) AppEnv() string {
	return s.AppTypeEnv
}

// Farm returns the farm running the application from the "FARM" environment variable.
// Examples: "local", "dolly", "production".
func (s *commonSettings) Farm() string {
	return s.FarmNameEnv
}

// ProductSuite returns the product suite this application belongs to from the "PRODUCT" environment variable.
// Examples: "engagement", "performance".
func (s *commonSettings) ProductSuite() string {
	return s.ProductSuiteEnv
}

// IsProduction returns true if "APP_ENV" == "production".
func (s *commonSettings) IsProduction() bool {
	return s.AppTypeEnv == "production"
}

// IsRunningInAWS returns true if "APP_ENV" != "local".
func (s *commonSettings) IsRunningInAWS() bool {
	return !s.IsRunningLocal()
}

// IsRunningLocal returns true if FARM" == "local".
func (s *commonSettings) IsRunningLocal() bool {
	return s.FarmNameEnv == "local"
}
