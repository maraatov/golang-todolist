package db

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(connectStr string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), connectStr)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New("file://internal/db/migrations", connectStr)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return dbPool, nil
}
