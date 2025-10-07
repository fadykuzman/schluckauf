package storage

import _ "database/sql"

type FileAction string

const (
	ActionPending FileAction = "pending"
	ActionKeep    FileAction = "keep"
	ActionTrash   FileAction = "trash"
)

type File struct {
	ID       int
	GroupID  int
	Path     string
	Filesize int64
	Action   FileAction
}

type FilesStats struct {
	TotalPending int
	TotalKeep    int
	TotalTrash   int
}

func (s *Storage) CreateFile(groupID int, path string, filesize int64) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO files (group_id, path, filesize) VALUES (?, ?,?)",
		groupID, path, filesize,
	)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return int(id), err
}

func (s *Storage) GetGroupFiles(groupID int) ([]File, error) {
	rows, err := s.db.Query(
		"SELECT id, group_id, path, filesize, action FROM files WHERE group_id=?",
		groupID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var files []File

	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.GroupID, &f.Path, &f.Filesize, &f.Action); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (s *Storage) UpdateFileAction(groupID int, fileID int, action string) error {
	_, err := s.db.Exec(
		"UPDATE files SET action = ? WHERE id = ?",
		action, fileID,
	)

	_, errGroup := s.db.Exec(
		" UPDATE groups SET updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		groupID,
	)

	if err != nil {
		return err
	}
	return errGroup
}

func (s *Storage) getFileStats(groupID int, fileID int, action string) (FilesStats, error) {
	return FilesStats{}, nil
}
