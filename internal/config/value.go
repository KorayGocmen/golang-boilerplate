package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	secretsManagerTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	parameterStoreTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/koraygocmen/golang-boilerplate/internal/aws"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
)

type Type string

const (
	TypeSecret = Type("secret")
	TypeParam  = Type("param")
)

type Param struct {
	Key     string
	Type    Type
	Panic   bool
	Default interface{}
}

func get(ctx context.Ctx, p Param) (string, error) {
	val := strings.TrimSpace(os.Getenv(p.Key))

	if val == "" {
		if p.Type == TypeSecret {
			val = fmt.Sprintf("secretsmanager:/%s/%s", strings.ToLower(env.ENV), strings.ToUpper(p.Key))
		}
		if p.Type == TypeParam {
			val = fmt.Sprintf("ssm:/%s/%s", strings.ToLower(env.ENV), strings.ToUpper(p.Key))
		}
	}

	if strings.HasPrefix(val, "secretsmanager:") {
		var err error
		val, err = aws.Client.SecretGet(ctx, strings.TrimPrefix(val, "secretsmanager:"))
		if err != nil {
			var pnfe *secretsManagerTypes.ResourceNotFoundException
			if errors.As(err, &pnfe) {
				goto done
			}

			err = fmt.Errorf("get config param error: %w", err)
			return "", err
		}
		goto done
	}

	if strings.HasPrefix(val, "ssm:") {
		var err error
		val, err = aws.Client.ParamGet(ctx, strings.TrimPrefix(val, "ssm:"))
		if err != nil {
			var pnfe *parameterStoreTypes.ParameterNotFound
			if errors.As(err, &pnfe) {
				goto done
			}

			err = fmt.Errorf("get config param error: %w", err)
			return "", err
		}
		goto done
	}

done:
	if val == "" {
		if p.Default != nil {
			val = fmt.Sprintf("%v", p.Default)
		} else if p.Panic {
			err := fmt.Errorf("config param is required: %s", p.Key)
			panic(err)
		}
	}

	return val, nil
}

func GetStr(ctx context.Ctx, p Param) string {
	val, err := get(ctx, p)
	if err != nil {
		if p.Panic {
			err := fmt.Errorf("get config param (str) error: %s: %w", p.Key, err)
			panic(err)
		}
		return ""
	}

	if val == "" {
		return ""
	}

	return val
}

func GetInt(ctx context.Ctx, p Param) int {
	val, err := get(ctx, p)
	if err != nil {
		return 0
	}

	if val == "" {
		return 0
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}

	return valInt
}

func GetBool(ctx context.Ctx, p Param) bool {
	val, err := get(ctx, p)
	if err != nil {
		return false
	}

	if val == "" {
		return false
	}

	valBool, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return valBool
}
