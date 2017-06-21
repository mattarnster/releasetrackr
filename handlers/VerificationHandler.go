package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"releasetrackr/helpers"
	"releasetrackr/models"
	"releasetrackr/responses"

	"gopkg.in/mgo.v2/bson"
)

// VerificationHandler handles verification of user emails
func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	// If we have something in the ?key= query field...
	if key != "" {
		sess, err := helpers.GetDbSession()
		if err != nil {
			panic("test")
		}

		user := &models.User{}
		c := sess.DB("releasetrackr").C("users")

		// Find the user that the verification field corresponds to
		userErr := c.Find(bson.M{"verificationcode": key, "verified": false}).One(&user)

		// If it's invalid, display an error back to the user
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

		// If not, we'll set the verified field to true
		change := bson.M{"$set": bson.M{"verified": true}}
		c.Update(user, change)

		log.Printf("[Handler][VerificationHandler] Verification token pass: %s - %s", key, r.RemoteAddr)

		// Display a success message to the user.
		json, _ := json.Marshal(&responses.SuccessResponse{
			Code:    200,
			Message: "Verification passed.",
		})

		w.WriteHeader(200)
		w.Write(json)
		return
	}
}
