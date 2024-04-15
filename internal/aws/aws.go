package aws

import (
	"fmt"

	awsgoConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
)

type client struct {
	secretsManagerClient *secretsmanager.Client
	cloudwatchLogsClient *cloudwatchlogs.Client

	ParamGet  func(ctx context.Ctx, key string) (string, error)
	SecretGet func(ctx context.Ctx, key string) (string, error)
	LogStream func(ctx context.Ctx, groupName, streamUnique string) (*Stream, error)
}

var (
	// Initialize AWS with nil values.
	Client = client{
		ParamGet:  paramGet(nil),
		SecretGet: secretGet(nil),
		LogStream: logStream(nil),
	}
)

func Init() error {
	// AWS config is loaded from the environment variables.
	cfg, err := awsgoConfig.LoadDefaultConfig(context.Background(), awsgoConfig.WithRegion(env.AWS_REGION))
	if err != nil {
		err = fmt.Errorf("failed to create aws config: %w", err)
		return err
	}

	ssmClient := ssm.NewFromConfig(cfg)
	secretsManagerClient := secretsmanager.NewFromConfig(cfg)
	cloudwatchLogsClient := cloudwatchlogs.NewFromConfig(cfg)

	Client = client{
		secretsManagerClient: secretsManagerClient,
		cloudwatchLogsClient: cloudwatchLogsClient,

		ParamGet:  paramGet(ssmClient),
		SecretGet: secretGet(secretsManagerClient),
		LogStream: logStream(cloudwatchLogsClient),
	}

	return nil
}
