package jobs

import (
	"log"

	"time"

	"releasetrackr/helpers"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ClearNonVerifiedUsers is a scheduled task to remove users from the DB who are not verified.
func ClearNonVerifiedUsers() {
	log.Println("[Job][ClearNonVerifiedUsers] Job started")
	sess, _ := helpers.GetDbSession()

	c := sess.DB("releasetrackr").C("users")

	var removeResult *mgo.ChangeInfo

	fromDate := time.Now().Add(-2 * time.Hour)
	toDate := time.Now().Add(-1 * time.Hour)

	log.Printf("[Job][ClearNonVerifiedUsers] Clearing users from %v to %v", fromDate, toDate)

	removeResult, _ = c.RemoveAll(
		bson.M{
			"verified": false,
			"created": bson.M{
				"$gt": fromDate,
				"$lt": toDate,
			},
		},
	)
	log.Printf("[Job][ClearNonVerifiedUsers] %d users cleared", removeResult.Removed)
}
