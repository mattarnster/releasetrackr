package jobs

import (
	"log"

	"gopkg.in/mgo.v2/bson"

	"releasetrackr/helpers"
	"releasetrackr/models"
)

// SendNewReleaseNotification sends a new release notification
func SendNewReleaseNotification(repo models.Repo, newRelease models.Release) {
	sess, err := helpers.GetDbSession()

	if err != nil {
		log.Panicf("Couldn't get DB session: %v", err.Error())
		return
	}

	log.Printf("[Job][SendNewReleaseNotification] Starting new release notifications job")

	c := sess.DB("releasetrackr").C("tracks")

	var tracks []models.Track

	c.Find(bson.M{"repoID": repo.ID}).All(&tracks)

	if len(tracks) == 0 {
		return
	}

	for _, track := range tracks {
		c = sess.DB("releasetrackr").C("users")
		var user models.User
		err := c.FindId(track.UserID).One(&user)
		if err != nil {
			log.Panicf("No user with this ID assigned with this track record. %v %v", track.ID, repo.ID)
		}

		log.Printf("Sending user notfication %s", user.Email)
		helpers.SendNotificationEmail(repo, user.Email, newRelease)
	}

}
