package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Ping(ctx context.Context, db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func DisplayMigrations(ctx context.Context, db *pgxpool.Pool) error {
	rows, err := db.Query(ctx, "SELECT * FROM schema_migrations order by version desc")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var version int
		var dirty bool
		if err := rows.Scan(&version, &dirty); err != nil {
			return err
		}
		log.Printf("version: %d | dirty: %t\n", version, dirty)
	}
	return nil
}
