package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shikihtm/blog-backend/internal/logger"
)

func Initialize(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error when open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error when ping sqlite: %w", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to migrate tables: %w", err)
	}

	db.SetMaxOpenConns(1)
	log.Println(logger.System("Database connection established and tables validated successfully."))
	return db, nil
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS post_stats (
		uuid TEXT PRIMARY KEY,
		slug TEXT UNIQUE NOT NULL,
		views INTEGER DEFAULT 0,
		likes INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}
