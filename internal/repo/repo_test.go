package repo

import (
	"os"
	"reflect"
	"testing"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/database"
	"github.com/koraygocmen/golang-boilerplate/internal/database/databasetest"
	"github.com/koraygocmen/golang-boilerplate/pkg/str"
)

func TestNew(t *testing.T) {
	repos, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("want: read repo dir error nil; got: %v", err)
	}

	dbtest := databasetest.Get()
	database.DB = dbtest.DB

	tx := New(context.Background())
	txReflect := reflect.ValueOf(tx)
	for _, repo := range repos {
		if repo.IsDir() {
			repoName := str.SnakeToCamel(repo.Name())
			r := reflect.Indirect(txReflect).FieldByName(repoName)
			if r.Kind() == 0 || r.IsZero() || r.IsNil() {
				t.Fatalf("want: repo tx %s initialized; got: nil", repoName)
			}
		}
	}
}
