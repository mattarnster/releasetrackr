package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Release is a release for a Repo
type Release struct {
	ID                 bson.ObjectId `json:"release_id" bson:"_id"`
	URL                string        `json:"url" bson:"url"`
	Tag                string        `json:"tag" bson:"tag"`
	Name               string        `json:"name" bson:"name"`
	ReleaseCreatedDate time.Time     `json:"release_created_at" bson:"release_created_at"`
	Body               string        `json:"release_body" bson:"release_body"`
}
