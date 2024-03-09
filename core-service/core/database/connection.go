package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alekLukanen/go-templ-htmx-example-app/core/settings"
)

func ConfigurePostgresDBConnection(ctx context.Context) (*pgxpool.Pool, error) {

	internalCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	connectionSettings := DatabaseConnectionURL("postgres")
	config, err := pgxpool.ParseConfig(connectionSettings)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 4

	db, err := pgxpool.NewWithConfig(internalCtx, config)
	if err != nil {
		return nil, err
	}

	if err := Ping(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func DatabaseConnectionURL(driverName string) string {
	connectionSettings := fmt.Sprintf(
		"%s://%s:%s@%s:5432/%s?sslmode=disable",
		driverName,
		settings.DATABASE_DB_USER,
		settings.DATABASE_DB_PASSWORD,
		settings.DATABASE_DB_HOST,
		settings.DATABASE_DB_NAME,
	)
	return connectionSettings
}
