package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Repo is a repository referenced by Track
type Repo struct {
	ID            bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Repo          string        `json:"repo" bson:"repo,omitempty"`
	LastReleaseID bson.ObjectID `json:"last_release_id" bson:"last_release_id,omitempty"`
}
