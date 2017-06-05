package models

import "gopkg.in/mgo.v2/bson"

// Repo is a repository referenced by Track
type Repo struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Repo          string        `json:"repo" bson:"repo,omitempty"`
	LastReleaseID bson.ObjectId `json:"last_release_id" bson:"last_release_id,omitempty"`
}
