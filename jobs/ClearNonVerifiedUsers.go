package jobs

import (
	"context"
	"log"
	"releasetrackr/db"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"time"
)

// ClearNonVerifiedUsers is a scheduled task to remove users from the DB who are not verified.
func ClearNonVerifiedUsers() {
	log.Println("[Job][ClearNonVerifiedUsers] Job started")
	sess, err := db.GetDbSession()
	if err != nil {
		panic(err)
	}

	c := sess.Database("releasetrackr").Collection("users")

	fromDate := time.Now().Add(-2 * time.Hour)
	toDate := time.Now().Add(-1 * time.Hour)

	log.Printf("[Job][ClearNonVerifiedUsers] Clearing users from %v to %v", fromDate, toDate)

	var removeResult *mongo.DeleteResult
	removeResult, err = c.DeleteMany(
		context.Background(),
		bson.M{
			"verified": false,
			"created": bson.M{
				"$gt": fromDate,
				"$lt": toDate,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	log.Printf("[Job][ClearNonVerifiedUsers] %d users cleared", removeResult.DeletedCount)
}
