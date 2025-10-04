package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Group struct {
	ID        int
	Hash      string
	Size      int64
	FileCount int
}

type File struct {
	ID         int
	GroupID    int
	Path       string
	Resolution string
	Filesize   int64
	Action     string
}

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
					file_count INTEGER);
		  CREATE TABLE IF NOT EXISTS files (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
				  group_id INTEGER,
		      path TEXT,
		      resolution TEXT,
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

func (s *Storage) CreateGroup(hash string, size int64, fileCount int) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO groups (hash, size, file_count) VALUES (?, ?, ?)",
		hash, size, fileCount,
	)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (s *Storage) CreateFile(groupID int, path string, filesize int64) error {
	_, err := s.db.Exec(
		"INSERT INTO files (group_id, path, filesize) VALUES (?, ?,?)",
		groupID, path, filesize,
	)
	return err
}

func (s *Storage) ListGroups() ([]Group, error) {
	rows, err := s.db.Query("SELECT id, hash, size, file_count FROM groups")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Hash, &g.Size, &g.FileCount); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (s *Storage) InsertSampleData() error {
	gid, _ := s.CreateGroup("sample-hash-1", 2156432, 3)
	s.CreateFile(gid, "/photos/IMG_1234.jpg", 2156432)
	s.CreateFile(gid, "/photos/backup/IMG_1234.jpg", 2156432)
	s.CreateFile(gid, "/photos/copy/IMG_1234.jpg", 1847000)

	gid2, _ := s.CreateGroup("sample-hash-2", 3500000, 2)
	s.CreateFile(gid2, "/photos/IMG_5678.jpg", 3500000)
	s.CreateFile(gid2, "/photos/old/IMG_5678.jpg", 3500000)

	return nil
}
