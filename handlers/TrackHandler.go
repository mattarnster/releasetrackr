package handlers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"releasetrackr/helpers"
	"releasetrackr/models"
	"releasetrackr/requests"
	"releasetrackr/responses"

	recaptcha "github.com/dpapathanasiou/go-recaptcha"
	uuid "github.com/nu7hatch/gouuid"
	"gopkg.in/mgo.v2/bson"
)

var repo = &models.Repo{}

// TrackHandler handles creation and verification of Track requests
func TrackHandler(w http.ResponseWriter, r *http.Request) {
	var recaptchaSecret = os.Getenv("RECAPTCHA_SECRET")

	// If the method isn't POST, then send them back
	// with a Bad Request (400)
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

	// If we couldn't decode the request, then we'll
	// send them back with a 400.
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

	// Validate the recaptcha call
	recaptcha.Init(recaptchaSecret)
	// Extract the IP from the request headers
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	// Determine whether or not they get to continue
	recaptchaResult := recaptcha.Confirm(ip, tr.RecaptchaResponse)
	// If not...
	if recaptchaResult == false {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  400,
			Error: "Recaptcha challenge failed!",
		})

		w.WriteHeader(400)
		w.Write(json)

		return
	}

	// If any of the required fields are empty, return an error.
	if tr.Email == "" || tr.Repo == "" || tr.RecaptchaResponse == "" {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  400,
			Error: "Missing required field(s)",
		})

		w.WriteHeader(400)
		w.Write(json)

		return
	}

	log.Printf("[Handler][TrackHandler] Incoming track request: %s from %s", tr.Repo, r.RemoteAddr)

	// Grab the DB session from the helpers.
	sess, err := helpers.GetDbSession()
	if err != nil {
		panic("Couldn't get DB session")
	}

	user := &models.User{}

	c := sess.DB("releasetrackr").C("users")

	// Try and find the existing user (if they exist)
	// Otherwise create the user.
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

		// If they're a new user, we'll tell them
		// that they need verification and send them off an email.
		json, _ := json.Marshal(&responses.SuccessResponse{
			Code:    403,
			Message: "Email verification required.",
		})

		helpers.SendVerificationEmail(tr.Email, verification.String())

		w.WriteHeader(403)
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

	c = sess.DB("releasetrackr").C("repos")

	// Find an existing repo by name
	repoErr := c.Find(bson.M{"repo": tr.Repo}).One(&repo)

	var isNewRepo = false
	var newRepo models.Repo

	// If we can't find that particular repo in the DB
	// then we'll make a new one.
	if repoErr != nil {
		newRepo = models.Repo{
			ID:   bson.NewObjectId(),
			Repo: tr.Repo,
		}

		isNewRepo = true

		err := c.Insert(&newRepo)
		if err != nil {
			panic("Unable to insert new repo")
		}

		log.Printf("[Handler][TrackHandler] New repo added: %s for %s", tr.Repo, tr.Email)
	}

	// If we made a new repo in the DB,
	// make sure we know what we're searching for.
	// This is a bit unnessecary to query for the newly
	// created repo, but whatever.
	var searchRepo bson.ObjectId
	if isNewRepo {
		searchRepo = newRepo.ID
	} else {
		searchRepo = repo.ID
	}

	// See if the user already has a subscription to
	// watch this repo for releases.
	c = sess.DB("releasetrackr").C("tracks")
	record := &models.Track{}
	dbtr := c.Find(
		bson.M{
			"userID": bson.ObjectId(user.ID),
			"repoID": bson.ObjectId(searchRepo),
		},
	).One(&record)

	// If they do, we will get no error from the DB,
	// but the user will get a 409 (Conflict)
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

	c = sess.DB("releasetrackr").C("tracks")

	// Make a new "track" for this repo & user.
	var repoID bson.ObjectId

	if isNewRepo {
		repoID = newRepo.ID
	} else {
		repoID = repo.ID
	}

	trModel := &models.Track{
		ID:     bson.NewObjectId(),
		UserID: user.ID,
		RepoID: repoID,
	}

	insErr := c.Insert(&trModel)
	if insErr != nil {
		log.Panicf("[Handler][TrackHandler] Could not insert new track request: %v", insErr)
	}

	log.Printf("[Handler][TrackHandler] New track request: %s from %s for %s", trModel.ID.String(), user.Email, tr.Repo)

	// Go back to the user letting them
	// know that the request was successful.
	json, _ := json.Marshal(&responses.SuccessResponse{
		Code:    201,
		Message: "Track request acknowledged.",
	})
	w.WriteHeader(201)
	w.Write(json)

	return
}
