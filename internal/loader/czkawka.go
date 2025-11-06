// Package loader provides functions to parse and load duplicate image scan outputs
package loader

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

type CzkawkaFileOutput map[string][][]FileInfo

type FileInfo struct {
	Path         string `json:"path"`
	ModifiedDate int64  `json:"modified_date"`
	Size         int64  `json:"size"`
	Hash         string `json:"hash"`
}

type DuplicateFileGroup struct {
	Hash      string
	Size      int64
	FileCount int
	Files     []FileInfo
}

type CzkawkaImageOutput [][]ImageInfo

type ImageInfo struct {
	Path         string  `json:"path"`
	ModifiedDate int64   `json:"modified_date"`
	Size         int64   `json:"size"`
	Hash         []int   `json:"hash"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	Similarity   float64 `json:"similarity"`
}

type DuplicateImageGroup struct {
	Hash       []int       `json:"hash"`
	Size       int64       `json:"size"`
	ImageCount int         `json:"image_count"`
	Images     []ImageInfo `json:"images"`
}

func ParseFileDuplicates(filepath string) ([]DuplicateFileGroup, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var czkawka CzkawkaFileOutput

	if err := json.Unmarshal(data, &czkawka); err != nil {
		return nil, err
	}

	var groups []DuplicateFileGroup

	for sizeKey, dupGroups := range czkawka {
		for _, fileGroup := range dupGroups {
			if len(fileGroup) < 2 {
				continue
			}

			firstFile := fileGroup[0]

			groupHash := firstFile.Hash

			if groupHash == "" {
				h := sha256.New()
				h.Write([]byte(sizeKey))
				for _, f := range fileGroup {
					h.Write([]byte(f.Path))
				}
				groupHash = fmt.Sprintf("%x", h.Sum(nil))
			}

			groups = append(groups, DuplicateFileGroup{
				Hash:      groupHash,
				Size:      firstFile.Size,
				FileCount: len(fileGroup),
				Files:     fileGroup,
			})

		}
	}
	return groups, nil
}

func ParseImageDuplicates(filepath string) ([]DuplicateImageGroup, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var czkawka CzkawkaImageOutput

	if err := json.Unmarshal(data, &czkawka); err != nil {
		return nil, err
	}

	var groups []DuplicateImageGroup

	for _, dupGroup := range czkawka {
		if len(dupGroup) < 2 {
			continue
		}

		firstFile := dupGroup[0]

		groupHash := firstFile.Hash

		if len(groupHash) == 0 {
			return nil, fmt.Errorf(
				"image hash is missing for file: %s ", firstFile.Path)
		}

		groups = append(groups, DuplicateImageGroup{
			Hash:       groupHash,
			Size:       firstFile.Size,
			ImageCount: len(dupGroup),
			Images:     dupGroup,
		})

	}
	return groups, nil
}
