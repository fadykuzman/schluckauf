// Package storage provides SQLite-based persistence for duplicate photo groups and files
package storage

import (
	"database/sql"
	"log"

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

func (s *Storage) InsertSampleData() error {
	gid, err1 := s.CreateGroup("sample-hash-1", 2156432, 3)
	s.CreateFile(gid, "/photos/IMG_1234.jpg", 2156432)
	s.CreateFile(gid, "/photos/backup/IMG_1234.jpg", 2156432)
	s.CreateFile(gid, "/photos/copy/IMG_1234.jpg", 1847000)

	if err1 != nil {
		log.Fatal("Could not insert Group 1 in sample data", err1)
	}

	gid2, err2 := s.CreateGroup("sample-hash-2", 3500000, 2)
	s.CreateFile(gid2, "/photos/IMG_5678.jpg", 3500000)
	s.CreateFile(gid2, "/photos/old/IMG_5678.jpg", 3500000)

	if err2 != nil {
		log.Fatal("Could not insert Group 2 in sample data", err2)
	}
	return nil
}
