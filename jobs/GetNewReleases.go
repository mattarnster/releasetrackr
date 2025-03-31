package jobs

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"releasetrackr/db"
	"releasetrackr/models"
)

// GetNewReleases gets new releases from the Github API
func GetNewReleases() {
	// Grab a bunch of repos
	sess, _ := db.GetDbSession()

	c := sess.Database("releasetrackr").Collection("repos")
	opts := options.Count().SetHint("_id_")
	count, err := c.CountDocuments(context.Background(), bson.D{}, opts)
	if err != nil {
		log.Printf("[Job][GetNewReleases] Failed to find repos %s", err.Error())
		return
	}

	if count == 0 {
		log.Println("[Job][GetNewReleases] No repos in DB")
		return
	}

	cur, _ := c.Find(context.Background(), bson.D{})
	defer cur.Close(context.Background())

	log.Printf("[Job][GetNewReleases] Result count: %v", count)
	// Get a list of all returned documents and print them out.
	// See the mongo.Cursor documentation for more examples of using cursors.
	var repos []models.Repo
	if err = cur.All(context.TODO(), &repos); err != nil {
		log.Panic(err)
	}
	for _, repo := range repos {
		log.Printf("[Job][GetNewReleases] Looking for release for %+v", repo.Repo)

		resp, err := http.Get("https://api.github.com/repos/" + repo.Repo + "/releases")
		if err != nil {
			log.Printf("[Job][GetNewReleases] API Request failed: %v", err.Error())
			continue
		}

		log.Printf("[Job][GetNewReleases] Github ratelimit will be hit in in %v requests.", resp.Header["X-Ratelimit-Remaining"])
		log.Printf("[Job][GetNewReleases] Ratelimit will reset at %v", resp.Header["X-Ratelimit-Reset"])

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("[Job][GetNewReleases] https://api.github.com/%s | Response Status Code: %v", repo.Repo, resp.StatusCode)
		}

		var releases []models.GitHubRelease

		body, _ := io.ReadAll(resp.Body)

		// Print the raw JSON for debugging
		log.Printf("[Job][GetNewReleases] Raw JSON response: %s", string(body))

		err = json.Unmarshal(body, &releases)
		if err != nil {
			log.Printf("[Job][GetNewReleases] Error unmarshaling JSON - likely invalid: %v", err.Error())

			// Let's try unmarshaling into a generic structure to see the actual format
			var rawData []map[string]interface{}
			jsonErr := json.Unmarshal(body, &rawData)
			if jsonErr == nil && len(rawData) > 0 {
				log.Printf("[Job][GetNewReleases] First release structure: %+v", rawData[0])
			}

			continue
		}

		if len(releases) == 0 {
			log.Printf("[Job][GetNewReleases] No releases found for %s", repo.Repo)
			continue
		}

		// Get the latest release (first in the array)
		latestRelease := releases[0]

		c = sess.Database("releasetrackr").Collection("releases")

		// Check if this release already exists in our database
		var existingRelease models.Release
		err = c.FindOne(context.Background(), bson.M{"release_id": float64(latestRelease.ID)}).Decode(&existingRelease)

		// Not found - Add it to the DB
		if err != nil {
			isNewRelease := true

			newRelease := models.Release{
				ID:                 bson.NewObjectID(),
				ReleaseID:          float64(latestRelease.ID),
				URL:                latestRelease.HTMLURL,
				Tag:                latestRelease.TagName,
				Name:               latestRelease.Name,
				ReleaseCreatedAt:   latestRelease.CreatedAt,
				ReleasePublishedAt: latestRelease.PublishedAt,
				Body:               latestRelease.Body,
				RepoID:             repo.ID,
			}

			log.Printf("[Job][GetNewReleases] New release record: %v", newRelease)

			_, err := c.InsertOne(context.Background(), &newRelease)
			if err != nil {
				log.Printf("[Job][GetNewReleases] Failed to insert new release: %v", err.Error())
				continue
			}

			// Update the repo with the new release ID
			c = sess.Database("releasetrackr").Collection("repos")

			res := c.FindOneAndUpdate(context.Background(),
				bson.M{
					"_id": repo.ID,
				},
				bson.M{
					"$set": bson.M{
						"last_release_id": newRelease.ID,
					},
				})
			if res.Err() != nil {
				log.Printf("[Job][GetNewReleases] Failed to update repo with new release ID: %v", res.Err())
				continue
			}

			if isNewRelease {
				SendNewReleaseNotification(repo, newRelease)
			}
		}
	}
}
