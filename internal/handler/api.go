package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fadykuzman/schluckauf/internal/storage"
)

type Handler struct {
	store *storage.Storage
}

func New(store *storage.Storage) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.store.ListGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handler) GetGroupFiles(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Group ID", http.StatusBadRequest)
		return
	}
	files, err := h.store.GetGroupFiles(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

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

type UpdateActionRequest struct {
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

	var req UpdateActionRequest

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
