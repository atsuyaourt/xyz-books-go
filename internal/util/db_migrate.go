package util

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func DBMigrationUp(migrationPath string, dbSource string) error {
	migration, err := migrate.New("file://"+migrationPath, dbSource)
	if err != nil {
		return err
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func DBMigrationDown(migrationPath string, dbSource string) error {
	migration, err := migrate.New("file://"+migrationPath, dbSource)
	if err != nil {
		return err
	}

	if err = migration.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
