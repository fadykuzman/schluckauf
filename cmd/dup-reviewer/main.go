package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fadykuzman/schluckauf/internal/loader"
	"github.com/fadykuzman/schluckauf/internal/storage"
)

func main() {
	store, err := storage.New("./data/duplicates.db")
	if err != nil {
		log.Fatal(err)
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

			groups, err := loader.ParseJSON(jsonPath)
			if err != nil {
				log.Fatal("Failed to parse JSON: ", err)
			}

			for _, group := range groups {
				gid, err := store.CreateGroup(group.Hash, group.Size, group.FileCount)
				if err != nil {
					log.Printf("Warning: failed to create group: %v", err)
					continue
				}

				for _, filePath := range group.Files {
					if err := store.CreateFile(gid, filePath, group.Size); err != nil {
						log.Printf("Warning: failed to create file: %v", err)
					}
				}
			}
			fmt.Printf("Loaded %d duplicate groups\n", len(groups))
		}
	}

	http.HandleFunc("GET /api/groups", func(w http.ResponseWriter, r *http.Request) {
		groups, err := store.ListGroups()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("GET /api/groups/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		groupID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid Group ID", http.StatusBadRequest)
			return
		}
		files, err := store.GetGroupFiles(groupID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	})

	http.Handle("/", http.FileServer(http.Dir("./web")))
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
