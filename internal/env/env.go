package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ENV                   string
	AWS_REGION            string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
)

func Init() {
	// Load env files if provided.
	if envFiles := os.Getenv("ENV_FILES"); envFiles != "" {
		if err := godotenv.Load(strings.Split(envFiles, ",")...); err != nil {
			panic(err)
		}
	}

	// Fields that has to come from env.
	if ENV = strings.ToUpper(os.Getenv("ENV")); ENV == "" {
		panic("ENV is required")
	}

	if AWS_REGION = os.Getenv("AWS_REGION"); AWS_REGION == "" {
		panic("AWS_REGION is required")
	}

	AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
}

// Github actions always sets CI to true.
func IsCI() bool {
	return os.Getenv("CI") == "true"
}

func IsDev() bool {
	env := strings.ToUpper(ENV)
	return (env == "" && !IsCI()) || env == EnvDev || env == EnvDevelopment
}

func IsStaging() bool {
	env := strings.ToUpper(ENV)
	return env == EnvStag || env == EnvStaging
}

func IsProd() bool {
	env := strings.ToUpper(ENV)
	return env == EnvProd || env == EnvProduction
}
