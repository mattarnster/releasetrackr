package jobs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

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

		var f interface{}

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &f)
		if err != nil {
			log.Fatalf("[Job][GetNewReleases] Error unmarshaling JSON - likely invalid.")
			return
		}

		objects := f.([]interface{})
		first := objects[0].(map[string]interface{})

		c = sess.DB("releasetrackr").C("releases")
		var existingRelease models.Release

		err = c.Find(bson.M{"release_id": first["id"].(float64)}).One(&existingRelease)

		if err != nil {
			// Release doesn't exist, add it to the DB
			newReleaseID := bson.NewObjectId()
			createdAtTime, caTErr := time.Parse(time.RFC3339Nano, first["created_at"].(string))
			if caTErr != nil {
				log.Fatalf("Created at time parse failed %v", caTErr.Error())
			}
			publishedAtTime, paTErr := time.Parse(time.RFC3339Nano, first["published_at"].(string))
			if paTErr != nil {
				log.Fatalf("Published at time parse failed: %v", paTErr.Error())
			}
			newRelease := models.Release{
				ID:                 newReleaseID,
				ReleaseID:          first["id"].(float64),
				URL:                first["html_url"].(string),
				Tag:                first["tag_name"].(string),
				Name:               first["name"].(string),
				ReleaseCreatedAt:   createdAtTime,
				ReleasePublishedAt: publishedAtTime,
				Body:               first["body"].(string),
			}

			log.Printf("[Job][GetNewReleases] New release record: %v", newRelease)

			c.Insert(&newRelease)
		}
	}
}
