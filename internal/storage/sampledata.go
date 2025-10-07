package storage

import (
	_ "database/sql"
	"log"
)

func (s *Storage) InsertSampleData() error {
	gid, err1 := s.CreateGroup("sample-hash-1", 2156432, 3)
	s.CreateFile(gid, "/photos/IMG_1234.jpg", 2156432)
	s.CreateFile(gid, "/photos/backup/IMG_1234.jpg", 2156432)
	s.CreateFile(gid, "/photos/copy/IMG_1234.jpg", 1847000)

	if err1 != nil {
		log.Fatal("Could not insert Group 1 in sample data", err1)
	}

	gid2, err2 := s.CreateGroup("sample-hash-2", 3500000, 2)
	s.CreateFile(gid2, "/photos/IMG_5678.jpg", 3500000)
	s.CreateFile(gid2, "/photos/old/IMG_5678.jpg", 3500000)

	if err2 != nil {
		log.Fatal("Could not insert Group 2 in sample data", err2)
	}
	return nil
}
