package handler

import (
	"errors"

	"github.com/sergeikus/go-rest-template/pkg/auth"
	"github.com/sergeikus/go-rest-template/pkg/storage"
)

// API is used to pass database interface to handlers
type API struct {
	DB   storage.DB
	Auth auth.Auth
}

const (
	// MsgStatusOK represents a success string message
	MsgStatusOK = "{\"status\": \"OK\"}"
)

// DataAdditionRequest represents for a key
type DataAdditionRequest struct {
	Data string `json:"data"`
}

var errDataAdditionRequestNoData = errors.New("data to be added must be non-empty string")

// Validate performs request validation
func (dar *DataAdditionRequest) Validate() error {
	if len(dar.Data) == 0 {
		return errDataAdditionRequestNoData
	}
	return nil
}
