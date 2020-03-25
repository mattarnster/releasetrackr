package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repo is a repository referenced by Track
type Repo struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Repo          string             `json:"repo" bson:"repo,omitempty"`
	LastReleaseID primitive.ObjectID `json:"last_release_id" bson:"last_release_id,omitempty"`
}
