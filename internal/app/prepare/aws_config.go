package prepare

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func AwsConfig(ctx context.Context) (aws.Config, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-central-1"))
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load SDK config, %v", err)
	}
	return awsCfg, nil
}

func LocalAwsConfig(ctx context.Context) (aws.Config, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("eu-central-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "http://localhost:8000"}, nil
		})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "local development dynamodb",
			},
		}),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load SDK config, %v", err)
	}
	return awsCfg, nil
}
