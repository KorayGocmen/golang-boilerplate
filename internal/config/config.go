package config

import (
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
)

var (
	Server   = ServerConfig{}
	Database = DatabaseConfig{}
	AWS      = AwsConfig{}
	Slack    = SlackConfig{}
	Log      = LogConfig{}
)

type ServerConfig struct {
	Addr    string
	Timeout struct {
		Read  int
		Write int
		Idle  int
	}
}

type DatabaseConfig struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DB      string
	SSLMode string

	MaxIdleConns  int
	MaxOpenConns  int
	MaxIdleTimeMs int
	MaxLifetimeMs int

	Migrations struct {
		Auto   bool
		Source string
		Type   string
	}
}

type AwsConfig struct {
	S3 struct {
		BucketDocuments string
	}
	SES struct {
		Source string
	}
}

type SlackConfig struct {
	Webhook struct {
		Errors string
		Events string
	}
}

type LogConfig struct {
	Level int
	Mode  string

	Syslog struct {
		Addr     string
		Protocol string
		Tag      string
	}

	File struct {
		Path string
	}

	AWS struct {
		LogGroup struct {
			Name   string
			Region string
		}
	}
}

func Load() {
	ctx := context.Background()

	Server.Addr = GetStr(ctx, Param{Key: "SERVER_ADDR", Type: TypeParam, Panic: true})
	Server.Timeout.Read = GetInt(ctx, Param{Key: "SERVER_TIMEOUT_READ", Type: TypeParam, Panic: true})
	Server.Timeout.Write = GetInt(ctx, Param{Key: "SERVER_TIMEOUT_WRITE", Type: TypeParam, Panic: true})
	Server.Timeout.Idle = GetInt(ctx, Param{Key: "SERVER_TIMEOUT_IDLE", Type: TypeParam, Panic: true})

	Database.Host = GetStr(ctx, Param{Key: "DATABASE_HOST", Type: TypeParam, Panic: true})
	Database.Port = GetStr(ctx, Param{Key: "DATABASE_PORT", Type: TypeParam, Panic: true})
	Database.DB = GetStr(ctx, Param{Key: "DATABASE_DB", Type: TypeParam, Panic: true})
	Database.User = GetStr(ctx, Param{Key: "DATABASE_USER", Type: TypeParam, Panic: true})
	Database.Pass = GetStr(ctx, Param{Key: "DATABASE_PASS", Type: TypeSecret, Panic: true})
	Database.SSLMode = GetStr(ctx, Param{Key: "DATABASE_SSLMODE", Type: TypeParam, Panic: true})

	Database.MaxIdleConns = GetInt(ctx, Param{Key: "DATABASE_MAX_IDLE_CONNS", Type: TypeParam, Panic: true})
	Database.MaxOpenConns = GetInt(ctx, Param{Key: "DATABASE_MAX_OPEN_CONNS", Type: TypeParam, Panic: true})
	Database.MaxIdleTimeMs = GetInt(ctx, Param{Key: "DATABASE_MAX_IDLE_TIME_MS", Type: TypeParam, Panic: true})
	Database.MaxLifetimeMs = GetInt(ctx, Param{Key: "DATABASE_MAX_LIFETIME_MS", Type: TypeParam, Panic: true})

	Database.Migrations.Auto = GetBool(ctx, Param{Key: "DATABASE_MIGRATIONS_AUTO", Type: TypeParam, Panic: true})
	Database.Migrations.Source = GetStr(ctx, Param{Key: "DATABASE_MIGRATIONS_SOURCE", Type: TypeParam, Panic: true})
	Database.Migrations.Type = GetStr(ctx, Param{Key: "DATABASE_MIGRATIONS_TYPE", Type: TypeParam, Panic: true})

	AWS.S3.BucketDocuments = GetStr(ctx, Param{Key: "AWS_S3_BUCKET_DOCUMENTS", Type: TypeParam, Panic: true})
	AWS.SES.Source = GetStr(ctx, Param{Key: "AWS_SES_SOURCE", Type: TypeParam, Panic: true})

	Slack.Webhook.Errors = GetStr(ctx, Param{Key: "SLACK_WEBHOOK_ERRORS", Type: TypeSecret, Panic: env.IsProd()})
	Slack.Webhook.Events = GetStr(ctx, Param{Key: "SLACK_WEBHOOK_EVENTS", Type: TypeSecret, Panic: env.IsProd()})

	Log.Mode = GetStr(ctx, Param{Key: "LOG_MODE", Type: TypeParam, Panic: true})
	Log.Level = GetInt(ctx, Param{Key: "LOG_LEVEL", Type: TypeParam, Panic: true})
	Log.Syslog.Addr = GetStr(ctx, Param{Key: "LOG_SYSLOG_ADDR", Type: TypeParam, Panic: false})
	Log.Syslog.Protocol = GetStr(ctx, Param{Key: "LOG_SYSLOG_PROTOCOL", Type: TypeParam, Panic: false})
	Log.Syslog.Tag = GetStr(ctx, Param{Key: "LOG_SYSLOG_TAG", Type: TypeParam, Panic: false})
	Log.AWS.LogGroup.Name = GetStr(ctx, Param{Key: "LOG_AWS_LOG_GROUP_NAME", Type: TypeParam, Panic: false})
	Log.AWS.LogGroup.Region = GetStr(ctx, Param{Key: "LOG_AWS_LOG_GROUP_REGION", Type: TypeParam, Panic: false})
	Log.File.Path = GetStr(ctx, Param{Key: "LOG_FILE_PATH", Type: TypeParam, Panic: false})
}
