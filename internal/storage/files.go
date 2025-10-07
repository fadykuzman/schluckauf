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
