package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()

	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connectDatabse() {

	client, ctx, cancel, err := Connect()
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancel)
	Ping(client, ctx)
	fmt.Println("Database connected.")
}

func Connect() (*mongo.Client, context.Context,
	context.CancelFunc, error) {
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func Ping(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func Query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	result, err = collection.Find(ctx, query,
		options.Find().SetProjection(field))
	return
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) bool {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertOne accept two argument of type Context
	// and of empty interface
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return false
	}
	return true
}

func DeleteOne(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) bool {
	collection := client.Database(dataBase).Collection(col)
	_, err := collection.DeleteOne(ctx, query)
	if err != nil {
		return false
	}
	return true
}

func UpdateOne(client *mongo.Client, ctx context.Context, dataBase,
	col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	collection := client.Database(dataBase).Collection(col)
	result, err = collection.ReplaceOne(ctx, filter, update)
	return
}
