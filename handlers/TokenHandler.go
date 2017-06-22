package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"releasetrackr/helpers"
	"releasetrackr/models"
	"releasetrackr/requests"
	"releasetrackr/responses"
	"time"

	"gopkg.in/mgo.v2/bson"

	uuid "github.com/nu7hatch/gouuid"
)

var tokenRequest *requests.TokenRequest
var user *models.User

// TokenHandler handles the validation of TokenRequest
// and then sending out an email if the request is valid.
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "Bad request method",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(json))

		return
	}

	decoder := json.NewDecoder(r.Body)

	// If we couldn't decode the request, then we'll
	// send them back with a 400.
	err := decoder.Decode(&tokenRequest)
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

	// Get the DB session
	sess, err := helpers.GetDbSession()
	if err != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "Internal server error - E_NO_DB_CONN",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(json))

		return
	}

	c := sess.DB("releasetrackr").C("users")

	// Find the user in question by their Email address
	usrNotFoundErr := c.Find(
		bson.M{
			"email": tokenRequest.Email,
		},
	).One(&user)

	if usrNotFoundErr != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  http.StatusNotFound,
			Error: "User email not valid or not found.",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(json))

		return
	}

	log.Printf("[TokenHandler] New Token request - %s", user.Email)

	newToken, err := uuid.NewV4()
	if err != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "Internal server error - E_TOK_GEN_FAIL",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(json))

		return
	}

	change := bson.M{
		"$set": bson.M{
			"session_token":  newToken.String(),
			"session_expiry": time.Now().Add(-18 * time.Hour), // 18 hour expiry
		},
	}
	dbUpdateErr := c.Update(
		bson.M{
			"_id": user.ID,
		}, change)
	if dbUpdateErr != nil {
		log.Printf("[TokenHandler] Error assigning new token values %v", dbUpdateErr.Error())
		return
	}

	log.Printf("[TokenHandler] Got changeInfo %v", dbUpdateErr)

	json, _ := json.Marshal(&responses.TokenResponse{
		Code:  http.StatusOK,
		Token: newToken.String(),
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))

	log.Printf("[TokenHandler] New token request fulfilled")

	return
}
