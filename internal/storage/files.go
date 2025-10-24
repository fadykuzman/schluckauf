package storage

import (
	_ "database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileAction string

const (
	ActionPending FileAction = "pending"
	ActionKeep    FileAction = "keep"
	ActionTrash   FileAction = "trash"
	ActionTrashed FileAction = "trashed"
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
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE files SET action = ? WHERE id = ?",
		action, fileID,
	)
	if err != nil {
		return err
	}

	_, errGroup := tx.Exec(
		" UPDATE groups SET updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		groupID,
	)

	if errGroup != nil {
		return errGroup
	}

	return tx.Commit()
}

type FileToTrash struct {
	ID       int
	FilePath string
}

type TrashFilesResponse struct {
	MovedCount      int
	FailedCount     int
	PartialFailures int
	TotalCount      int
	Errors          []string
}

func (s *Storage) TrashFiles() (TrashFilesResponse, error) {
	rows, err := s.db.Query(`
		SELECT id, path FROM files WHERE action = 'trash'
		`)
	if err != nil {
		return TrashFilesResponse{}, err
	}
	defer rows.Close()

	var filesToTrash []FileToTrash

	for rows.Next() {
		var f FileToTrash
		if err := rows.Scan(&f.ID, &f.FilePath); err != nil {
			return TrashFilesResponse{}, err
		}
		filesToTrash = append(filesToTrash, f)
	}

	log.Print("Moving files to trash")

	var movedCount int
	var failedCount int
	var partialFailures int
	var errors []string

	timestamp := time.Now().Format("2006-01-02_15-04-05")

	for _, f := range filesToTrash {
		log.Printf("Moving file %d to trash", f.ID)
		destPath, err := moveFileToTrash(f, timestamp)

		if err != nil {
			log.Printf("Error moving file %d to trash", f.ID)
			errors = append(errors, fmt.Sprintf("Couldn't move file %s to trash. %s", f.FilePath, err))
			failedCount++
		} else {
			log.Printf("Moved file %d to trash", f.ID)
			err := s.updateDBForTrashedFile(f.ID, destPath)

			if err != nil {
				errors = append(errors, fmt.Sprintf("File moved but couldn't update database for file %s: %s", f.FilePath, err))
				partialFailures++
			} else {
				movedCount++
			}
		}
	}

	response := TrashFilesResponse{
		MovedCount:      movedCount,
		FailedCount:     failedCount,
		PartialFailures: partialFailures,
		TotalCount:      movedCount + failedCount + partialFailures,
		Errors:          errors,
	}

	return response, nil
}

func (s *Storage) updateDBForTrashedFile(fileID int, newPath string) error {
	log.Printf("Updating file %d to be %s", fileID, ActionTrashed)
	result, err := s.db.Exec(`
				UPDATE files 
				SET action = ?, 
				path = ? 
				WHERE id = ?`,
		ActionTrashed,
		newPath,
		fileID,
	)

	fmt.Printf("%s", result)
	if err != nil {
		log.Printf("Error while Updating file %d to be %s", fileID, ActionTrashed)
		return err
	}
	log.Printf("Updated file %d to be %s", fileID, ActionTrashed)

	return nil
}

func moveFileToTrash(f FileToTrash, timestamp string) (string, error) {
	destPath := filepath.Join("./trash", timestamp, f.FilePath)
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return "", err
	}

	if err := os.Rename(f.FilePath, destPath); err != nil {
		return "", err
	}

	return destPath, nil
}

func (s *Storage) DeletePendingData() error {
	_, err := s.db.Exec("DELETE FROM files WHERE action = 'pending")
	if err != nil {
		return err
	}

	_, err = s.db.Exec("DELETE FROM groups WHERE id NOT IN (SELECT DISTINCT group_id FROM files)")
	if err != nil {
		return err
	}
	return nil
}
