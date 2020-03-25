package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the user model for the DB schema
type User struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email            string             `json:"email" bson:"email,omitempty" valid:"email"`
	VerificationCode string             `json:"code" bson:"verificationcode"`
	Verified         bool               `json:"verified" bson:"verified"`
	CreatedAt        time.Time          `json:"created_at" bson:"created"`
}
