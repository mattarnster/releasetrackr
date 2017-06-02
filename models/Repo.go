package models

import "gopkg.in/mgo.v2/bson"

// Repo is a repository referenced by Track
type Repo struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Repo string        `json:"repo" bson:"repo"`
}
