package storage

// DB represents a storage interface which can be
// in memory or an external database
type DB interface {
	Connect() error
}

const (
	// DatabaseTypeInMemory defines in memory database
	DatabaseTypeInMemory = "in-memory"
)
