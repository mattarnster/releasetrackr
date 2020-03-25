package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Track is the schema for the repo tracking functionality
type Track struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"userID" bson:"userID"`
	RepoID primitive.ObjectID `json:"repoID" bson:"repoID"`
}
