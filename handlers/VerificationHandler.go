package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	uuid "github.com/nu7hatch/gouuid"

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

		sessionTok, _ := uuid.NewV4()

		// If not, we'll set the verified field to true
		change := bson.M{
			"$set": bson.M{
				"verified":       true,
				"session_token":  sessionTok,
				"session_expiry": time.Now().Add(-18 * time.Hour), // 18 hour expiry
			},
		}
		c.Update(user, change)

		log.Printf("[Handler][VerificationHandler] Verification token pass: %s - %s - %s", key, user.Email, r.RemoteAddr)

		http.Redirect(w, r, os.Getenv("RT_DOMAIN")+"/?st="+sessionTok.String(), 301)

		return
	}
}
