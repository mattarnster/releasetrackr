package jobs

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"releasetrackr/helpers"
	"releasetrackr/models"
)

var repos []models.Repo

var existingRelease models.Release
var isNewRelease = false
var newRelease models.Release

// GetNewReleases gets new releases from the Github API
func GetNewReleases() {
	// Grab a bunch of repos
	sess, _ := helpers.GetDbSession()

	c := sess.Database("releasetrackr").Collection("repos")

	count, err := c.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Printf("[Job][GetNewReleases] Failed to find repos")
		return
	}

	if count == 0 {
		log.Println("[Job][GetNewReleases] No repos in DB")
		return
	}

	cur, _ := c.Find(context.Background(), bson.D{})
	defer cur.Close(context.Background())

	log.Printf("[Job][GetNewReleases] Result count: %v", count)
	for cur.Next(context.Background()) {
		var repo models.Repo
		err := cur.Decode(&repo)

		log.Printf("[Job][GetNewReleases] Looking for release for %+v", repo.Repo)

		resp, err := http.Get("https://api.github.com/repos/" + repo.Repo + "/releases")
		if err != nil {
			log.Printf("[Job][GetNewReleases] API Request failed: %v", err.Error())
		}

		log.Printf("[Job][GetNewReleases] Github ratelimit will be hit in in %v requests.", resp.Header["X-Ratelimit-Remaining"])
		log.Printf("[Job][GetNewReleases] Ratelimit will reset at %v", resp.Header["X-Ratelimit-Reset"])

		defer resp.Body.Close()

		var f interface{}

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &f)
		if err != nil {
			log.Fatalf("[Job][GetNewReleases] Error unmarshaling JSON - likely invalid.")
			return
		}

		objects := f.([]interface{})

		if len(objects) == 0 {
			log.Printf("[Job][GetNewReleases] No releases found for %s", repo.Repo)
			return
		}

		first := objects[0].(map[string]interface{})

		c = sess.Database("releasetrackr").Collection("releases")

		err = c.FindOne(context.Background(), bson.M{"release_id": first["id"].(float64)}).Decode(&existingRelease)

		// Not found - Add it to the DB
		if err != nil {
			isNewRelease = true

			createdAtTime, caTErr := time.Parse(time.RFC3339Nano, first["created_at"].(string))
			if caTErr != nil {
				log.Fatalf("Created at time parse failed %v", caTErr.Error())
			}
			publishedAtTime, paTErr := time.Parse(time.RFC3339Nano, first["published_at"].(string))
			if paTErr != nil {
				log.Fatalf("Published at time parse failed: %v", paTErr.Error())
			}

			newRelease = models.Release{
				ID:                 primitive.NewObjectID(),
				ReleaseID:          first["id"].(float64),
				URL:                first["html_url"].(string),
				Tag:                first["tag_name"].(string),
				Name:               first["name"].(string),
				ReleaseCreatedAt:   createdAtTime,
				ReleasePublishedAt: publishedAtTime,
				Body:               first["body"].(string),
				RepoID:             repo.ID,
			}

			log.Printf("[Job][GetNewReleases] New release record: %v", newRelease)

			result, err := c.InsertOne(context.Background(), &newRelease)

			repo.LastReleaseID = newRelease.ID

			c = sess.Database("releasetrackr").Collection("repos")

			res := c.FindOneAndUpdate(context.Background(),
				bson.M{
					"_id": repo.ID,
				},
				bson.M{
					"$set": bson.M{
						"last_release_id": result.InsertedID,
					},
				})
			if res == nil {
				panic(err)
			}

			SendNewReleaseNotification(repo, newRelease)
		}
	}
}
