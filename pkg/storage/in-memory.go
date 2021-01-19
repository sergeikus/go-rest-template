package storage

import (
	"fmt"
	"sync"

	"github.com/sergeikus/go-rest-template/pkg/types"
)

// InMemoryStorage represents an in-memory database
type InMemoryStorage struct {
	data map[int]types.Data
	// Locking mutex for addition
	// To mitigate race condition
	mutex sync.Mutex
	// Simulates primary key index
	index int
}

// Connect simulates connection to database
func (ims *InMemoryStorage) Connect() error {
	ims.data = make(map[int]types.Data)
	ims.mutex = sync.Mutex{}
	ims.index = 1
	return nil
}

// Store stores data
func (ims *InMemoryStorage) Store(data string) (id int, err error) {
	if len(data) == 0 {
		return id, fmt.Errorf("data must be non-empty string")
	}
	ims.mutex.Lock()
	defer ims.mutex.Unlock()
	if _, exist := ims.data[ims.index]; exist {
		return id, fmt.Errorf("key with '%d' ID already exists", ims.index)
	}
	ims.data[ims.index] = types.Data{ID: ims.index, String: data}
	id = ims.index

	ims.index++
	return id, nil
}

// GetAll returns all data from the dable
func (ims *InMemoryStorage) GetAll() ([]types.Data, error) {
	return nil, fmt.Errorf("Not implemented")
}

// GetKey returns data for a paricular key
func (ims *InMemoryStorage) GetKey(key int) (d types.Data, err error) {
	d, exist := ims.data[key]
	if !exist {
		return d, fmt.Errorf("data with '%d' key does not exist", key)
	}
	return d, nil
}

// Close does nothing but it's main implementation is to close database connection
func (ims *InMemoryStorage) Close() {
	return
}

// VerifyUserCredentials checks user login in database
func (ims *InMemoryStorage) VerifyUserCredentials(username, passwordHash string) (types.User, error) {
	return types.User{}, fmt.Errorf("not implemented")
}

// GetUserSalt salt returns user password salt
func (ims *InMemoryStorage) GetUserSalt(username string) (salt string, err error) {
	return salt, fmt.Errorf("not implemented")
}

// RegisterUser user registers new user
func (ims *InMemoryStorage) RegisterUser(user types.User) (id int, err error) {
	return id, fmt.Errorf("not implemented")
}
