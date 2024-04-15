package aws

import (
	"fmt"

	awsgo "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
)

func paramGet(ssmClient *ssm.Client) func(ctx context.Ctx, key string) (string, error) {
	return func(ctx context.Ctx, key string) (string, error) {
		if ssmClient == nil {
			err := fmt.Errorf("param get error: ssm client is nil, maybe aws is not initialized?")
			return "", err
		}

		input := &ssm.GetParameterInput{
			Name:           awsgo.String(key),
			WithDecryption: awsgo.Bool(true),
		}
		result, err := ssmClient.GetParameter(ctx, input)
		if err != nil {
			err = fmt.Errorf("param get error: %w", err)
			return "", err
		}

		return *result.Parameter.Value, nil
	}
}
