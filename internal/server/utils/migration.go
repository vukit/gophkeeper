package utils

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/vukit/gophkeeper/internal/server/logger"

	// Register packages for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationUp производит миграции для БД
func MigrationUp(source, dsn string, mLogger *logger.Logger) {
	m, err := migrate.New(source, dsn)
	if err != nil {
		mLogger.Fatal(err.Error())

		return
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		mLogger.Fatal(err.Error())
	}
}

// MigrationDown отменяет миграции для БД
func MigrationDown(source, dsn string, mLogger *logger.Logger) {
	m, err := migrate.New(source, dsn)
	if err != nil {
		mLogger.Fatal(err.Error())

		return
	}

	if err := m.Down(); err != nil {
		mLogger.Fatal(err.Error())
	}
}
