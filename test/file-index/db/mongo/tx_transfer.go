package db

import (
	"context"
	"training/file-index/pb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (store *MongoStore) TransferTx(ctx context.Context, arg *pb.FileAttr) (*mongo.InsertOneResult, error) {

	collection := store.collection
	result, err := collection.InsertOne(ctx, arg, options.InsertOne())
	if err != nil {
		return nil, err
	}
	return result, err
}

func (store *MongoStore) UpdateTx(ctx context.Context, arg *pb.FileAttr) (*mongo.UpdateResult, error) {
	collection := store.collection
	result, err := collection.UpdateOne(ctx, arg, options.Update())
	if err != nil {
		return nil, err
	}
	return result, err
}

func (store *MongoStore) DeleteTx(ctx context.Context, name string) (*mongo.DeleteResult, error) {
	collection := store.collection

	filter := bson.M{"name": name}

	result, err := collection.DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		return nil, err
	}
	return result, err
}
