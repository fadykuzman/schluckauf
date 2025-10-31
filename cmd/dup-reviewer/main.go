package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fadykuzman/schluckauf/internal/handler"
	"github.com/fadykuzman/schluckauf/internal/storage"
)

func main() {
	store, err := storage.New("./data/duplicates.db")
	if err != nil {
		log.Fatal(fmt.Errorf("error: %+v", err))
	}
	defer store.Close()

	h := handler.New(store)

	http.HandleFunc("GET /api/groups", h.ListImageGroups)
	http.HandleFunc("/health", h.Health)
	http.HandleFunc("GET /api/groups/{id}", h.GetGroupImages)
	http.HandleFunc("/api/image", h.ServeImage)
	http.HandleFunc("POST /api/groups/{gid}/files/{fid}", h.UpdateImageAction)
	http.HandleFunc("GET /api/groups/stats", h.GetGroupStats)
	http.HandleFunc("POST /api/files/actions/trash", h.TrashImages)
	http.HandleFunc("POST /api/scan", h.ScanDirectory)

	http.Handle("/", http.FileServer(http.Dir("./web")))
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
