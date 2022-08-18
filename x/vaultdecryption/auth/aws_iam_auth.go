package auth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	vaultapi "github.com/hashicorp/vault/api"
)

const (
	defaultStsRegion = "us-east-1"
	loginPath        = "auth/aws/login"
)

type AWSIamAuth struct {
	vaultRole string
	roleArn   string
}

func NewAWSIamAuth(vaultRole string, roleArn string) *AWSIamAuth {
	return &AWSIamAuth{
		vaultRole: vaultRole,
		roleArn:   roleArn,
	}
}

func (a *AWSIamAuth) Login(ctx context.Context, client *vaultapi.Client) (*vaultapi.Secret, error) {
	var awsSession *session.Session
	var err error
	if awsSession, err = session.NewSession(); err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	loginData, err := awsutil.GenerateLoginData(
		stscreds.NewCredentials(awsSession, a.roleArn),
		"",
		defaultStsRegion,
		hclog.Default(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to generate login data for AWS auth endpoint: %w", err)
	}

	loginData["role"] = a.vaultRole

	secret, err := client.Logical().WriteWithContext(ctx, loginPath, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with AWS auth: %w", err)
	}

	return secret, nil
}
