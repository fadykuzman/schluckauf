package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/fadykuzman/schluckauf/internal/loader"
)

type ScanRequest struct {
	Directory string `json:"directory"`
}

type ScanResponse struct {
	Success    bool   `json:"success"`
	GroupCount int    `json:"groupCount"`
	Message    string `json:"message"`
}

func (h *Handler) ScanDirectory(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate the directory path
	info, err := os.Stat(req.Directory)
	if err != nil {
		http.Error(w, "Directory does not exist", http.StatusBadRequest)
		return
	}

	if !info.IsDir() {
		http.Error(w, "Path is not a directory", http.StatusBadRequest)
		return
	}

	// check if czkawka_cli is installed
	if _, err := exec.LookPath("czkawka_cli"); err != nil {
		http.Error(w, "Czkawka CLI is not installed. Install with: cargo install czkawka_cli", http.StatusInternalServerError)
		return
	}

	// create temp file for JSON output
	tempFile, err := os.CreateTemp("", "czkawka-scan-*.json")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// execute the cli command
	cmd := exec.Command("czkawka_cli", "dup", "-d", req.Directory, "--export-json", tempFile.Name())

	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Scan failed: %s", string(output)), http.StatusInternalServerError)
		return
	}

	// parse the JSON results
	groups, err := loader.ParseJSON(tempFile.Name())
	if err != nil {
		http.Error(w, "Failed to parse scan results", http.StatusInternalServerError)
		return
	}

	// Clear pending data
	if err := h.store.DeletePendingData(); err != nil {
		http.Error(w, "Failed to clear pending data", http.StatusInternalServerError)
		return
	}

	// Load new scan results into database
	for _, group := range groups {
		gid, err := h.store.CreateGroup(group.Hash, group.Size, group.FileCount)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create group with hash: %s", group.Hash), http.StatusInternalServerError)
			return
		}

		for _, file := range group.Files {
			_, err := h.store.CreateFile(gid, file.Path, group.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create file with gid %s and path %s", gid, file.Path), http.StatusInternalServerError)
				return
			}
		}
	}

	// Return success response
	resp := ScanResponse{
		Success:    true,
		GroupCount: len(groups),
		Message:    "Scan successfully done",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
