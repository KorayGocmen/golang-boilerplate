package config

import (
	"testing"
)

func TestParse(t *testing.T) {
	dbConf := DatabaseConfig{
		Host:    "localhost",
		Port:    "5432",
		DB:      "db",
		User:    "user",
		Pass:    "pass",
		SSLMode: "disable",
	}
	want := "postgres://user:pass@localhost:5432/db?sslmode=disable"
	if got := SqlOpts(dbConf); got != want {
		t.Fatalf("want: %s; got: %s", want, got)
	}

	dbConf = DatabaseConfig{
		Host:    "localhost",
		Port:    "5432",
		DB:      "db",
		User:    "user",
		SSLMode: "disable",
	}
	wantMissingPass := "postgres://user@localhost:5432/db?sslmode=disable"
	if gotMissingPass := SqlOpts(dbConf); gotMissingPass != wantMissingPass {
		t.Fatalf("want: %s; got: %s", wantMissingPass, gotMissingPass)
	}

	dbConf = DatabaseConfig{
		Host: "localhost",
		Port: "5432",
		DB:   "db",
		User: "user",
		Pass: "pass",
	}
	wantMissingSSLMode := "postgres://user:pass@localhost:5432/db"
	if gotMissingSSLMode := SqlOpts(dbConf); gotMissingSSLMode != wantMissingSSLMode {
		t.Fatalf("want: %s; got: %s", wantMissingSSLMode, gotMissingSSLMode)
	}
}
