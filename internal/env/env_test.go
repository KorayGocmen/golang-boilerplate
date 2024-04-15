package env

import (
	"os"
	"testing"
)

func testLoadPanic(testcase string, t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Fatalf("%s did not panic", testcase)
}

func setRequiredFields() {
	os.Setenv("ENV", "TEST")

	os.Setenv("AWS_REGION", "aws_region")
	os.Setenv("AWS_ACCESS_KEY_ID", "aws_access_key_id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "aws_secret_access_key")
}

func TestInit(t *testing.T) {
	setRequiredFields()

	Init()

	// Env.
	if ENV != "TEST" {
		t.Fatalf("Env = %s; want TEST", ENV)
	}

	// AWS.
	if AWS_REGION != "aws_region" {
		t.Fatalf("AWS_REGION = %s; want aws_region", AWS_REGION)
	}
	if AWS_ACCESS_KEY_ID != "aws_access_key_id" {
		t.Fatalf("AWS_ACCESS_KEY_ID= %s; want aws_access_key_id", AWS_ACCESS_KEY_ID)
	}
	if AWS_SECRET_ACCESS_KEY != "aws_secret_access_key" {
		t.Fatalf("AWS_SECRET_ACCESS_KEY = %s; want aws_secret_access_key", AWS_SECRET_ACCESS_KEY)
	}
}

func TestInitPanic(t *testing.T) {
	os.Setenv("ENV", "")
	testLoadPanic("Env is nil", t, Init)
	os.Setenv("ENV", "TEST")

	// AWS required fields.
	os.Setenv("AWS_REGION", "")
	testLoadPanic("AWS.Region is nil", t, Init)
	os.Setenv("AWS_REGION", "aws_region")
}

func TestIsDev(t *testing.T) {
	setRequiredFields()

	// Override ENV.
	os.Setenv("ENV", "DEV")

	if Init(); !IsDev() {
		t.Fatalf(`Env = "DEV"; want true`)
	}

	os.Setenv("ENV", "development")
	if Init(); !IsDev() {
		t.Fatalf(`Env = "development"; want true`)
	}

	os.Setenv("ENV", "test")
	if Init(); IsDev() {
		t.Fatalf(`Env = "test"; want false`)
	}
}

func TestIsProd(t *testing.T) {
	setRequiredFields()

	// Override ENV.
	os.Setenv("ENV", "PROD")

	if Init(); !IsProd() {
		t.Fatalf(`Env = "prod"; want true`)
	}

	os.Setenv("ENV", "production")
	if Init(); !IsProd() {
		t.Fatalf(`Env = "production"; want true`)
	}

	os.Setenv("ENV", "test")
	if Init(); IsProd() {
		t.Fatalf(`Env = "test"; want false`)
	}
}
