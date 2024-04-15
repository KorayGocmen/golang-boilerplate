package databasetest

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/koraygocmen/golang-boilerplate/internal/config"
	"github.com/koraygocmen/golang-boilerplate/internal/database"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
	"github.com/koraygocmen/golang-boilerplate/internal/logger"
	"github.com/koraygocmen/golang-boilerplate/pkg/generate"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

type DatabaseTest struct {
	DB       *database.Database
	Pool     *dockertest.Pool
	Resource *dockertest.Resource
}

var (
	dbOnce sync.Once
	db     *DatabaseTest
)

// Get is a singleton for the database test to either
// initialize or return the existing database test.
func Get() *DatabaseTest {
	dbOnce.Do(func() {
		var err error
		if db, err = create(); err != nil {
			log.Fatalf("test main error: %v", err)
		}
	})
	return db
}

func create() (*DatabaseTest, error) {
	dbTest := &DatabaseTest{}

	var err error
	dbTest.Pool, err = dockertest.NewPool("")
	if err != nil {
		err = fmt.Errorf("database test create error: dockertest new pool error: %w", err)
		return nil, err
	}

	// Default max wait is 30 seconds.
	dbTest.Pool.MaxWait = time.Duration(30) * time.Second

	dbConfigTest := config.DatabaseConfig{
		DB:            "root",
		User:          "root",
		Pass:          "pass",
		SSLMode:       "disable",
		MaxIdleConns:  10,
		MaxOpenConns:  10,
		MaxIdleTimeMs: 0,
		MaxLifetimeMs: 300000,
	}

	// Run the database.
	dbTest.Resource, err = dbTest.Pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       fmt.Sprintf("test-%s", generate.AlphaCode(10, false)),
			Repository: "postgres",
			Tag:        "14.6-alpine",
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", dbConfigTest.User),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", dbConfigTest.Pass),
				fmt.Sprintf("POSTGRES_DB=%s", dbConfigTest.DB),
			},
		})
	if err != nil {
		err = fmt.Errorf("database test create error: dockertest pool run error: %w", err)
		return nil, err
	}

	// Set the resource host and port.
	resourceHostPort := dbTest.Resource.GetHostPort("5432/tcp")
	dbConfigTest.Host, dbConfigTest.Port, _ = net.SplitHostPort(resourceHostPort)

	// Exponential backoff-retry, since the application might not be ready right away.
	if err := dbTest.Pool.Retry(func() error {
		db, err := sql.Open("postgres", config.SqlOpts(dbConfigTest))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		dbTest.Purge()
		err = fmt.Errorf("database test create error: database ping error: %w", err)
		return nil, err
	}

	loggr, err := logger.New(logger.Config{
		Mode: string(logger.ModeNone),
	})
	if err != nil {
		dbTest.Purge()
		err = fmt.Errorf("database test create error: logger new error: %w", err)
		return nil, err
	}

	if dbTest.DB, err = database.Connect(loggr, dbConfigTest); err != nil {
		dbTest.Purge()
		err = fmt.Errorf("database test create error: %w", err)
		return nil, err
	}

	// Enable full query logging if in development environment.
	if env.IsDev() {
		_, err = dbTest.DB.SQL.Query(fmt.Sprintf("ALTER DATABASE %s SET log_statement = 'all';", dbConfigTest.DB))
		if err != nil {
			dbTest.Purge()
			err = fmt.Errorf("database test create error: %w", err)
			return nil, err
		}
	}

	return dbTest, nil
}

func (d *DatabaseTest) Purge() error {
	if err := d.Pool.Purge(d.Resource); err != nil {
		err = fmt.Errorf("database test purge error: dockertest pool purge error: %w", err)
		return err
	}
	return nil
}
