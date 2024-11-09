package db

import (
	"context"
	"fmt"
	"training/file-index/pb"

	"go.mongodb.org/mongo-driver/mongo"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	TransferTx(ctx context.Context, arg *pb.FileAttr) (*mongo.InsertOneResult, error)
	UpdateTx(ctx context.Context, arg *pb.FileAttr) (*mongo.UpdateResult, error)
	DeleteTx(ctx context.Context, id string) (*mongo.DeleteResult, error)
}

// MongoStore provides all functions to execute Mongo queries and transactions
type MongoStore struct {
	mgClient   *mongo.Client
	collection *mongo.Collection
}

// NewStore creates a new store
func NewStore(mgoClient *mongo.Client, collection *mongo.Collection) Store {
	return &MongoStore{
		mgClient:   mgoClient,
		collection: collection,
	}
}

// Close the connection
func (store *MongoStore) CloseConnection(client *mongo.Client, context context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		if err := client.Disconnect(context); err != nil {
			panic(err)
		}
		fmt.Println("Close connection is called")
	}()
}
