package jobs

import (
	"log"

	"github.com/mattarnster/releasetrackr/helpers"
	"github.com/mattarnster/releasetrackr/models"
)

var repos []models.Repo

// GetNewReleases gets new releases from the Github API
func GetNewReleases() {
	// Grab a bunch of repos
	sess, _ := helpers.GetDbSession()

	c := sess.DB("releasetracker").C("repos")
	c.Find(nil).All(&repos)

	if len(repos) == 0 {
		log.Println("[Job][GetNewReleases] No repos in DB")
		return
	}

	// Then start firing off requests to the API
	for _, repo := range repos {
		log.Printf("Looking for release for %s", repo.ID.String())
	}
}
