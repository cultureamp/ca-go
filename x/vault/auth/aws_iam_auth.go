package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
	vaultapi "github.com/hashicorp/vault/api"
)

const (
	loginPath = "auth/aws/login"
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

func (auth *AWSIamAuth) Login(ctx context.Context, client *vaultapi.Client) (*vaultapi.Secret, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		region = ""
	}
	loginData, err := awsutil.GenerateLoginData(
		stscreds.NewCredentials(awsSession, auth.roleArn),
		"",
		region,
		hclog.Default(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to generate login data for AWS auth endpoint: %w", err)
	}

	loginData["role"] = auth.vaultRole

	secret, err := client.Logical().WriteWithContext(ctx, loginPath, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with AWS auth: %w", err)
	}

	return secret, nil
}
