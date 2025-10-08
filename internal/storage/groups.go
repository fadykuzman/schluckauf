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
	ID            int
	Hash          string
	Size          int64
	FileCount     int
	UpdatedAt     *time.Time
	Status        GroupStatus
	ThumbnailPath string
}

type GroupStats struct {
	Pending int
	Decided int
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
		`SELECT g.id, g.hash, g.size, g.file_count, g.updated_at, f.path as thumbnail_path,
			CASE
				WHEN SUM(CASE WHEN f.action = 'pending' THEN 1 ELSE 0 END) > 0
		    THEN 'pending'
		    ELSE 'decided'
		  END as status
		FROM groups g
		LEFT JOIN files f ON g.id = f.group_id
		WHERE f.id = (SELECT MIN(id) FROM files WHERE group_id = g.id)
		GROUP BY g.id
		ORDER BY
		  CASE WHEN status = 'pending' THEN 0 ELSE 1 END,
		  updated_at DESC NULLS LAST
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
			&g.ThumbnailPath,
			&g.Status,
		); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (s *Storage) GetGroupStats() (GroupStats, error) {
	rows, err := s.db.Query(`
		SELECT status, COUNT(*) as count
		FROM (
		  SELECT g.id,
				CASE
		      WHEN SUM(CASE WHEN f.action = 'pending' THEN 1 ELSE 0 END) > 0
		      THEN 'pending'
		      ELSE 'decided'
	 			END as status
			FROM groups g
		  LEFT JOIN files f ON g.id = f.group_id
		  GROUP BY g.id
		) as group_statuses
		GROUP BY status
		`)
	if err != nil {
		return GroupStats{}, err
	}

	defer rows.Close()

	var gs GroupStats

	for rows.Next() {
		var status string
		var count int

		if err := rows.Scan(&status, &count); err != nil {
			return GroupStats{}, err
		}

		switch status {
		case "pending":
			gs.Pending = count
		case "decided":
			gs.Decided = count
		}

	}

	return gs, nil
}
