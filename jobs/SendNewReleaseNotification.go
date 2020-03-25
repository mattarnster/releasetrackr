package jobs

import (
	"context"
	"log"

	"releasetrackr/helpers"
	"releasetrackr/models"

	"go.mongodb.org/mongo-driver/bson"
)

// SendNewReleaseNotification sends a new release notification
func SendNewReleaseNotification(repo models.Repo, newRelease models.Release) {
	sess, err := helpers.GetDbSession()

	if err != nil {
		log.Panicf("Couldn't get DB session: %v", err.Error())
		return
	}

	log.Printf("[Job][SendNewReleaseNotification] Starting new release notifications job")

	c := sess.Database("releasetrackr").Collection("tracks")

	count, err := c.CountDocuments(context.Background(), bson.D{})

	cur, err := c.Find(context.Background(), bson.M{"repoID": repo.ID})
	//defer cur.Close(context.Background())

	if count == 0 {
		return
	}

	for cur.Next(context.TODO()) {
		log.Println("[+] Got here")
		var track models.Track
		_ = cur.Decode(&track)

		c = sess.Database("releasetrackr").Collection("users")
		var user models.User

		res := c.FindOne(context.Background(), bson.M{"_id": track.UserID}).Decode(&user)
		if res != nil {
			log.Panicf("No user with this ID assigned with this track record. %v %v", track.ID, repo.ID)
		}

		log.Printf("Sending user notfication %s", user.Email)
		helpers.SendNotificationEmail(repo, user.Email, newRelease)
	}
	cur.Close(context.TODO())
}
