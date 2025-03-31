package db

import (
	"context"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// DbConnection Will hold the global DB Connection
var DbConnection *mongo.Client
var lock = &sync.Mutex{}

// GetDbSession returns the currently active DB connection (if not, then it creates one)
func GetDbSession() (*mongo.Client, error) {
	if DbConnection == nil {
		lock.Lock()
		defer lock.Unlock()
		opts := options.Client().ApplyURI(connectionString())
		// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(opts)
		if err != nil {
			return nil, err
		}
		DbConnection = client
	}
	return DbConnection, nil
}

// KillDbSession kills the currently active DB session
func KillDbSession() {
	if DbConnection != nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		DbConnection.Disconnect(ctx)
	}
}

// Build the connection string
func connectionString() string {
	buildHost := os.Getenv("MONGO_HOST")

	return buildHost
}
