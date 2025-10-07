package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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

func (h *Handler) GetGroupStats(w http.ResponseWriter, r *http.Request) {
	gs, err := h.store.GetGroupStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gs)
}
