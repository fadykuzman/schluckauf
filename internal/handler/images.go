package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fadykuzman/schluckauf/internal/storage"
)

func (h *Handler) ServeImage(w http.ResponseWriter, r *http.Request) {
	requestedPath := r.URL.Query().Get("path")
	if requestedPath == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	cleanPath := filepath.Clean(requestedPath)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	baseDir := os.Getenv("PHOTOS_DIR")
	if baseDir == "" {
		baseDir = "./test"
	}
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Server configuration error (%v)", err), http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(absPath, absBase) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil || info.IsDir() {
		http.Error(w, "File not Found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, absPath)
}

type UpdateImageActionRequest struct {
	Action storage.ImageAction `json:"action"`
}

func (h *Handler) UpdateImageAction(w http.ResponseWriter, r *http.Request) {
	gidStr := r.PathValue("gid")
	groupID, gerr := strconv.Atoi(gidStr)
	if gerr != nil {
		http.Error(w, "Invalid groupID", http.StatusBadRequest)
		return
	}

	fidStr := r.PathValue("fid")
	fileID, err := strconv.Atoi(fidStr)
	if err != nil {
		http.Error(w, "Invalid File ID", http.StatusBadRequest)
		return
	}

	var req UpdateImageActionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Action != storage.ActionKeep && req.Action != storage.ActionTrash {
		http.Error(w, "Action must be 'keep' or 'trash'", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateImageAction(groupID, fileID, req.Action); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) TrashImages(w http.ResponseWriter, r *http.Request) {
	response, err := h.store.TrashImages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
