package handler

import (
	"fmt"

	"github.com/sergeikus/go-rest-template/pkg/storage"
)

// API is used to pass database interface to handlers
type API struct {
	DB storage.DB
}

// KeyAdditionRequest represents for a key
type KeyAdditionRequest struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

// Validate performs request validation
func (kar *KeyAdditionRequest) Validate() error {
	if len(kar.Key) == 0 {
		return fmt.Errorf("key must be non-empty string")
	}
	return nil
}
