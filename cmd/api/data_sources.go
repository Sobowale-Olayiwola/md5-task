package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DataSources struct {
	DB         *mongo.Database
	Client     *mongo.Client
	InMemoryDB *redis.Client
}

// InitDS establishes connections to fields in dataSources
func initDS() (*DataSources, error) {
	log.Printf("Initializing data sources\n")
	// Initialize MongoDB connection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	//Verify MongoDB connection
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	mongoDB := client.Database(os.Getenv("MONGO_DATABASE"))

	//Initalize RedisDB connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_CONNECTION_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx = context.TODO()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis", err.Error())
	} else {
		log.Println("Connected to Redis")
	}

	return &DataSources{
		DB:         mongoDB,
		Client:     client,
		InMemoryDB: rdb,
	}, nil
}

// close to be used in graceful server shutdown
func (d *DataSources) close() error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	if err := d.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("error closing MongoDB: %w", err)
	}

	return nil
}
