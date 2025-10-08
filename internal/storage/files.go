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

func (s *Storage) CreateFile(groupID int, path string, filesize int64) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO files (group_id, path, filesize) VALUES (?, ?,?)",
		groupID, path, filesize,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
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

func (s *Storage) UpdateFileAction(groupID int, fileID int, action FileAction) error {
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

type FileToTrash struct {
	ID int
}

type FilesToTrash struct {
	Files []FileToTrash
}

func (s *Storage) TrashFiles() (FilesToTrash, error) {
	rows, err := s.db.Query(`
		SELECT id FROM files WHERE action = 'trash'
		`)
	if err != nil {
		return FilesToTrash{}, err
	}

	defer rows.Close()

	var files []FileToTrash

	for rows.Next() {
		var f FileToTrash
		if err := rows.Scan(&f.ID); err != nil {
			return FilesToTrash{}, err
		}
		files = append(files, f)
	}

	return FilesToTrash{files}, nil
}
