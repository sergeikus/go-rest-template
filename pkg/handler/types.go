package handler

import (
	"github.com/sergeikus/go-rest-template/pkg/storage"
)

// API is used to pass database interface to handlers
type API struct {
	DB storage.DB
}

// GetRequest represents for a key
type GetRequest struct {
	Key string `json:"string"`
}
