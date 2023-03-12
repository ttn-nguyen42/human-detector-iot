package database

import (
	"context"
	"iot_api/utils"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once

/*
A singleton that stores the database client, used in the data layer
to make queries and write data.
*/
var client *mongo.Client = nil

func GetClient() *mongo.Client {
	if client != nil {
		return client
	}
	// Initializes the singleton
	// Make sure that for all goroutines/threads
	// the singleton is only initialized once
	once.Do(func() {
		client = Connect()
	})
	return client
}

// Returns a context that prevents the database query
// from taking too much time
func GetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	return ctx, cancel
}

// Creates a connection with the database
// and make sure the connection is working
func Connect() *mongo.Client {
	driver, err := utils.GetMongoDriver()
	if err != nil {
		// Fatal exits the application on error
		logrus.Fatal(err.Error())
		return nil
	}
	newClient, err := mongo.NewClient(options.Client().ApplyURI(driver))
	if err != nil {
		logrus.Fatal(err.Error())
		return nil
	}
	ctx, cancel := GetContext()
	defer cancel()
	err = newClient.Connect(ctx)
	if err != nil {
		logrus.Fatal(err.Error())
		return nil
	}
	err = newClient.Ping(ctx, nil)
	if err != nil {
		logrus.Fatal(err.Error())
		return nil
	}
	logrus.Info("Database connected successfully")
	return newClient
}

// Close the connection to the database
func Close() {
	ctx, cancel := GetContext()
	defer cancel()
	GetClient().Disconnect(ctx)
	logrus.Info("Database connection closed")
}

// Decorator design pattern
// It adds behaviour like adding error handling, ID conversion
// and makes the code more testable by minimizing addition of behaviour in each decoration layer
type MongoCollection[E interface{}] struct {
	// Implements Collection[E interface{}]
	Col *mongo.Collection
}

/*
A wrapper for MongoDB's InsertOne.
Uses the same parameters but instead, returns the key as a string for convenience
*/
func (c *MongoCollection[E]) InsertOne(ctx context.Context, document E) (string, error) {
	result, err := c.Col.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	// Ignore the error
	// Result of InsertOne from the official driver always returns a primitype.ObjectID
	id, _ := result.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

/*
A wrapper for MongoDB's FindOne.
Uses the same parameters but also requires a pointer to an object of entity.
The result of the query will be copied over to that object
*/
func (c *MongoCollection[E]) FindOne(ctx context.Context, result *E, filter interface{}) error {
	err := c.Col.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

/*
A wrapper for MongoDB's FindAll.
Uses the same parameters but returns a list of objects of entities instead of cursors
*/
func (c *MongoCollection[E]) FindAll(ctx context.Context, filter interface{}) ([]E, error) {
	cursor, err := c.Col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var results []E
	for cursor.Next(ctx) {
		var result E
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

/*
A wrapper for MongoDB's UpdateOne
*/
func (c *MongoCollection[E]) UpdateOne(ctx context.Context, filter interface{}, upt interface{}) (*UpdateResult, error) {
	result, err := c.Col.UpdateOne(ctx, filter, upt)
	if err != nil {
		return nil, err
	}
	ret := &UpdateResult{
		MatchedFilter: int(result.MatchedCount),
		MatchedField: int(result.ModifiedCount),
	}
	return ret, nil
}