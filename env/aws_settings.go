package env

import (
	senv "github.com/caarlos0/env/v10"
)

// AWSSettings implements AWS settings.
// This is an interface so that clients can mock out this behaviour in tests.
type AWSSettings interface {
	AwsProfile() string
	AwsRegion() string
	AwsAccountID() string
	IsXrayTracingEnabled() bool
}

// awsSettings that drive behavior.
type awsSettings struct {
	// These have to be public so that "github.com/caarlos0/env/v10" can populate them
	Aws_Profile   string `env:"AWS_PROFILE"    envDefault:"default"`
	Aws_Region    string `env:"AWS_REGION"     envDefault:"dev"`
	Aws_AccountID string `env:"AWS_ACCOUNT_ID" envDefault:"<unknown>"`
	Xray_Logging  bool   `env:"XRAY_LOGGING"   envDefault:"true"`
}

func newAWSSettings() *awsSettings {
	settings := awsSettings{}
	if err := senv.Parse(&settings); err != nil {
		panic(err)
	}

	return &settings
}

// AwsProfile returns the AWS profile from the "AWS_PROFILE" environment variable.
func (s *awsSettings) AwsProfile() string {
	return s.Aws_Profile
}

// AwsRegion returns the AWS region from the "AWS_REGION" environment variable.
func (s *awsSettings) AwsRegion() string {
	return s.Aws_Region
}

// AwsAccountID returns the AWS region from the "AWS_ACCOUNT_ID" environment variable.
func (s *awsSettings) AwsAccountID() string {
	return s.Aws_AccountID
}

// IsXrayTracingEnabled returns "true" if the "XRAY_LOGGING" environment variable is turned on.
func (s *awsSettings) IsXrayTracingEnabled() bool {
	return s.Xray_Logging
}
