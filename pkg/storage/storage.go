package storage

import "github.com/sergeikus/go-rest-template/pkg/types"

// DB represents a storage interface which can be
// in memory or an external database
type DB interface {
	Connect() error
	Close()

	Store(data string) (int, error)
	GetAll() ([]types.Data, error)
	GetKey(key int) (types.Data, error)

	// Authentication actions
	VerifyUserCredentials(username, passwordHash string) (types.User, error)
	GetUserSalt(username string) (string, error)

	// User management
	RegisterUser(types.User) (int, error)
}

const (
	// DatabaseTypeInMemory defines an in-memory database
	DatabaseTypeInMemory = "in-memory"
	// DatabaseTypePostgre defines a postgre database
	DatabaseTypePostgre = "postgres"
)
