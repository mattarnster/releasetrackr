package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Release is a release for a Repo
type Release struct {
	ID                 bson.ObjectID `json:"id" bson:"_id"`
	ReleaseID          float64       `json:"release_id" bson:"release_id"`
	URL                string        `json:"html_url" bson:"url"`
	Tag                string        `json:"tag" bson:"tag"`
	Name               string        `json:"name" bson:"name"`
	ReleasePublishedAt time.Time     `json:"release_published_at" bson:"release_published_at"`
	ReleaseCreatedAt   time.Time     `json:"release_created_at" bson:"release_created_at"`
	Body               string        `json:"release_body" bson:"release_body"`
	RepoID             bson.ObjectID `json:"repo_id" bson:"repo_id"`
}
