package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	db "training/file-index/db/mongo"
	"training/file-index/pb"

	"go.mongodb.org/mongo-driver/mongo"
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
func NewInMemoryFileStore(mgClient *mongo.Client, collection *mongo.Collection) *InMemoryFileStore {
	return &InMemoryFileStore{
		store: db.NewStore(mgClient, collection),
	}
}

// Save saves the file to the store
func (store *InMemoryFileStore) Save(file *pb.FileAttr) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	fmt.Println("Insert file", file)
	store.store.TransferTx(context.TODO(), file)
	return nil
}

func (store *InMemoryFileStore) Update(file *pb.FileAttr) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	fmt.Println("Update file", file)
	store.store.UpdateTx(context.TODO(), file)
	return nil
}

func (store *InMemoryFileStore) Delete(id string) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	fmt.Println("Delete file", id)
	store.store.DeleteTx(context.TODO(), id)
	return nil
}
