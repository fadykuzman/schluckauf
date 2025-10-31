package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fadykuzman/schluckauf/internal/handler"
	"github.com/fadykuzman/schluckauf/internal/loader"
	"github.com/fadykuzman/schluckauf/internal/storage"
)

func main() {
	store, err := storage.New("./data/duplicates.db")
	if err != nil {
		log.Fatal(fmt.Errorf("error: %+v", err))
	}
	defer store.Close()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--sample":
			store.InsertSampleData()
			fmt.Println("Sample data inserted")

		case "--load":
			if len(os.Args) < 3 {
				log.Fatal("Usage: --load <path-to-json>")
			}
			jsonPath := os.Args[2]

			groups, err := loader.ParseFileDuplicates(jsonPath)
			if err != nil {
				log.Fatal("Failed to parse JSON: ", err)
			}

			for _, group := range groups {
				gid, err := store.CreateGroup(group.Hash, group.Size, group.FileCount)
				if err != nil {
					log.Printf("Warning: failed to create group: %v", err)
					continue
				}

				for _, file := range group.Files {
					if _, err := store.CreateFile(gid, file.Path, group.Size); err != nil {
						log.Printf("Warning: failed to create file: %v", err)
					}
				}
			}
			fmt.Printf("Loaded %d duplicate groups\n", len(groups))
		}
	}

	h := handler.New(store)

	http.HandleFunc("GET /api/groups", h.ListImageGroups)
	http.HandleFunc("/health", h.Health)
	http.HandleFunc("GET /api/groups/{id}", h.GetGroupImages)
	http.HandleFunc("/api/image", h.ServeImage)
	http.HandleFunc("POST /api/groups/{gid}/files/{fid}", h.UpdateFileAction)
	http.HandleFunc("GET /api/groups/stats", h.GetGroupStats)
	http.HandleFunc("POST /api/files/actions/trash", h.TrashFiles)
	http.HandleFunc("POST /api/scan", h.ScanDirectory)

	http.Handle("/", http.FileServer(http.Dir("./web")))
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
