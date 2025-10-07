// Package storage provides SQLite-based persistence for duplicate photo groups and files
package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		  CREATE TABLE IF NOT EXISTS groups (
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
		      FOREIGN KEY(group_id) REFERENCES groups(id)
		  );
		  CREATE INDEX IF NOT EXISTS idx_group_action ON files(group_id, action);
		`)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
