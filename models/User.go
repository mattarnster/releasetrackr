package models

import (
	"context"
	"errors"
	"releasetrackr/db"
	"releasetrackr/jwttoken"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the user model for the DB schema
type User struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email            string             `json:"email" bson:"email,omitempty" valid:"email"`
	Password         string             `json:"-" bson:"password"`
	VerificationCode string             `json:"-" bson:"verificationcode"`
	Verified         bool               `json:"verified" bson:"verified"`
	CreatedAt        time.Time          `json:"created_at" bson:"created"`
}

// GetUserFromJWT retuns a user from a JWT token
func (*User) GetUserFromJWT(token string) (*User, error) {
	sess, err := db.GetDbSession()
	if err != nil {
		return nil, errors.New("Couldn't connect to the database")
	}

	c := sess.Database("releasetrackr").Collection("users")
	var user User

	tok, tokErr := jwttoken.GetJWT(token)

	if tokErr != nil {
		return nil, errors.New("Token error")
	}

	userResult := c.FindOne(context.Background(), bson.M{
		"email": tok.Claims.(*jwttoken.CustomClaims).Email,
	})

	if userResult.Err() != nil {
		return nil, errors.New("User record not found")
	}

	userResult.Decode(&user)

	return &user, nil
}
