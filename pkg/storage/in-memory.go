package storage

// InMemoryStorage represents an in-memory database
type InMemoryStorage struct {
	data map[string]string
}

// Connect simulates connection to database
func (ims *InMemoryStorage) Connect() error {
	ims.data = make(map[string]string)
	return nil
}
