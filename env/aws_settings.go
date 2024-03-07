package env

import (
	senv "github.com/caarlos0/env/v10"
)

// awsSettings that drive behavior.
type awsSettings struct {
	AwsProfile   string `env:"AWS_PROFILE"    envDefault:"default"`
	AwsRegion    string `env:"AWS_REGION,required,notEmpty"`
	AwsAccountID string `env:"AWS_ACCOUNT_ID"`
	XrayLogging  bool   `env:"XRAY_LOGGING"   envDefault:"true"`
}

func newAWSSettings() *awsSettings {
	settings := awsSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// GetAwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func (s *awsSettings) GetAwsProfile() string {
	return s.AwsProfile
}

// GetAwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func (s *awsSettings) GetAwsRegion() string {
	return s.AwsRegion
}

// GetAwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
func (s *awsSettings) GetAwsAccountID() string {
	return s.AwsAccountID
}

// IsXrayTracingEnabled returns "true" if the "XRAY_LOGGING" environment variable is turned on.
func (s *awsSettings) IsXrayTracingEnabled() bool {
	return s.XrayLogging
}
