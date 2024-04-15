package service

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/database"
	"github.com/koraygocmen/golang-boilerplate/internal/database/databasetest"
	"github.com/koraygocmen/golang-boilerplate/pkg/str"
)

func TestNew(t *testing.T) {
	services, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("want: read service dir error nil; got: %v", err)
	}

	dbtest := databasetest.Get()
	database.DB = dbtest.DB

	tx := Service.Transaction(context.Background(), 30*time.Second)
	txReflect := reflect.ValueOf(tx)
	for _, service := range services {
		if service.IsDir() {
			serviceName := str.SnakeToCamel(service.Name())
			r := reflect.Indirect(txReflect).FieldByName(serviceName)
			if r.Kind() == 0 || r.IsZero() || r.IsNil() {
				t.Fatalf("want: service tx %s initialized; got: nil", serviceName)
			}
		}
	}
}
