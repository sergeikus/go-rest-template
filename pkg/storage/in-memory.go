package storage

import (
	"fmt"
	"sync"
)

// InMemoryStorage represents an in-memory database
type InMemoryStorage struct {
	data map[string]string
	// Locking mutex for addition
	// To mitigate race condition
	mutex sync.Mutex
}

// Connect simulates connection to database
func (ims *InMemoryStorage) Connect() error {
	ims.data = make(map[string]string)
	ims.mutex = sync.Mutex{}
	return nil
}

// Store stores data
func (ims *InMemoryStorage) Store(key, data string) error {
	if len(key) == 0 {
		return fmt.Errorf("key must be non-empty string")
	}
	ims.mutex.Lock()
	defer ims.mutex.Unlock()
	if _, exist := ims.data[key]; exist {
		return fmt.Errorf("key with '%s' ID already exists", key)
	}
	ims.data[key] = data
	return nil
}
