// package db

package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type DB interface {
	GetConnection() *sql.DB
	Close() error
}

type shopDb struct {
	pgServer *sql.DB
}

func (db *shopDb) GetConnection() *sql.DB {
	return db.pgServer
}

func (db *shopDb) Close() error {
	return db.pgServer.Close()
}

func NewShopDb(connString string) (DB, error) {

	pgDb, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure the connection pool
	pgDb.SetMaxOpenConns(25)
	pgDb.SetMaxIdleConns(25)
	pgDb.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := pgDb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully.")
	return &shopDb{pgServer: pgDb}, nil
}
