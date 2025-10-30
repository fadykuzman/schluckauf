// Package storage provides SQLite-based persistence for duplicate photo groups and files
package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Database: %w", err)
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to set journal mode: %w", err)
	}

	_, err = db.Exec("PRAGMA busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to set busy_timeout: %w", err)
	}
	_, err = db.Exec(`
		  CREATE TABLE IF NOT EXISTS file_groups (
		      id INTEGER PRIMARY KEY AUTOINCREMENT,
		      hash TEXT UNIQUE,
		      size INTEGER,
					file_count INTEGER,
		      updated_at TIMESTAMP NULL
			);
		  CREATE TABLE IF NOT EXISTS files (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
				  group_id INTEGER,
		      path TEXT,
					filesize INTEGER,
		      action TEXT DEFAULT 'pending',
		      FOREIGN KEY(group_id) REFERENCES file_groups(id)
		  );
		  CREATE INDEX IF NOT EXISTS idx_file_group_action ON files(group_id, action);
		  
		  CREATE TABLE IF NOT EXISTS image_groups (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				hash TEXT,
				size INTEGER,
				image_count INTEGER,
				updated_at TIMESTAMP NULL
			);
		  CREATE TABLE IF NOT EXISTS images (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				group_id INTEGER,
				path TEXT,
				image_size INTEGER,
				width INTEGER,
				height INTEGER,
				similarity REAL,
				action TEXT DEFAULT 'pending',
				FOREIGN KEY(group_id) REFERENCES image_groups(id)
		  );
	    CREATE INDEX IF NOT EXISTS idx_image_group_action ON images(group_id, action);
		`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
