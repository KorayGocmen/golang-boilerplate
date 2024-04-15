package config

import (
	"os"
	"testing"

	secretsManagerTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	parameterStoreTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/koraygocmen/golang-boilerplate/internal/aws"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
)

func setRequiredFields() {
	// Set environment variables.
	os.Setenv("SERVER_ADDR", "server_addr")
	os.Setenv("SERVER_TIMEOUT_READ", "1")
	os.Setenv("SERVER_TIMEOUT_WRITE", "2")
	os.Setenv("SERVER_TIMEOUT_IDLE", "3")

	os.Setenv("DATABASE_HOST", "database_host")
	os.Setenv("DATABASE_PORT", "4")
	os.Setenv("DATABASE_USER", "database_user")
	os.Setenv("DATABASE_PASS", "database_pass")
	os.Setenv("DATABASE_DB", "database_db")
	os.Setenv("DATABASE_SSLMODE", "enable")

	os.Setenv("DATABASE_MAX_IDLE_CONNS", "10")
	os.Setenv("DATABASE_MAX_OPEN_CONNS", "10")
	os.Setenv("DATABASE_MAX_IDLE_TIME_MS", "0")
	os.Setenv("DATABASE_MAX_LIFETIME_MS", "1800000")

	os.Setenv("DATABASE_MIGRATIONS_AUTO", "true")
	os.Setenv("DATABASE_MIGRATIONS_SOURCE", "database_migrations_source")
	os.Setenv("DATABASE_MIGRATIONS_TYPE", "database_migrations_type")

	os.Setenv("AWS_S3_BUCKET_DOCUMENTS", "aws_s3_bucket_documents")
	os.Setenv("AWS_SES_SOURCE", "aws_ses_source")

	os.Setenv("SLACK_WEBHOOK_ERRORS", "slack_webhook_errors")
	os.Setenv("SLACK_WEBHOOK_EVENTS", "slack_webhook_events")

	os.Setenv("LOG_MODE", "console")
	os.Setenv("LOG_LEVEL", "5")
	os.Setenv("LOG_FILE_PATH", "log_file_path")
	os.Setenv("LOG_SYSLOG_ADDR", "log_syslog_addr")
	os.Setenv("LOG_SYSLOG_PROTOCOL", "log_syslog_protocol")
	os.Setenv("LOG_SYSLOG_TAG", "log_syslog_tag")
	os.Setenv("LOG_AWS_LOG_GROUP_NAME", "log_aws_log_group_name")
	os.Setenv("LOG_AWS_LOG_GROUP_REGION", "log_aws_log_group_region")
	os.Setenv("LOG_AWS_LOG_STREAM_SUFFIX", "log_aws_log_stream_suffix")
}

func testLoadPanic(testcase string, t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Fatalf("%s did not panic", testcase)
}

func TestLoad(t *testing.T) {
	setRequiredFields()

	// Load
	Load()

	// Server.
	if Server.Addr != "server_addr" {
		t.Fatalf("Server.Addr = %s; want server_addr", Server.Addr)
	}
	if Server.Timeout.Read != 1 {
		t.Fatalf("Server.Timeout.Read = %d; want 1", Server.Timeout.Read)
	}
	if Server.Timeout.Write != 2 {
		t.Fatalf("Server.Timeout.Write = %d; want 2", Server.Timeout.Write)
	}
	if Server.Timeout.Idle != 3 {
		t.Fatalf("Server.Timeout.Idle = %d; want 3", Server.Timeout.Idle)
	}

	// Database.
	if Database.Host != "database_host" {
		t.Fatalf("Database.Host = %s; want database_host", Database.Host)
	}
	if Database.Port != "4" {
		t.Fatalf("Database.Port = %s; want 4", Database.Port)
	}
	if Database.User != "database_user" {
		t.Fatalf("Database.User = %s; want database_user", Database.User)
	}
	if Database.Pass != "database_pass" {
		t.Fatalf("Database.Pass = %s; want database_pass", Database.Pass)
	}
	if Database.DB != "database_db" {
		t.Fatalf("Database.DB = %s; want database_db", Database.DB)
	}
	if Database.SSLMode != "enable" {
		t.Fatalf("Database.SSLMode = %s; want enable", Database.SSLMode)
	}
	if Database.Migrations.Auto != true {
		t.Fatalf("Database.Migrations.Auto = %t; want true", Database.Migrations.Auto)
	}
	if Database.Migrations.Source != "database_migrations_source" {
		t.Fatalf("Database.Migrations.Source = %s; want database_migrations_source", Database.Migrations.Source)
	}
	if Database.Migrations.Type != "database_migrations_type" {
		t.Fatalf("Database.Migrations.Type = %s; want database_migrations_type", Database.Migrations.Type)
	}

	// AWS.
	if AWS.S3.BucketDocuments != "aws_s3_bucket_documents" {
		t.Fatalf("AWS.S3.Bucket.Documents = %s; want aws_s3_bucket_documents", AWS.S3.BucketDocuments)
	}
	if AWS.SES.Source != "aws_ses_source" {
		t.Fatalf("AWS.SES.Source = %s; want aws_ses_source", AWS.SES.Source)
	}

	// Slack.
	if Slack.Webhook.Errors != "slack_webhook_errors" {
		t.Fatalf("Slack.WebhookErrorsURL = %s; want slack_webhook_errors_url", Slack.Webhook.Errors)
	}
	if Slack.Webhook.Events != "slack_webhook_events" {
		t.Fatalf("Slack.WebhookEventsURL = %s; want slack_webhook_events_url", Slack.Webhook.Events)
	}

	// Log.
	if Log.Mode != "console" {
		t.Fatalf("Log.Mode = %s; want log_mode", Log.Mode)
	}
	if Log.Level != 5 {
		t.Fatalf("Log.Level = %d; want 5", Log.Level)
	}
	if Log.Syslog.Addr != "log_syslog_addr" {
		t.Fatalf("Log.Syslog.Addr = %s; want log_syslog_addr", Log.Syslog.Addr)
	}
	if Log.Syslog.Protocol != "log_syslog_protocol" {
		t.Fatalf("Log.Syslog.Protocol = %s; want log_syslog_protocol", Log.Syslog.Protocol)
	}
	if Log.Syslog.Tag != "log_syslog_tag" {
		t.Fatalf("Log.Syslog.Tag = %s; want log_syslog_tag", Log.Syslog.Tag)
	}
	if Log.File.Path != "log_file_path" {
		t.Fatalf("Log.File = %s; want log_file_path", Log.File.Path)
	}
	if Log.AWS.LogGroup.Name != "log_aws_log_group_name" {
		t.Fatalf("Log.AWS.LogGroup.Name = %s; want log_aws_log_group_name", Log.AWS.LogGroup.Name)
	}
	if Log.AWS.LogGroup.Region != "log_aws_log_group_region" {
		t.Fatalf("Log.AWS.LogGroup.Region = %s; want log_aws_log_group_region", Log.AWS.LogGroup.Region)
	}
}

