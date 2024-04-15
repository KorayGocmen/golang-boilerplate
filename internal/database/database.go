package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/koraygocmen/golang-boilerplate/internal/config"
	"github.com/koraygocmen/golang-boilerplate/internal/logger"
	"github.com/koraygocmen/golang-boilerplate/migrations"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database struct {
	GORM *gorm.DB
	SQL  *sql.DB
}

var (
	DB *Database
)

func Connect(loggr *logger.Writer, conf config.DatabaseConfig) (*Database, error) {
	gormLogLevel := config.Log.Level
	if gormLogLevel > int(gormlogger.Info) {
		gormLogLevel = int(gormlogger.Info)
	}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         loggr,
	}

	DBGORM, err := gorm.Open(postgres.Open(config.SqlOpts(conf)), gormConfig)
	if err != nil {
		err = fmt.Errorf("gorm open db error: %w", err)
		return nil, err
	}

	DBSQL, err := DBGORM.DB()
	if err != nil {
		err = fmt.Errorf("gorm get db error: %w", err)
		return nil, err
	}

	// Set up connection pool.
	DBSQL.SetMaxIdleConns(conf.MaxIdleConns)
	DBSQL.SetMaxOpenConns(conf.MaxOpenConns)
	DBSQL.SetConnMaxIdleTime(time.Duration(conf.MaxIdleTimeMs) * time.Millisecond)
	DBSQL.SetConnMaxLifetime(time.Duration(conf.MaxLifetimeMs) * time.Millisecond)

	// Set up migrations via goose.
	goose.SetLogger(loggr)
	goose.SetDialect("postgres")
	goose.SetBaseFS(migrations.FS)

	return &Database{
		GORM: DBGORM,
		SQL:  DBSQL,
	}, nil
}
