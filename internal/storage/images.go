package storage

import (
	_ "database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ImageAction string

type Image struct {
	ID        int         `json:"id"`
	GroupID   int         `json:"groupId"`
	Path      string      `json:"path"`
	Imagesize int64       `json:"imageSize"`
	Action    ImageAction `json:"action"`
}

type ImageToTrash struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}

type TrashImagesResponse struct {
	MovedCount      int      `json:"movedCount"`
	FailedCount     int      `json:"failedCount"`
	PartialFailures int      `json:"partialfailures"`
	TotalCount      int      `json:"totalCount"`
	Errors          []string `json:"errors"`
}

func (s *Storage) CreateImage(groupID int, path string, filesize int64) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO images (group_id, path, image_size) VALUES (?, ?,?)",
		groupID, path, filesize,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert image: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insertId for image: %w", err)
	}
	return int(id), nil
}

func (s *Storage) GetGroupImages(groupID int) ([]Image, error) {
	rows, err := s.db.Query(
		`SELECT id, group_id, path, image_size, action 
						FROM images 
						WHERE group_id=?
						AND action != 'trashed'`,
		groupID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var images []Image

	for rows.Next() {
		var f Image
		if err := rows.Scan(&f.ID, &f.GroupID, &f.Path, &f.Imagesize, &f.Action); err != nil {
			return nil, err
		}
		images = append(images, f)
	}

	return images, nil
}

func (s *Storage) UpdateImageAction(groupID int, fileID int, action ImageAction) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE images SET action = ? WHERE id = ?",
		action, fileID,
	)
	if err != nil {
		return err
	}

	_, errGroup := tx.Exec(
		" UPDATE image_groups SET updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		groupID,
	)

	if errGroup != nil {
		return errGroup
	}

	return tx.Commit()
}

func (s *Storage) TrashImages() (TrashImagesResponse, error) {
	rows, err := s.db.Query(`
		SELECT id, path FROM images WHERE action = 'trash'
		`)
	if err != nil {
		return TrashImagesResponse{}, fmt.Errorf("failed to query images to trash: %w", err)
	}
	defer rows.Close()

	var imagesToTrash []ImageToTrash

	for rows.Next() {
		var image ImageToTrash
		if err := rows.Scan(&image.ID, &image.Path); err != nil {
			return TrashImagesResponse{}, fmt.Errorf("failed to scan image to trash row into ImageToTrash struct (%w)", err)
		}
		imagesToTrash = append(imagesToTrash, image)
	}

	log.Print("Moving files to trash")

	var movedCount int
	var failedCount int
	var partialFailures int
	var errors []string

	timestamp := time.Now().Format("2006-01-02_15-04-05")

	for _, image := range imagesToTrash {
		log.Printf("Moving file %d to trash", image.ID)
		destPath, err := moveImageToTrash(image, timestamp)

		if err != nil {
			log.Printf("Error moving file %d to trash", image.ID)
			errors = append(errors, fmt.Sprintf("Couldn't move file %s to trash. %s", image.Path, err))
			failedCount++
		} else {
			log.Printf("Moved file %d to trash", image.ID)
			err := s.updateDBForTrashedImage(image.ID, destPath)

			if err != nil {
				errors = append(errors, fmt.Sprintf("File moved but couldn't update database for file %s: %s", image.Path, err))
				partialFailures++
			} else {
				movedCount++
			}
		}
	}

	response := TrashImagesResponse{
		MovedCount:      movedCount,
		FailedCount:     failedCount,
		PartialFailures: partialFailures,
		TotalCount:      movedCount + failedCount + partialFailures,
		Errors:          errors,
	}

	return response, nil
}

func (s *Storage) updateDBForTrashedImage(fileID int, newPath string) error {
	log.Printf("Updating image %d to be %s", fileID, ActionTrashed)
	_, err := s.db.Exec(`
				UPDATE images 
				SET action = ?, 
				path = ? 
				WHERE id = ?`,
		ActionTrashed,
		newPath,
		fileID,
	)
	if err != nil {
		log.Printf("Error while Updating image %d to be %s", fileID, ActionTrashed)
		return err
	}
	log.Printf("Updated image %d to be %s", fileID, ActionTrashed)

	return nil
}

func moveImageToTrash(image ImageToTrash, timestamp string) (string, error) {
	trashPath := os.Getenv("TRASH_DIR")
	if trashPath == "" {
		trashPath = "./trash"
	}
	destPath := filepath.Join(trashPath, timestamp, image.Path)
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return "", err
	}

	if err := os.Rename(image.Path, destPath); err != nil {
		return "", err
	}

	return destPath, nil
}

func (s *Storage) DeletePendingImages() error {
	_, err := s.db.Exec("DELETE FROM images WHERE action = 'pending'")
	if err != nil {
		return fmt.Errorf("failed to delete pending images %w", err)
	}

	_, err = s.db.Exec("DELETE FROM image_groups WHERE id NOT IN (SELECT DISTINCT group_id FROM images)")
	if err != nil {
		return fmt.Errorf("failed to delete image groups with no images in the images table %w", err)
	}
	return nil
}
