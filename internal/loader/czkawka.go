package loader

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

type CzkawkaOutput map[string][][]FileInfo

type FileInfo struct {
	Path         string `json:"path"`
	ModifiedDate int64  `json:"modified_date"`
	Size         int64  `json:"size"`
	Hash         string `json:"hash"`
}

type DuplicateGroup struct {
	Hash      string
	Size      int64
	FileCount int
	Files     []string
}

func ParseJSON(filepath string) ([]DuplicateGroup, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var czkawka CzkawkaOutput

	if err := json.Unmarshal(data, &czkawka); err != nil {
		return nil, err
	}

	var groups []DuplicateGroup

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

			var paths []string
			for _, file := range fileGroup {
				paths = append(paths, file.Path)
			}

			groups = append(groups, DuplicateGroup{
				Hash:      groupHash,
				Size:      firstFile.Size,
				FileCount: len(fileGroup),
				Files:     paths,
			})

		}
	}
	return groups, nil
}
