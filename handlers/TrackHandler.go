package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mattarnster/releasetrackr/helpers"
	"github.com/mattarnster/releasetrackr/models"
	"github.com/mattarnster/releasetrackr/requests"
	"github.com/mattarnster/releasetrackr/responses"
	uuid "github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
)

// TrackHandler handles creation and verification of Track requests
func TrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  400,
			Error: "Bad Request",
		})
		w.WriteHeader(400)
		w.Write(json)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var tr = &requests.TrackRequest{}

	err := decoder.Decode(tr)
	if err != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  400,
			Error: "Request format invalid.",
		})

		w.WriteHeader(400)
		w.Write(json)

		return
	}

	defer r.Body.Close()

	log.Printf("[Handler][TrackHandler] Incoming track request: %s from %s", tr.Repo, r.RemoteAddr)

	sess, err := helpers.GetDbSession()
	if err != nil {
		panic("Couldn't get DB session")
	}

	user := &models.User{}

	c := sess.DB("releasetrackr").C("users")

	userErr := c.Find(bson.M{"email": tr.Email}).One(&user)
	if userErr != nil {

		uid := bson.NewObjectId()
		verification, _ := uuid.NewV4()

		c.Insert(&models.User{
			ID:               uid,
			Email:            tr.Email,
			VerificationCode: verification.String(),
			Verified:         false,
			CreatedAt:        time.Now(),
		})

		log.Printf("[Handler][TrackHandler] New user, needs verification: %s, %s - {%s}", uid, tr.Email, verification.String())
		helpers.SendVerificationEmail(tr.Email, verification.String())

		json, _ := json.Marshal(&responses.SuccessResponse{
			Code:    202,
			Message: "Email verification required.",
		})

		w.WriteHeader(202)
		w.Write(json)

		return

	}
	// Existing user, make sure they're verified first...
	if user.Verified == false {
		response, _ := json.Marshal(&responses.ErrorResponse{
			Code:  403,
			Error: "Verification required - Check your email.",
		})
		log.Println("[Handler][TrackHandler] Responding with verification required.")
		w.WriteHeader(403)
		w.Write(response)
		return
	}
	// Already a user, stop them from making another record of the same.
	c = sess.DB("releasetrackr").C("tracks")
	record := &models.Track{}
	dbtr := c.Find(bson.M{"userID": bson.ObjectId(user.ID), "repo": tr.Repo}).One(&record)
	if dbtr == nil {
		response, _ := json.Marshal(&responses.ErrorResponse{
			Code:  409,
			Error: "You've already subscribed to be notified about this repository.",
		})

		log.Printf("[Handler][TrackHandler] User already subscribed to repo: %s - %s", user.Email, tr.Repo)

		w.WriteHeader(409)
		w.Write(response)
		return
	}

	trID := bson.NewObjectId()
	c.Insert(&models.Track{
		ID:     trID,
		UserID: user.ID,
		Repo:   tr.Repo,
	})

	log.Printf("[Handler][TrackHandler] New track request: %s from %s for %s", trID.String(), user.Email, tr.Repo)

	json, _ := json.Marshal(&responses.SuccessResponse{
		Code:    201,
		Message: "Track request acknowledged.",
	})
	w.WriteHeader(201)
	w.Write(json)

	return
}
