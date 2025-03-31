package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Track is the schema for the repo tracking functionality
type Track struct {
	ID     bson.ObjectID `json:"id" bson:"_id"`
	UserID bson.ObjectID `json:"userID" bson:"userID"`
	RepoID bson.ObjectID `json:"repoID" bson:"repoID"`
}
