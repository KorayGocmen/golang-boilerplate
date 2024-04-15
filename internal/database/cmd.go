package database

import (
	"fmt"

	"github.com/koraygocmen/golang-boilerplate/internal/config"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
	"github.com/koraygocmen/golang-boilerplate/migrations"
	"github.com/koraygocmen/golang-boilerplate/seeds"
	"github.com/pressly/goose/v3"
)

func (d *Database) List() error {
	migrationFiles, err := migrations.FS.ReadDir(migrations.Path)
	if err != nil {
		err = fmt.Errorf("read seeds source error: %w", err)
		return err
	}

	for _, migrationFile := range migrationFiles {
		fmt.Println(migrationFile.Name())
	}

	return nil
}

func (d *Database) Status() error {
	if err := goose.Status(d.SQL, migrations.Path); err != nil {
		err = fmt.Errorf("goose status error: %w", err)
		return err
	}

	return nil
}

func (d *Database) Create(name string) error {
	var (
		source = config.Database.Migrations.Source
		typ    = config.Database.Migrations.Type
	)

	if err := goose.Create(d.SQL, source, name, typ); err != nil {
		err = fmt.Errorf("goose create error: %w", err)
		return err
	}

	return nil
}

func (d *Database) Up(ctx context.Ctx) error {
	if err := goose.Up(d.SQL, migrations.Path); err != nil {
		err = fmt.Errorf("goose up error: %w", err)
		return err
	}

	return nil
}

func (d *Database) UpByOne(ctx context.Ctx) error {
	if err := goose.UpByOne(d.SQL, migrations.Path); err != nil {
		err = fmt.Errorf("goose up-by-one error: %w", err)
		return err
	}

	return nil
}

func (d *Database) UpTo(ctx context.Ctx, version int64) error {
	if err := goose.UpTo(d.SQL, migrations.Path, version); err != nil {
		err = fmt.Errorf("goose up-to error: %w", err)
		return err
	}

	return nil
}

func (d *Database) Down(ctx context.Ctx) error {
	if err := goose.Down(d.SQL, migrations.Path); err != nil {
		err = fmt.Errorf("goose down error: %w", err)
		return err
	}

	return nil
}

func (d *Database) DownTo(ctx context.Ctx, version int64) error {
	if err := goose.DownTo(d.SQL, migrations.Path, version); err != nil {
		err = fmt.Errorf("goose down-to error: %w", err)
		return err
	}

	return nil
}

func (d *Database) Reset(ctx context.Ctx) error {
	if env.IsProd() {
		return fmt.Errorf("goose reset is not allowed in production environment")
	}

	if err := goose.Reset(d.SQL, migrations.Path); err != nil {
		err = fmt.Errorf("goose reset error: %w", err)
		return err
	}

	return nil
}

func (d *Database) Seed(ctx context.Ctx) error {
	seedFiles, err := seeds.FS.ReadDir(seeds.Path)
	if err != nil {
		err = fmt.Errorf("read seeds source error: %w", err)
		return err
	}

	for _, seedFile := range seedFiles {
		seedFilePath := fmt.Sprintf("%s/%s", seeds.Path, seedFile.Name())
		seedsData, err := seeds.FS.ReadFile(seedFilePath)
		if err != nil {
			err = fmt.Errorf("read seeds source file error: %w", err)
			return err
		}

		if err := d.GORM.Exec(string(seedsData)).Error; err != nil {
			err = fmt.Errorf("exec seeds source file error: %w", err)
			return err
		}
	}

	return nil
}
