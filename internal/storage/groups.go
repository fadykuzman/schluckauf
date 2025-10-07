package storage

import (
	_ "database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type GroupStatus string

const (
	StatusPending GroupStatus = "pending"
	StatusDecided GroupStatus = "decided"
)

type Group struct {
	ID           int
	Hash         string
	Size         int64
	FileCount    int
	UpdatedAt    *time.Time
	Status       GroupStatus
	PendingCount int
	DecidedCount int
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

func (s *Storage) ListGroups() ([]Group, error) {
	groupRows, err := s.db.Query("SELECT id, hash, size, file_count, updated_at FROM groups")
	if err != nil {
		return nil, err
	}
	defer groupRows.Close()

	var groups []Group
	for groupRows.Next() {
		var g Group

		if err := groupRows.Scan(
			&g.ID,
			&g.Hash,
			&g.Size,
			&g.FileCount,
			&g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
