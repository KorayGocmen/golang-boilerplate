package aws

import (
	"fmt"

	awsgo "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
)

func secretGet(secretsManagerClient *secretsmanager.Client) func(ctx context.Ctx, key string) (string, error) {
	return func(ctx context.Ctx, key string) (string, error) {
		if secretsManagerClient == nil {
			err := fmt.Errorf("secret get error: secrets manager client is nil, maybe aws is not initialized?")
			return "", err
		}

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     awsgo.String(key),
			VersionStage: awsgo.String("AWSCURRENT"),
		}
		result, err := secretsManagerClient.GetSecretValue(ctx, input)
		if err != nil {
			err = fmt.Errorf("secret get error: %s: %v", key, err)
			return "", err
		}

		return *result.SecretString, nil
	}
}
