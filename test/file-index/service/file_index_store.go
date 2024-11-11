package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	db "training/file-index/db/sqlc"
	"training/file-index/pb"
)

var ErrAlreadyExists = errors.New("record already exists")

type FileStore interface {
	Save(files *pb.FileAttr) error
	Update(files *pb.FileAttr) error
	Delete(id string) error
}

// InMemoryLaptopStore stores laptop in memory
type InMemoryFileStore struct {
	mutex sync.RWMutex
	store db.Store
}

// NewInMemoryFileStore returns a new InMemoryFileStore
func NewInMemoryFileStore(db db.Store) *InMemoryFileStore {
	return &InMemoryFileStore{
		store: db,
	}
}

// Save saves a file attribute in the store
func (store *InMemoryFileStore) Save(files *pb.FileAttr) error {
	store.mutex.Lock()
	ctx := context.Background()
	arg := db.InsertFileParams{
		Name:       files.Name,
		Extension:  files.Type,
		Size:       files.Size,
		Attributes: files.Path,
		Content:    files.Content,
		CreatedAt:  files.CreatedAt.AsTime(),
		ModifiedAt: files.CreatedAt.AsTime(),
		AccessedAt: files.ModifiedAt.AsTime(),
	}
	fmt.Println("Insert file with name: ", files.Name)
	store.store.InsertFile(ctx, arg)
	defer store.mutex.Unlock()
	return nil
}

// Update updates a file attribute in the store
func (store *InMemoryFileStore) Update(files *pb.FileAttr) error {
	// Implementation here
	return nil
}

// Delete deletes a file attribute by its ID
func (store *InMemoryFileStore) Delete(id string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// Your logic to delete the file attribute by ID
	return nil
}
