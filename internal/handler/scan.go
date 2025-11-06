package handler

import (
	"encoding/json"
	"fmt"
	"log"
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

	scansDir := os.Getenv("SCANS_DIR")
	if scansDir == "" {
		scansDir = "./scans"
	}

	os.MkdirAll(scansDir, 0o770)

	tempFile, err := os.CreateTemp(scansDir, "czkawka-scan-*.json")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}

	// defer os.Remove(tempFile.Name())
	tempFile.Close()

	// execute the cli command
	cmd := exec.Command("czkawka_cli", "image", "-d", req.Directory, "-C", tempFile.Name())

	output, err := cmd.CombinedOutput()
	groups, parseErr := loader.ParseImageDuplicates(tempFile.Name())
	if parseErr != nil {
		http.Error(w, fmt.Sprintf("Scan failed: %s (parse error: %v)", string(output), parseErr), http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Printf("warning: czkawka exited with error (%v) but produced valid output", err)
	}

	// parse the JSON results
	if len(groups) == 0 {
		resp := ScanResponse{
			Success:    true,
			GroupCount: 0,
			Message:    "No Duplicates found",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Clear pending data
	if err := h.store.DeletePendingImages(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("error deleting pending images: %v", err),
			http.StatusInternalServerError)
		return
	}

	// Load new scan results into database
	for _, group := range groups {
		gid, err := h.store.CreateImageGroup(
			group.Hash,
			group.Size,
			group.ImageCount)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create group with hash: %d", group.Hash), http.StatusInternalServerError)
			return
		}

		for _, file := range group.Images {
			_, err := h.store.CreateImage(gid, file.Path, group.Size)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create image with gid %d and path %s: \n error (%v)", gid, file.Path, err), http.StatusInternalServerError)
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
