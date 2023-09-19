package awsconfig

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
)

func getIsLocalBool() bool {
	value, err := strconv.ParseBool(os.Getenv("IS_LOCAL_DEV"))
	if err != nil {
		return false
	}
	return value
}

func GetAwsConfig(ctx context.Context) (aws.Config, error) {
	var cfg aws.Config
	var err error
	if getIsLocalBool() {
		hostName := net.JoinHostPort(os.Getenv("LOCALSTACK_HOST"), "4566")
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           fmt.Sprintf("http://%s", hostName),
				PartitionID:   "aws",
				SigningRegion: region,
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithEndpointResolverWithOptions(customResolver))
		if err != nil {
			return cfg, errors.Wrap(err, "failed to load local AWS Configurations")
		}
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
		if err != nil {
			return cfg, errors.Wrap(err, "failed to load AWS Configurations")
		}
	}

	return cfg, nil
}
