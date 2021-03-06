package jobs

import (
	"context"
	"log"
	"releasetrackr/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"time"
)

// ClearNonVerifiedUsers is a scheduled task to remove users from the DB who are not verified.
func ClearNonVerifiedUsers() {
	log.Println("[Job][ClearNonVerifiedUsers] Job started")
	sess, _ := db.GetDbSession()

	c := sess.Database("releasetrackr").Collection("users")

	fromDate := time.Now().Add(-2 * time.Hour)
	toDate := time.Now().Add(-1 * time.Hour)

	log.Printf("[Job][ClearNonVerifiedUsers] Clearing users from %v to %v", fromDate, toDate)

	var removeResult *mongo.DeleteResult
	removeResult, _ = c.DeleteMany(
		context.Background(),
		bson.M{
			"verified": false,
			"created": bson.M{
				"$gt": fromDate,
				"$lt": toDate,
			},
		},
	)
	log.Printf("[Job][ClearNonVerifiedUsers] %d users cleared", removeResult.DeletedCount)
}
