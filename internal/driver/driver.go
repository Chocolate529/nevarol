package driver

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// ConnectSQL creates database connection pool
func ConnectSQL(dsn string) (*DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %v", err)
	}

	// Set connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	log.Println("Database connection established")

	return &DB{Pool: pool}, nil
}

// RunMigrations runs SQL migration files
func (db *DB) RunMigrations() error {
	migrationFile := "./migrations/001_init.sql"

	// Read migration file
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("unable to read migration file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute migration
	_, err = db.Pool.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("unable to execute migration: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