func TestLoadPanic(t *testing.T) {
	aws.Client.ParamGet = func(ctx context.Ctx, key string) (string, error) {
		return "", &parameterStoreTypes.ParameterNotFound{}
	}

	aws.Client.SecretGet = func(ctx context.Ctx, key string) (string, error) {
		return "", &secretsManagerTypes.ResourceNotFoundException{}
	}

	// Server required fields.
	os.Setenv("SERVER_ADDR", "")
	testLoadPanic("Server.Addr is nil", t, Load)
	os.Setenv("SERVER_ADDR", "server_addr")

	os.Setenv("SERVER_TIMEOUT_READ", "")
	testLoadPanic("Server.Timeout.Read is nil", t, Load)
	os.Setenv("SERVER_TIMEOUT_READ", "1")

	os.Setenv("SERVER_TIMEOUT_WRITE", "")
	testLoadPanic("Server.Timeout.Write is nil", t, Load)
	os.Setenv("SERVER_TIMEOUT_WRITE", "2")

	os.Setenv("SERVER_TIMEOUT_IDLE", "")
	testLoadPanic("Server.Timeout.Idle is nil", t, Load)
	os.Setenv("SERVER_TIMEOUT_IDLE", "3")

	// Database required fields.
	os.Setenv("DATABASE_HOST", "")
	testLoadPanic("Database.Host is nil", t, Load)
	os.Setenv("DATABASE_HOST", "database_host")

	os.Setenv("DATABASE_PORT", "")
	testLoadPanic("Database.Port is nil", t, Load)
	os.Setenv("DATABASE_PORT", "4")

	os.Setenv("DATABASE_DB", "")
	testLoadPanic("Database.DB is nil", t, Load)
	os.Setenv("DATABASE_DB", "database_db")

	os.Setenv("DATABASE_SSLMODE", "")
	testLoadPanic("Database.SSLMode is nil", t, Load)
	os.Setenv("DATABASE_SSLMODE", "enable")

	os.Setenv("DATABASE_MIGRATIONS_AUTO", "")
	testLoadPanic("Database.Migrations.Auto is nil", t, Load)
	os.Setenv("DATABASE_MIGRATIONS_AUTO", "true")

	os.Setenv("DATABASE_MIGRATIONS_SOURCE", "")
	testLoadPanic("Database.Migrations.Source is nil", t, Load)
	os.Setenv("DATABASE_MIGRATIONS_SOURCE", "database_migrations_source")

	os.Setenv("DATABASE_MIGRATIONS_TYPE", "")
	testLoadPanic("Database.Migrations.Type is nil", t, Load)
	os.Setenv("DATABASE_MIGRATIONS_TYPE", "database_migrations_type")

	// AWS required fields.
	os.Setenv("AWS_S3_BUCKET_DOCUMENTS", "")
	testLoadPanic("AWS.S3.BucketDocuments is nil", t, Load)
	os.Setenv("AWS_S3_BUCKET_DOCUMENTS", "aws_s3_bucket_documents")

	os.Setenv("AWS_SES_SOURCE", "")
	testLoadPanic("AWS.SES.Source is nil", t, Load)
	os.Setenv("AWS_SES_SOURCE", "aws_ses_source")

	// Log required fields.
	os.Setenv("LOG_MODE", "")
	testLoadPanic("Log.Mode is nil", t, Load)
	os.Setenv("LOG_MODE", "log_mode")

	os.Setenv("LOG_LEVEL", "")
	testLoadPanic("Log.Level is nil", t, Load)
	os.Setenv("LOG_LEVEL", "5")
}
