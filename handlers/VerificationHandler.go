package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattarnster/releasetrackr/helpers"
	"github.com/mattarnster/releasetrackr/models"
	"github.com/mattarnster/releasetrackr/responses"

	"gopkg.in/mgo.v2/bson"
)

// VerificationHandler handles verification of user emails
func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key != "" {
		sess, err := helpers.GetDbSession()
		if err != nil {
			panic("test")
		}

		user := &models.User{}
		c := sess.DB("releasetrackr").C("users")

		userErr := c.Find(bson.M{"verificationcode": key, "verified": false}).One(&user)
		if userErr != nil {
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  400,
				Error: "The token you want to verify is invalid",
			})
			w.WriteHeader(400)
			w.Write(json)
			log.Printf("[Handler][VerificationHandler] Verification token fail: %s", r.RemoteAddr)
			return
		}

		change := bson.M{"$set": bson.M{"verified": true}}
		c.Update(user, change)

		log.Printf("[Handler][VerificationHandler] Verification token pass: %s - %s", key, r.RemoteAddr)

		json, _ := json.Marshal(&responses.SuccessResponse{
			Code:    200,
			Message: "Verification passed.",
		})

		w.WriteHeader(200)
		w.Write(json)
		return
	}
}
