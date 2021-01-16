package storage

import "github.com/sergeikus/go-rest-template/pkg/types"

// DB represents a storage interface which can be
// in memory or an external database
type DB interface {
	Connect() error
	Store(data string) (int, error)
	GetAll() ([]types.Data, error)
	GetKey(key int) (types.Data, error)
	Close()
}

const (
	// DatabaseTypeInMemory defines an in-memory database
	DatabaseTypeInMemory = "in-memory"
	// DatabaseTypePostgre defines a postgre database
	DatabaseTypePostgre = "postgres"
)
