package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
	"github.com/rakunlabs/muz"
)

var (
	ConnMaxLifetime = 15 * time.Minute
	MaxIdleConns    = 3
	MaxOpenConns    = 3
)

// Connect attempts to connect to database server.
func Connect(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", "postgres://postgres@localhost:5432/postgres")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)

	return db, nil
}

func MigrateDB(ctx context.Context, db *sqlx.DB) error {
	m := muz.Migrate{
		Path:      "migrations",
		Extension: ".sql",
	}

	driver := &muz.SQLDriver{
		DB:      db.DB,
		Dialect: muz.DialectPostgres,
		Table:   "migrations",
		Logger:  slog.Default(),
	}

	if err := m.Migrate(ctx, driver); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
