package models

import "gopkg.in/mgo.v2/bson"
import "time"

// User is the user model for the DB schema
type User struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email            string        `json:"email" bson:"email,omitempty" valid:"email"`
	VerificationCode string        `json:"code" bson:"verificationcode"`
	Verified         bool          `json:"verified" bson:"verified"`
	CreatedAt        time.Time     `json:"created_at" bson:"created"`
}
