package storage

import (
	"encoding/json"
	"time"

	_ "modernc.org/sqlite"
)

type GroupStatus string

const (
	StatusPending  GroupStatus = "pending"
	StatusDecided  GroupStatus = "decided"
	StatusArchived GroupStatus = "archived"
)

type ImageGroup struct {
	ID            int         `json:"id"`
	Hash          string      `json:"hash"`
	Size          int64       `json:"size"`
	ImageCount    int         `json:"imageCount"`
	UpdatedAt     *time.Time  `json:"updatedAt"`
	Status        GroupStatus `json:"status"`
	ThumbnailPath string      `json:"thumbnailPath"`
}

type ImageGroupStats struct {
	Pending            int `json:"pending"`
	Decided            int `json:"decided"`
	ImagesToTrashCount int `json:"imagesToTrashCount"`
}

func (s *Storage) CreateImageGroup(hash []int, size int64, fileCount int) (int, error) {
	hashJSON, _ := json.Marshal(hash)

	result, err := s.db.Exec(
		"INSERT INTO image_groups (hash, size, image_count) VALUES (?, ?, ?)",
		hashJSON, size, fileCount,
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

func (s *Storage) ListImageGroups() ([]ImageGroup, error) {
	groupRows, err := s.db.Query(
		`SELECT g.id, g.hash, g.size, g.image_count, g.updated_at, i.path as thumbnail_path,
			CASE
				WHEN SUM(CASE WHEN i.action = 'pending' THEN 1 ELSE 0 END) > 0
		    THEN 'pending'
		    ELSE 'decided'
		  END as status
		FROM image_groups g
		LEFT JOIN images i ON g.id = i.group_id
		WHERE i.id = (SELECT MIN(id) FROM images WHERE group_id = g.id)
		GROUP BY g.id
		ORDER BY
		  CASE WHEN status = 'pending' OR status = 'trashed' THEN 0 ELSE 1 END,
		  updated_at DESC NULLS LAST
		`)
	if err != nil {
		return nil, err
	}
	defer groupRows.Close()

	var groups []ImageGroup
	for groupRows.Next() {
		var g ImageGroup

		if err := groupRows.Scan(
			&g.ID,
			&g.Hash,
			&g.Size,
			&g.ImageCount,
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

func (s *Storage) GetImageGroupStats() (ImageGroupStats, error) {
	rows, err := s.db.Query(`
		SELECT status, COUNT(*) as count
		FROM (
		  SELECT g.id,
				CASE
		      WHEN SUM(CASE WHEN i.action = 'pending' THEN 1 ELSE 0 END) > 0
		      THEN 'pending'
		      ELSE 'decided'
	 			END as status
			FROM image_groups g
		  LEFT JOIN images i ON g.id = i.group_id
		  GROUP BY g.id
		) as group_statuses
		GROUP BY status
		`)
	if err != nil {
		return ImageGroupStats{}, err
	}

	defer rows.Close()

	var gs ImageGroupStats

	for rows.Next() {
		var status string
		var count int

		if err := rows.Scan(&status, &count); err != nil {
			return ImageGroupStats{}, err
		}

		switch status {
		case "pending":
			gs.Pending = count
		case "decided":
			gs.Decided = count
		}
	}

	row := s.db.QueryRow("SELECT COUNT(*) FROM images WHERE action = 'trash'")

	if err := row.Scan(&gs.ImagesToTrashCount); err != nil {
		return gs, err
	}

	return gs, nil
}
