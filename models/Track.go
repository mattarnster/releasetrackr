package models

import "gopkg.in/mgo.v2/bson"

// Track is the schema for the repo tracking functionality
type Track struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	UserID bson.ObjectId `json:"userID" bson:"userID"`
	RepoID bson.ObjectId `json:"repoID" bson:"repoID"`
}
