package jobs

import (
	"log"
	"net/http"

	"github.com/mattarnster/releasetrackr/helpers"
	"github.com/mattarnster/releasetrackr/models"
)

var repos []models.Repo

// GetNewReleases gets new releases from the Github API
func GetNewReleases() {
	// Grab a bunch of repos
	sess, _ := helpers.GetDbSession()

	c := sess.DB("releasetrackr").C("repos")
	c.Find(nil).All(&repos)

	log.Printf("[Job][GetNewReleases] Result count: %v", len(repos))

	if len(repos) == 0 {
		log.Println("[Job][GetNewReleases] No repos in DB")
		return
	}

	// Then start firing off requests to the API
	for _, repo := range repos {
		log.Printf("[Job][GetNewReleases] Looking for release for %+v", repo.Repo)

		resp, err := http.Get("https://api.github.com/repos/" + repo.Repo + "/releases")
		if err != nil {
			log.Printf("[Job][GetNewReleases] API Request failed: %v", err.Error())
		}

		defer resp.Body.Close()

		// Fill the record with the data from the JSON
		var record models.GithubReleases

		//var f interface{}

		//err := json.Unmarshal(resp.Body, &f)

		log.Printf("[Job][GetNewReleases] New record: %v", record)
	}
}
