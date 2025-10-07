// Package handler provides REST API operations for duplicate photo groups and files
package handler

import (
	"github.com/fadykuzman/schluckauf/internal/storage"
)

type Handler struct {
	store *storage.Storage
}

func New(store *storage.Storage) *Handler {
	return &Handler{store: store}
}
