package storage

import (
	"time"

	_ "modernc.org/sqlite"
)

type GroupStatus string

const (
	StatusPending GroupStatus = "pending"
	StatusDecided GroupStatus = "decided"
)

type Group struct {
	ID        int
	Hash      string
	Size      int64
	FileCount int
	UpdatedAt *time.Time
	Status    GroupStatus
}

func (s *Storage) CreateGroup(hash string, size int64, fileCount int) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO groups (hash, size, file_count) VALUES (?, ?, ?)",
		hash, size, fileCount,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *Storage) ListGroups() ([]Group, error) {
	groupRows, err := s.db.Query(
		`SELECT g.id, g.hash, g.size, g.file_count, g.updated_at 
			CASE
				WHEN SUM(CASE WHEN action = 'pending' THEN 1 ELSE 0 END) > 0
		    THEN 'pending'
		    ELSE 'decided'
		  END as status
		FROM groups g
		LEFT JOIN files f ON g.id = f.group_id
		GROUP BY g.id
		`)
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
			&g.UpdatedAt,
			&g.Status,
		); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
