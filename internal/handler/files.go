package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

	info, err := os.Stat(absPath)
	if err != nil || info.IsDir() {
		http.Error(w, "File not Found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, absPath)
}

type UpdateFileActionRequest struct {
	Action storage.FileAction `json:"action"`
}

func (h *Handler) UpdateFileAction(w http.ResponseWriter, r *http.Request) {
	gidStr := r.PathValue("gid")
	groupID, gerr := strconv.Atoi(gidStr)
	if gerr != nil {
		http.Error(w, "Invalid groupID", http.StatusBadRequest)
	}

	fidStr := r.PathValue("fid")
	fileID, err := strconv.Atoi(fidStr)
	if err != nil {
		http.Error(w, "Invalid File ID", http.StatusBadRequest)
	}

	var req UpdateFileActionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Action != storage.ActionKeep && req.Action != storage.ActionTrash {
		http.Error(w, "Action must be 'keep' or 'trash'", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateFileAction(groupID, fileID, req.Action); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) TrashFiles(w http.ResponseWriter, r *http.Request) {
	filesToTrash, err := h.store.TrashFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filesToTrash)
}
