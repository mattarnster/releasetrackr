package helpers

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DbConnection Will hold the global DB Connection
var DbConnection *mongo.Client

// GetDbSession returns the currently active DB connection (if not, then it creates one)
func GetDbSession() (*mongo.Client, error) {
	// if DbConnection == nil {
	// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// 	err := client.Dial(connectionString())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	DbConnection = session
	// }
	// return DbConnection, nil
	if DbConnection == nil {
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT")))
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
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
	buildPort := "27017"
	if os.Getenv("MONGO_PORT") != "" {
		buildPort = os.Getenv("MONGO_PORT")
	}

	return buildHost + ":" + string(buildPort)
}
