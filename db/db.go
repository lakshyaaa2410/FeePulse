// db/db.go
package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBInterface interface {
	GetPool() *pgxpool.Pool
	Close()
}

type DB struct {
	pool *pgxpool.Pool
}

func InitializeDB() (*DB, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Close() {
	db.pool.Close()
}
